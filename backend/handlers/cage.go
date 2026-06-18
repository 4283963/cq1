package handlers

import (
	"net/http"
	"ocean-ranch/models"
	"ocean-ranch/repository"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CageHandler struct {
	repo *repository.CageRepository
}

func NewCageHandler(repo *repository.CageRepository) *CageHandler {
	return &CageHandler{repo: repo}
}

func (h *CageHandler) GetAll(c *gin.Context) {
	cages, err := h.repo.FindAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": cages,
	})
}

func (h *CageHandler) GetByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cage, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cage not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": cage,
	})
}

func (h *CageHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	cage, err := h.repo.FindByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "cage not found"})
		return
	}

	var body models.Cage
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cage.WaterTemp = body.WaterTemp
	cage.DissolvedOxygen = body.DissolvedOxygen
	cage.Status = body.Status

	if err := h.repo.Update(cage); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": cage,
	})
}
