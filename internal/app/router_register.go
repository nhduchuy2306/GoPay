package app

import (
	"gopay/internal/auth"
	"gopay/internal/mail"
	"gopay/internal/middleware"
	"gopay/internal/notification"
	"gopay/internal/transaction"
	"gopay/internal/user"
	"gopay/internal/wallet"
	"gopay/internal/webhook"
	"log"
	"os"
	"strconv"

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
	mailHost := os.Getenv("MAIL_HOST")
	if mailHost == "" {
		mailHost = "localhost"
	}
	mailPort, err := strconv.Atoi(os.Getenv("MAIL_PORT"))
	if err != nil || mailPort == 0 {
		mailPort = 1025
	}
	mailService := mail.NewService(mailHost, mailPort)
	emailSender := notification.NewEmailSender(mailService)
	notificationService := notification.NewService(emailSender)
	log.Println("Init Mail", mailService)

	// User
	userRepo := user.NewRepository(db)
	userService := user.NewService(userRepo)

	// Wallet
	walletRepo := wallet.NewRepository(db)
	walletService := wallet.NewService(walletRepo)

	// Auth
	authService := auth.NewService(userService, walletService, notificationService)
	authHandler := auth.NewHandler(authService)
	authHandler.RegisterRouters(group)

	// User
	userHandler := user.NewHandler(userService)
	userHandler.RegisterRouters(group)

	walletHandler := wallet.NewHandler(walletService)
	walletHandler.RegisterRouters(group)

	// Transaction
	transactionRepo := transaction.NewRepository(db)
	transactionService := transaction.NewService(transactionRepo, walletService, notificationService)
	transactionHandler := transaction.NewHandler(transactionService)
	transactionHandler.RegisterRouters(group)

	// Webhook
	webhookRepo := webhook.NewRepository(db)
	webhookWorker := webhook.NewWorker()
	webhookService := webhook.NewService(webhookRepo, webhookWorker)
	webhookHandler := webhook.NewHandler(webhookService)
	webhookHandler.RegisterRouters(group)
}
