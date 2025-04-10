package handler

import (
	"net/http"
	"strconv"

	"order-service/internal/model"
	"order-service/internal/usecase"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	Usecase *usecase.OrderUsecase
}

func (h *OrderHandler) RegisterRoutes(r *gin.Engine) {
	r.POST("/orders", h.Create)
	r.GET("/orders/:id", h.GetByID)
	r.PATCH("/orders/:id", h.UpdateStatus)
	r.GET("/orders", h.ListAllOrByUser)
}

func (h *OrderHandler) Create(c *gin.Context) {
	var o model.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	o.Status = "pending"
	if err := h.Usecase.Create(&o); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, o)
}

func (h *OrderHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	order, err := h.Usecase.GetByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var payload struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.Usecase.UpdateStatus(id, payload.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

func (h *OrderHandler) ListAllOrByUser(c *gin.Context) {
	userID := c.Query("user_id")
	var (
		orders []model.Order
		err    error
	)

	if userID != "" {
		orders, err = h.Usecase.ListByUser(userID)
	} else {
		orders, err = h.Usecase.ListAll()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

func (h *OrderHandler) ListByUser(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}
	orders, err := h.Usecase.ListByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
