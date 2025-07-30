package handlers

import (
	"log"
	"net/http"
	"sermorpheus-engine-test/internal/services"
	"sermorpheus-engine-test/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TransactionHandler struct {
	transactionService *services.TransactionService
	customerService    *services.CustomerService
	blockchainService  *services.BlockchainService
}

func NewTransactionHandler(transactionService *services.TransactionService, customerService *services.CustomerService, blockchainService *services.BlockchainService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
		customerService:    customerService,
		blockchainService:  blockchainService,
	}
}

type CreateTransactionRequest struct {
	CustomerEmail string    `json:"customer_email" binding:"required,email"`
	CustomerName  string    `json:"customer_name" binding:"required"`
	CustomerPhone string    `json:"customer_phone"`
	EventID       uuid.UUID `json:"event_id" binding:"required"`
	Quantity      int       `json:"quantity" binding:"required,gt=0"`
}

func (th *TransactionHandler) CreateTransaction(c *gin.Context) {
	var req CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	customer, err := th.customerService.GetOrCreateCustomer(req.CustomerEmail, req.CustomerName, req.CustomerPhone)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to process customer", err.Error())
		return
	}

	transactionReq := &services.CreateTransactionRequest{
		CustomerID: customer.ID,
		EventID:    req.EventID,
		Quantity:   req.Quantity,
	}

	transaction, err := th.transactionService.CreateTransaction(transactionReq)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to create transaction", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Transaction created successfully", gin.H{
		"transaction":      transaction,
		"payment_address":  transaction.PaymentAddress,
		"usdt_amount":      transaction.USDTAmount,
		"payment_deadline": transaction.PaymentLockedAt.Add(30 * 60 * 1000000000),
	})
}

func (th *TransactionHandler) GetTransaction(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	transaction, err := th.transactionService.GetTransactionByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Transaction not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Transaction retrieved successfully", transaction)
}

func (th *TransactionHandler) ConfirmPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	var req struct {
		TxHash string  `json:"tx_hash" binding:"required"`
		Amount float64 `json:"amount,omitempty"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	if err := th.transactionService.ConfirmPayment(id, req.TxHash); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Failed to confirm payment", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Payment confirmed successfully", nil)
}

func (th *TransactionHandler) CheckPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid transaction ID", err.Error())
		return
	}

	transaction, err := th.transactionService.GetTransactionByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Transaction not found", err.Error())
		return
	}

	if transaction.Status != "pending" {
		utils.SuccessResponse(c, http.StatusOK, "Transaction already processed", gin.H{
			"status": transaction.Status,
		})
		return
	}

	go func() {
		found := th.blockchainService.CheckRecentTransfer(transaction.ID, transaction.USDTAmount, transaction.PaymentAddress)
		if found {
			log.Printf("Manual check found payment for transaction %s", transaction.ID)
		} else {
			log.Printf("Manual check: no payment found for transaction %s", transaction.ID)
		}
	}()

	utils.SuccessResponse(c, http.StatusOK, "Payment check initiated", gin.H{
		"transaction_id": transaction.ID,
		"checking":       true,
	})
}
