package router

import (
	"go_starter/internal/handler"
	"go_starter/internal/middleware"
	"go_starter/internal/repository"
	"go_starter/internal/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg interface{}, logger *zap.Logger) {
	api := r.Group("/api")

	// User module
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc, logger)
	authHandler := handler.NewAuthHandler(userSvc, logger)
	userGroup := api.Group("/users")
	{
		userGroup.POST("", userHandler.Create)
		userGroup.GET("", userHandler.List)
		userGroup.GET("/paginate", userHandler.Paginate)
		userGroup.GET("/:id", userHandler.GetById)
		userGroup.PUT("/:id", userHandler.Update)
		userGroup.DELETE("/:id", userHandler.Delete)
	}

	// Auth module
	authGroup := api.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)
		authGroup.GET("/profile", middleware.AuthMiddleware(), authHandler.GetProfile)
	}

}
