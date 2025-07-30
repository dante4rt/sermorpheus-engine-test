package services

import (
	"errors"
	"sermorpheus-engine-test/internal/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventService struct {
	db *gorm.DB
}

func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

func (es *EventService) CreateEvent(event *models.Event) error {
	event.AvailableQuota = event.Quota

	if err := es.db.Create(event).Error; err != nil {
		return err
	}
	return nil
}

func (es *EventService) GetEventByID(id uuid.UUID) (*models.Event, error) {
	var event models.Event
	if err := es.db.First(&event, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &event, nil
}

func (es *EventService) GetEvents(limit, offset int) ([]models.Event, error) {
	var events []models.Event
	query := es.db.Order("created_at DESC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&events).Error; err != nil {
		return nil, err
	}
	return events, nil
}

func (es *EventService) UpdateEventQuota(eventID uuid.UUID, quantity int) error {
	return es.db.Transaction(func(tx *gorm.DB) error {
		var event models.Event
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			First(&event, "id = ?", eventID).Error; err != nil {
			return err
		}

		if event.AvailableQuota < quantity {
			return errors.New("insufficient tickets available")
		}

		event.AvailableQuota -= quantity
		return tx.Save(&event).Error
	})
}

func (es *EventService) RestoreEventQuota(eventID uuid.UUID, quantity int) error {
	return es.db.Transaction(func(tx *gorm.DB) error {
		var event models.Event
		if err := tx.Set("gorm:query_option", "FOR UPDATE").
			First(&event, "id = ?", eventID).Error; err != nil {
			return err
		}

		event.AvailableQuota += quantity
		if event.AvailableQuota > event.Quota {
			event.AvailableQuota = event.Quota
		}

		return tx.Save(&event).Error
	})
}
