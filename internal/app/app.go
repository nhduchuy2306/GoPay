package app

import (
	"gopay/internal/album"
	"gopay/internal/auth"
	"gopay/internal/user"
	"gopay/pkg"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type App struct {
	engine *gin.Engine
}

func New() *App {
	engine := gin.Default()

	// Check Cors
	engine.Use(cors.Default())

	// Connect to database
	db := pkg.CreateDbConnection()

	// Create Api Group
	api := engine.Group("/api/v1")

	// Register modules
	registerAllModules(db, api)

	return &App{engine: engine}
}

func (a *App) Run(addr string) error {
	return a.engine.Run(addr)
}

func registerAllModules(db *gorm.DB, group *gin.RouterGroup) {
	// Album
	albumRepo := album.NewRepository(db)
	albumService := album.NewService(albumRepo)
	albumHandler := album.NewHandler(albumService)
	albumHandler.RegisterRouters(group)

	// User
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRouters(group)

	// Auth
	authService := auth.NewService(userService)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRouters(group)
}
