package app

import (
	"context"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func New(ctx context.Context) *App {
	// Init Env
	InitEnv()

	// Init Engine
	engine := gin.Default()

	// Check Cors
	engine.Use(cors.Default())

	// Init database
	db := InitDB()

	// Start background workers
	StartBackgroundWorkers(ctx, db)

	// Register modules
	RegisterAllModules(db, engine)

	return &App{engine: engine}
}

func (a *App) Run(addr string) error {
	return a.engine.Run(addr)
}
