package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Event struct {
	ID             uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name           string          `gorm:"not null" json:"name"`
	Description    string          `json:"description"`
	Location       string          `gorm:"not null" json:"location"`
	Schedule       time.Time       `gorm:"not null" json:"schedule"`
	PriceIDR       float64         `gorm:"not null" json:"price_idr"`
	Quota          int             `gorm:"not null" json:"quota"`
	AvailableQuota int             `gorm:"not null" json:"available_quota"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
	DeletedAt      *gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
	Transactions   []Transaction   `json:"transactions,omitempty"`
	Tickets        []Ticket        `json:"tickets,omitempty"`
}

type Customer struct {
	ID           uuid.UUID     `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Email        string        `gorm:"uniqueIndex;not null" json:"email"`
	Name         string        `gorm:"not null" json:"name"`
	Phone        string        `json:"phone"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	Transactions []Transaction `json:"transactions,omitempty"`
	Tickets      []Ticket      `json:"tickets,omitempty"`
}

type Transaction struct {
	ID                     uuid.UUID               `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CustomerID             uuid.UUID               `gorm:"type:uuid;not null" json:"customer_id"`
	EventID                uuid.UUID               `gorm:"type:uuid;not null" json:"event_id"`
	Quantity               int                     `gorm:"not null" json:"quantity"`
	TotalIDR               float64                 `gorm:"not null" json:"total_idr"`
	USDTRate               float64                 `gorm:"not null" json:"usdt_rate"`
	USDTAmount             float64                 `gorm:"not null" json:"usdt_amount"`
	PaymentAddress         string                  `json:"payment_address"`
	Status                 string                  `gorm:"default:'pending'" json:"status"`
	PaymentLockedAt        *time.Time              `json:"payment_locked_at"`
	PaymentConfirmedAt     *time.Time              `json:"payment_confirmed_at"`
	CreatedAt              time.Time               `json:"created_at"`
	UpdatedAt              time.Time               `json:"updated_at"`
	Customer               Customer                `json:"customer,omitempty"`
	Event                  Event                   `json:"event,omitempty"`
	Tickets                []Ticket                `json:"tickets,omitempty"`
	BlockchainTransactions []BlockchainTransaction `json:"blockchain_transactions,omitempty"`
}

type Ticket struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID   `gorm:"type:uuid;not null" json:"transaction_id"`
	EventID       uuid.UUID   `gorm:"type:uuid;not null" json:"event_id"`
	CustomerID    uuid.UUID   `gorm:"type:uuid;not null" json:"customer_id"`
	TicketCode    string      `gorm:"uniqueIndex;not null" json:"ticket_code"`
	Status        string      `gorm:"default:'active'" json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Transaction   Transaction `json:"transaction,omitempty"`
	Event         Event       `json:"event,omitempty"`
	Customer      Customer    `json:"customer,omitempty"`
}

type PaymentAddress struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Address    string    `gorm:"uniqueIndex;not null" json:"address"`
	PrivateKey string    `gorm:"not null" json:"private_key"`
	IsUsed     bool      `gorm:"default:false" json:"is_used"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type USDTRate struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	IDRToUSDTRate float64   `gorm:"not null" json:"idr_to_usdt_rate"`
	CreatedAt     time.Time `json:"created_at"`
}

type BlockchainTransaction struct {
	ID            uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TransactionID uuid.UUID   `gorm:"type:uuid;not null" json:"transaction_id"`
	TxHash        string      `gorm:"uniqueIndex" json:"tx_hash"`
	FromAddress   string      `json:"from_address"`
	ToAddress     string      `json:"to_address"`
	Amount        float64     `json:"amount"`
	Confirmations int         `gorm:"default:0" json:"confirmations"`
	Status        string      `gorm:"default:'pending'" json:"status"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	Transaction   Transaction `json:"transaction,omitempty"`
}
