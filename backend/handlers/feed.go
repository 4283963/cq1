package handlers

import (
	"net/http"
	"ocean-ranch/models"
	"ocean-ranch/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type FeedHandler struct {
	cageRepo *repository.CageRepository
	feedRepo *repository.FeedRepository
}

func NewFeedHandler(cageRepo *repository.CageRepository, feedRepo *repository.FeedRepository) *FeedHandler {
	return &FeedHandler{cageRepo: cageRepo, feedRepo: feedRepo}
}

func (h *FeedHandler) Create(c *gin.Context) {
	var req models.FeedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cage, err := h.cageRepo.FindByID(req.CageID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cage not found"})
		return
	}

	operator := req.Operator
	if operator == "" {
		operator = "web_user"
	}

	record := &models.FeedRecord{
		CageID:     cage.ID,
		FeedAmount: req.FeedAmount,
		Operator:   operator,
		CreatedAt:  time.Now(),
	}

	if err := h.feedRepo.Create(record); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "投喂指令已发送",
		"data":    record,
	})
}

func (h *FeedHandler) GetByCageID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("cage_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cage_id"})
		return
	}

	records, err := h.feedRepo.FindByCageID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": records,
	})
}

func (h *FeedHandler) GetAll(c *gin.Context) {
	records, err := h.feedRepo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": records,
	})
}
