package router

import (
	"go_starter/internal/handler"
	"go_starter/internal/repository"
	"go_starter/internal/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRoutes(r *gin.Engine, db *gorm.DB, cfg interface{}) {
	api := r.Group("/api")

	// User module
	userRepo := repository.NewUserRepository(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)

	userGroup := api.Group("/user")
	{
		userGroup.POST("", userHandler.Create)
		userGroup.GET("", userHandler.List)
		userGroup.GET("/paginate", userHandler.Paginate)
		userGroup.GET("/:id", userHandler.GetById)
		userGroup.PUT("/:id", userHandler.Update)
		userGroup.DELETE("/:id", userHandler.Delete)
	}

}
