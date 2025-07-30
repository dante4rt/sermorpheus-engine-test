package services

import (
	"log"
	"sermorpheus-engine-test/internal/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DatabaseService struct {
	DB *gorm.DB
}

func NewDatabaseService(databaseURL string) *DatabaseService {
	db, err := gorm.Open(postgres.Open(databaseURL), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = db.AutoMigrate(
		&models.Event{},
		&models.Customer{},
		&models.Transaction{},
		&models.Ticket{},
		&models.PaymentAddress{},
		&models.USDTRate{},
		&models.BlockchainTransaction{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Database connected and migrated successfully")
	return &DatabaseService{DB: db}
}

func (ds *DatabaseService) Close() {
	sqlDB, err := ds.DB.DB()
	if err != nil {
		log.Println("Error getting database instance:", err)
		return
	}
	sqlDB.Close()
}
