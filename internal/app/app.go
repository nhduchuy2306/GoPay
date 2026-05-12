package app

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
}

func New() *App {
	// Init Env
	InitEnv()

	// Init Engine
	engine := gin.Default()

	// Check Cors
	engine.Use(cors.Default())

	// Init database
	db := InitDB()

	// Register modules
	RegisterAllModules(db, engine)

	return &App{engine: engine}
}

func (a *App) Run(addr string) error {
	return a.engine.Run(addr)
}
