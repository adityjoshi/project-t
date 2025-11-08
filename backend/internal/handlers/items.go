package handlers

import (
	"net/http"
	"synapse/internal/models"
	"synapse/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ItemHandler struct {
	itemService     *services.ItemService
	relationService *services.RelationService
}

func NewItemHandler(itemService *services.ItemService, relationService *services.RelationService) *ItemHandler {
	return &ItemHandler{
		itemService:     itemService,
		relationService: relationService,
	}
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	var req models.CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate required fields
	if req.Title == "" && req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "title or content is required"})
		return
	}

	// Set default type
	if req.Type == "" {
		req.Type = "text"
	}

	item, err := h.itemService.CreateItem(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	item, err := h.itemService.GetItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) GetAllItems(c *gin.Context) {
	items, err := h.itemService.GetAllItems(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Ensure we always return an array, not null
	if items == nil {
		items = []models.Item{}
	}

	c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.itemService.DeleteItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item deleted"})
}

func (h *ItemHandler) GetRelatedItems(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	limit := 5
	related, err := h.relationService.FindRelatedItems(c.Request.Context(), id, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, related)
}

func (h *ItemHandler) RefreshImage(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.itemService.RefreshImageForItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return updated item
	item, err := h.itemService.GetItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) RefreshSummary(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	if err := h.itemService.RefreshSummaryForItem(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return current item (summary will be updated asynchronously)
	item, err := h.itemService.GetItem(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Summary regeneration started",
		"item":    item,
	})
}

