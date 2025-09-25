package main

import (
	"go_starter/internal/model"
	_ "go_starter/internal/repository"
	"go_starter/internal/router"
	_ "go_starter/internal/service"
	"go_starter/internal/util"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := util.LoadENV()
	db := util.ConnectDB(cfg)
	model.AutoMigrate(db)

	// Khởi tạo Gin
	r := gin.Default()

	// Đăng ký router (truyền vào DB, config nếu cần)
	router.SetupRoutes(r, db, cfg)

	// Run server
	r.Run(":" + cfg.Server.Port)
}
