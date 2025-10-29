package handler

import (
	"go_starter/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	emailService *service.EmailService
}

func NewEmailHandler() *EmailHandler {
	return &EmailHandler{
		emailService: service.NewEmailService(),
	}
}

func (h *EmailHandler) SendTestEmail(c *gin.Context) {
	var req struct {
		To      string `json:"to" binding:"required,email"`
		Subject string `json:"subject" binding:"required"`
		Body    string `json:"body" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.emailService.SendMail(req.To, req.Subject, req.Body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send email:" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Message": "email sent successfully"})
}
