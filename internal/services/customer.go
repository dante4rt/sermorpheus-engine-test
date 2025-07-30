package services

import (
	"sermorpheus-engine-test/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CustomerService struct {
	db *gorm.DB
}

func NewCustomerService(db *gorm.DB) *CustomerService {
	return &CustomerService{db: db}
}

func (cs *CustomerService) CreateCustomer(customer *models.Customer) error {
	return cs.db.Create(customer).Error
}

func (cs *CustomerService) GetCustomerByID(id uuid.UUID) (*models.Customer, error) {
	var customer models.Customer
	if err := cs.db.First(&customer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (cs *CustomerService) GetCustomerByEmail(email string) (*models.Customer, error) {
	var customer models.Customer
	if err := cs.db.First(&customer, "email = ?", email).Error; err != nil {
		return nil, err
	}
	return &customer, nil
}

func (cs *CustomerService) GetOrCreateCustomer(email, name, phone string) (*models.Customer, error) {
	customer, err := cs.GetCustomerByEmail(email)
	if err == nil {
		return customer, nil
	}

	if err == gorm.ErrRecordNotFound {
		newCustomer := &models.Customer{
			Email: email,
			Name:  name,
			Phone: phone,
		}

		if err := cs.CreateCustomer(newCustomer); err != nil {
			return nil, err
		}
		return newCustomer, nil
	}

	return nil, err
}
