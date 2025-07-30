package handlers

import (
	"net/http"
	"sermorpheus-engine-test/internal/services"
	"sermorpheus-engine-test/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CustomerHandler struct {
	customerService *services.CustomerService
}

func NewCustomerHandler(customerService *services.CustomerService) *CustomerHandler {
	return &CustomerHandler{customerService: customerService}
}

func (ch *CustomerHandler) GetCustomer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid customer ID", err.Error())
		return
	}

	customer, err := ch.customerService.GetCustomerByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Customer not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Customer retrieved successfully", customer)
}
