package handlers

import (
	"net/http"
	"sermorpheus-engine-test/internal/models"
	"sermorpheus-engine-test/internal/services"
	"sermorpheus-engine-test/internal/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EventHandler struct {
	eventService *services.EventService
}

func NewEventHandler(eventService *services.EventService) *EventHandler {
	return &EventHandler{eventService: eventService}
}

type CreateEventRequest struct {
	Name        string  `json:"name" binding:"required"`
	Description string  `json:"description"`
	Location    string  `json:"location" binding:"required"`
	Schedule    string  `json:"schedule" binding:"required"`
	PriceIDR    float64 `json:"price_idr" binding:"required,gt=0"`
	Quota       int     `json:"quota" binding:"required,gt=0"`
}

func (eh *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid request data", err.Error())
		return
	}

	schedule, err := utils.ParseTimeISO(req.Schedule)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid schedule format", "Use ISO 8601 format")
		return
	}

	event := &models.Event{
		Name:        req.Name,
		Description: req.Description,
		Location:    req.Location,
		Schedule:    *schedule,
		PriceIDR:    req.PriceIDR,
		Quota:       req.Quota,
	}

	if err := eh.eventService.CreateEvent(event); err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to create event", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, "Event created successfully", event)
}

func (eh *EventHandler) GetEvents(c *gin.Context) {

	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		offset = 0
	}

	events, err := eh.eventService.GetEvents(limit, offset)
	if err != nil {
		utils.ErrorResponse(c, http.StatusInternalServerError, "Failed to fetch events", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Events retrieved successfully", gin.H{
		"events": events,
		"limit":  limit,
		"offset": offset,
	})
}

func (eh *EventHandler) GetEventByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "Invalid event ID", err.Error())
		return
	}

	event, err := eh.eventService.GetEventByID(id)
	if err != nil {
		utils.ErrorResponse(c, http.StatusNotFound, "Event not found", err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Event retrieved successfully", event)
}
