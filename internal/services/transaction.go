package services

import (
	"errors"
	"fmt"
	"math"
	"sermorpheus-engine-test/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionService struct {
	db           *gorm.DB
	eventService *EventService
	rateService  *RateService
	platformFee  float64
}

func NewTransactionService(db *gorm.DB, eventService *EventService, rateService *RateService, platformFee float64) *TransactionService {
	return &TransactionService{
		db:           db,
		eventService: eventService,
		rateService:  rateService,
		platformFee:  platformFee,
	}
}

type CreateTransactionRequest struct {
	CustomerID uuid.UUID `json:"customer_id"`
	EventID    uuid.UUID `json:"event_id"`
	Quantity   int       `json:"quantity"`
}

func (ts *TransactionService) CreateTransaction(req *CreateTransactionRequest) (*models.Transaction, error) {
	var result *models.Transaction

	err := ts.db.Transaction(func(tx *gorm.DB) error {

		event, err := ts.eventService.GetEventByID(req.EventID)
		if err != nil {
			return errors.New("event not found")
		}

		if event.AvailableQuota < req.Quantity {
			return errors.New("insufficient tickets available")
		}

		totalIDR := event.PriceIDR * float64(req.Quantity)

		rate, err := ts.rateService.GetCurrentRate()
		if err != nil {
			return errors.New("failed to get exchange rate")
		}

		baseUSDTAmount := totalIDR / rate.IDRToUSDTRate
		feeAmount := baseUSDTAmount * (ts.platformFee / 100)
		finalUSDTAmount := baseUSDTAmount + feeAmount

		finalUSDTAmount = math.Round(finalUSDTAmount*1000000) / 1000000

		if err := ts.eventService.UpdateEventQuota(req.EventID, req.Quantity); err != nil {
			return err
		}

		paymentAddr, err := ts.getAvailablePaymentAddress(tx)
		if err != nil {
			return errors.New("no payment address available")
		}

		transaction := &models.Transaction{
			CustomerID:      req.CustomerID,
			EventID:         req.EventID,
			Quantity:        req.Quantity,
			TotalIDR:        totalIDR,
			USDTRate:        rate.IDRToUSDTRate,
			USDTAmount:      finalUSDTAmount,
			PaymentAddress:  paymentAddr,
			Status:          "pending",
			PaymentLockedAt: func() *time.Time { t := time.Now(); return &t }(),
		}

		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		if err := ts.generateTickets(tx, transaction); err != nil {
			return err
		}

		result = transaction
		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (ts *TransactionService) GetTransactionByID(id uuid.UUID) (*models.Transaction, error) {
	var transaction models.Transaction
	if err := ts.db.Preload("Customer").Preload("Event").Preload("Tickets").
		First(&transaction, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (ts *TransactionService) UpdateTransactionStatus(id uuid.UUID, status string) error {
	return ts.db.Model(&models.Transaction{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (ts *TransactionService) ConfirmPayment(transactionID uuid.UUID, txHash string) error {
	return ts.db.Transaction(func(tx *gorm.DB) error {
		var transaction models.Transaction
		if err := tx.First(&transaction, "id = ?", transactionID).Error; err != nil {
			return err
		}

		if transaction.Status != "pending" {
			return errors.New("transaction is not in pending status")
		}

		now := time.Now()
		transaction.Status = "paid"
		transaction.PaymentConfirmedAt = &now

		if err := tx.Save(&transaction).Error; err != nil {
			return err
		}

		blockchainTx := &models.BlockchainTransaction{
			TransactionID: transactionID,
			TxHash:        txHash,
			ToAddress:     transaction.PaymentAddress,
			Amount:        transaction.USDTAmount,
			Status:        "confirmed",
		}

		return tx.Create(blockchainTx).Error
	})
}

func (ts *TransactionService) getAvailablePaymentAddress(tx *gorm.DB) (string, error) {
	var paymentAddr models.PaymentAddress
	if err := tx.Where("is_used = ?", false).First(&paymentAddr).Error; err != nil {
		return "", err
	}

	paymentAddr.IsUsed = true
	if err := tx.Save(&paymentAddr).Error; err != nil {
		return "", err
	}

	return paymentAddr.Address, nil
}

func (ts *TransactionService) generateTickets(tx *gorm.DB, transaction *models.Transaction) error {
	for i := 0; i < transaction.Quantity; i++ {
		ticket := &models.Ticket{
			TransactionID: transaction.ID,
			EventID:       transaction.EventID,
			CustomerID:    transaction.CustomerID,
			TicketCode:    ts.generateTicketCode(transaction.ID, i),
			Status:        "active",
		}

		if err := tx.Create(ticket).Error; err != nil {
			return err
		}
	}
	return nil
}

func (ts *TransactionService) generateTicketCode(transactionID uuid.UUID, index int) string {
	return fmt.Sprintf("TIX-%s-%d", transactionID.String()[:8], index+1)
}
