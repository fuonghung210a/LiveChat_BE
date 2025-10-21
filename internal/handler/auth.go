package handler

import (
	"go_starter/internal/service"
	"go_starter/internal/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AuthHandler struct {
	userSvc *service.UserService
	logger  *zap.Logger
}

func NewAuthHandler(userSvc *service.UserService, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		userSvc: userSvc,
		logger:  logger,
	}
}

// Register creates a new user account
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid registration request",
			zap.String("error", err.Error()),
			zap.String("email", req.Email),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("User registration attempt",
		zap.String("email", req.Email),
		zap.String("name", req.Name),
	)

	// Create user (password will be hashed in service layer)
	user, err := h.userSvc.CreateUser(req.Name, req.Email, req.Password)
	if err != nil {
		h.logger.Error("Failed to create user",
			zap.String("error", err.Error()),
			zap.String("email", req.Email),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate token",
			zap.String("error", err.Error()),
			zap.Uint("user_id", user.ID),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	h.logger.Info("User registered successfully",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
	)

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// Login authenticates a user and returns a JWT token
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid login request",
			zap.String("error", err.Error()),
		)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.logger.Info("Login attempt",
		zap.String("email", req.Email),
	)

	// Find user by email
	user, err := h.userSvc.GetUserByEmail(req.Email)
	if err != nil {
		h.logger.Warn("Login failed - user not found",
			zap.String("email", req.Email),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Verify password
	if !util.CheckPassword(req.Password, user.Password) {
		h.logger.Warn("Login failed - invalid password",
			zap.String("email", req.Email),
			zap.Uint("user_id", user.ID),
		)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := util.GenerateToken(user.ID, user.Email)
	if err != nil {
		h.logger.Error("Failed to generate token",
			zap.String("error", err.Error()),
			zap.Uint("user_id", user.ID),
		)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	h.logger.Info("Login successful",
		zap.Uint("user_id", user.ID),
		zap.String("email", user.Email),
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "login successful",
		"token":   token,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

// GetProfile returns the current user's profile
func (h *AuthHandler) GetProfile(c *gin.Context) {
	// Get user ID from context (set by AuthMiddleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("Unauthorized profile access attempt")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	h.logger.Debug("Fetching user profile",
		zap.Any("user_id", userID),
	)

	user, err := h.userSvc.GetUserById(int64(userID.(uint)))
	if err != nil {
		h.logger.Error("Failed to fetch user profile",
			zap.String("error", err.Error()),
			zap.Any("user_id", userID),
		)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	h.logger.Info("User profile retrieved",
		zap.Uint("user_id", user.ID),
	)

	c.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}
