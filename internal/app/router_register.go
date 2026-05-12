package app

import (
	"gopay/internal/auth"
	"gopay/internal/mail"
	"gopay/internal/middleware"
	"gopay/internal/transaction"
	"gopay/internal/user"
	"gopay/internal/wallet"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ModuleRegistrar func(db *gorm.DB, group *gin.RouterGroup)

var prefixRouterRegisterMap = map[string]ModuleRegistrar{
	"v1": registerAllV1,
}

func RegisterAllModules(db *gorm.DB, engine *gin.Engine) {
	api := engine.Group("/api", middleware.BearerMiddleWare())
	for apiType, register := range prefixRouterRegisterMap {
		group := api.Group("/" + apiType)
		register(db, group)
	}
}

func registerAllV1(db *gorm.DB, group *gin.RouterGroup) {
	// Mail
	mailService := mail.NewService("localhost", 1025)
	log.Println("Init Mail", mailService)

	// User
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRouters(group)

	// Auth
	authService := auth.NewService(userService)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRouters(group)

	// Wallet
	walletRepo := wallet.NewRepository(db)
	walletService := wallet.NewService(walletRepo)
	walletHandler := wallet.NewHandler(walletService)
	walletHandler.RegisterRouters(group)

	// Transaction
	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo)
	transactionHandler := transaction.NewHandler(transactionService)
	transactionHandler.RegisterRouters(group)
}
