package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Gin
	"github.com/gin-gonic/gin"
	"github.com/socialpay/socialpay/src/pkg/config"
	"github.com/socialpay/socialpay/src/pkg/env"

	// Database

	// Repositories
	transactionRepo "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"

	// Use Cases
	SocialpayUsecase "github.com/socialpay/socialpay/src/pkg/socialpayapi/usecase"
	tipService "github.com/socialpay/socialpay/src/pkg/socialpayapi/usecase"

	// Controllers

	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	paymentController "github.com/socialpay/socialpay/src/pkg/shared/payment/controller/gin"
	SocialpayController "github.com/socialpay/socialpay/src/pkg/socialpayapi/adapter/controller/gin"

	// Payment Processors
	awashProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/awash"
	cbeProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/cbe"
	cybersourceProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/cybersource"
	etswitch "github.com/socialpay/socialpay/src/pkg/shared/payment/etswitch"
	mpesaProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/mpesa"
	telebirrProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/telebirr"

	// Transaction Types
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"

	// Swagger
	// _ "github.com/socialpay/socialpay/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// [API KEY]
	apikeyHandler "github.com/socialpay/socialpay/src/pkg/apikey_mgmt/adapter/controller/gin"
	apikeyRepo "github.com/socialpay/socialpay/src/pkg/apikey_mgmt/adapter/gateway/repo/sqlc"
	apikeyUsecase "github.com/socialpay/socialpay/src/pkg/apikey_mgmt/usecase"

	// [V2 MERCHANT]
	v2MerchantHandler "github.com/socialpay/socialpay/src/pkg/v2_merchant/adapter/controller/gin"
	v2MerchantRepo "github.com/socialpay/socialpay/src/pkg/v2_merchant/adapter/gateway/repo/sqlc"
	v2MerchantUsecase "github.com/socialpay/socialpay/src/pkg/v2_merchant/usecase"

	// [TRANSACTION]
	transactionHandler "github.com/socialpay/socialpay/src/pkg/transaction/adapter/controller/gin"
	transactionUsecase "github.com/socialpay/socialpay/src/pkg/transaction/usecase"

	// [MIDDLEWARE]
	"github.com/socialpay/socialpay/src/pkg/shared/middleware"

	// [AUTH]
	authRepo "github.com/socialpay/socialpay/src/pkg/auth/adapter/gateway/repo/psql"
	"github.com/socialpay/socialpay/src/pkg/auth/infra/storage/psql"
	authUsecase "github.com/socialpay/socialpay/src/pkg/auth/usecase"

	// SMS
	"github.com/socialpay/socialpay/src/pkg/auth/adapter/gateway/sms"

	// [WEBHOOK]
	webhookController "github.com/socialpay/socialpay/src/pkg/webhook/adapter/controller"
	webhookConsumer "github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/kafka/consumer"
	webhookRepo "github.com/socialpay/socialpay/src/pkg/webhook/adapter/gateway/repository"
	webhookUsecase "github.com/socialpay/socialpay/src/pkg/webhook/usecase"

	// [QR]
	qrHandler "github.com/socialpay/socialpay/src/pkg/qr/adapter/controller/gin"
	qrRepo "github.com/socialpay/socialpay/src/pkg/qr/core/repository"
	qrUsecase "github.com/socialpay/socialpay/src/pkg/qr/usecase"

	// Add this import

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// [WALLET]
	controller "github.com/socialpay/socialpay/src/pkg/wallet/adapter/controller"
	repository "github.com/socialpay/socialpay/src/pkg/wallet/adapter/gateway/repository"
	usecase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"

	// [NOTIFICATIONS]
	notifications "github.com/socialpay/socialpay/src/pkg/notifications"
	notificationUsecase "github.com/socialpay/socialpay/src/pkg/notifications/usecase"
)

// @title           Social Pay API V2
// @version         2.0
// @description     Social Pay API V2 documentation
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@socialpay.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v2

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
// @description API key authentication for API key management endpoints

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication for protected endpoints. Use format: "Bearer {token}"
// @bearerFormat JWT

// @securityDefinitions.apikey MerchantID
// @in header
// @name X-MERCHANT-ID
// @description Merchant ID authentication for merchant endpoints

// @security BearerAuth

func main() {
	// Check environment variables first
	env.LoadEnv()
	env.CheckEnv()

	log := log.New(os.Stdout, "[SOCIALPAY-V2]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Output Adapters]
	// [DB] Postgres
	db, err := psql.New(log)
	if err != nil {
		log.Fatal(err.Error())
	}

	log.Println("Database connected")

	// Run database migrations
	log.Println("Running database migrations...")

	// Get current working directory
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}
	log.Printf("Current working directory: %s", workDir)

	// // Construct absolute path for migrations
	// migrationsPath := fmt.Sprintf("file://%s/db/migrations", workDir)
	// log.Printf("Migrations path: %s", migrationsPath)

	// Construct database URL
	// dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_PASS"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("SSL_MODE"),
	// )
	// log.Printf("Database URL: postgres://%s:***@%s:%s/%s?sslmode=%s",
	// 	os.Getenv("DB_USER"),
	// 	os.Getenv("DB_HOST"),
	// 	os.Getenv("DB_PORT"),
	// 	os.Getenv("DB_NAME"),
	// 	os.Getenv("SSL_MODE"),
	// )

	// // Check if migrations directory exists
	// if _, err := os.Stat(fmt.Sprintf("%s/db/migrations", workDir)); os.IsNotExist(err) {
	// 	log.Fatal("Migrations directory does not exist at:", fmt.Sprintf("%s/db/migrations", workDir))
	// }

	// Initialize migrate instance
	// log.Println("Initializing migrate...")
	// m, err := migrate.New(
	// 	migrationsPath,
	// 	dbURL,
	// )
	// if err != nil {
	// 	log.Fatal("Failed to initialize migrations:", err)
	// }

	// // Run migrations
	// log.Println("Running migrations...")
	// if err := m.Up(); err != nil && err != migrate.ErrNoChange {
	// 	log.Fatal("Failed to run migrations:", err)
	// }

	// log.Println("Database migrations completed successfully")

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize Auth components (shared with v1)
	_authRepo, err := authRepo.NewPsqlRepo(log, db)
	if err != nil {
		log.Fatal("Failed to initialize auth repository:", err)
	}

	_authUseCase := authUsecase.New(log, _authRepo, sms.New(log))

	// Beign [API KEY]
	_apikeyRepo := apikeyRepo.NewRepository(db)
	_apikeyUseCase := apikeyUsecase.NewAPIKeyUseCase(_apikeyRepo)

	// Initialize middleware provider
	middlewareProvider := middleware.NewMiddlewareProvider(_authUseCase, _apikeyUseCase)

	// Configure CORS at router level
	router.Use(middlewareProvider.CORS)

	// Handle OPTIONS requests globally
	router.Use(func(c *gin.Context) {
		if c.Request.Method == "OPTIONS" {
			c.Status(http.StatusNoContent)
			c.Abort()
			return
		}
		c.Next()
	})

	// Serve static files

	// Setup Swagger
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v2 routes
	v2 := router.Group("/api/v2")

	v2.Static("/static", "./public")

	// [TRANSACTION]
	_transactionRepo := transactionRepo.NewTransactionRepository(db)
	_hostedPaymentRepo := transactionRepo.NewHostedPaymentRepository(db)
	_transactionUseCase := transactionUsecase.NewTransactionUsecase(
		_transactionRepo,
	)
	_transactionHandler := transactionHandler.NewTransactionHistoryHandler(
		_transactionUseCase,
		middlewareProvider.JWTAuth,
	)

	// Register transaction routes
	_transactionHandler.RegisterRoutes(v2)

	// [API KEY]
	_apikeyHandler := apikeyHandler.NewHandler(
		_apikeyUseCase,
		_apikeyRepo,
		middlewareProvider.JWTAuth,
		middlewareProvider.Public,
	)
	_apikeyHandler.RegisterRoutes(v2)

	// [V2 MERCHANT]
	_v2MerchantRepo := v2MerchantRepo.NewMerchantRepository(db)
	_v2MerchantUseCase := v2MerchantUsecase.NewMerchantUseCase(_v2MerchantRepo)
	_v2MerchantHandler := v2MerchantHandler.NewHandler(
		_v2MerchantUseCase,
		_v2MerchantRepo,
	)
	_v2MerchantHandler.RegisterRoutes(v2)

	// [WALLET]
	walletRepo := repository.NewWalletRepository(db)
	walletUseCase := usecase.NewWalletUseCase(walletRepo)
	walletController := controller.NewWalletController(walletUseCase, middlewareProvider.JWTAuth)
	// @Summary Wallet endpoints
	// @Description Endpoints for managing merchant wallets
	// @Tags wallet
	walletController.RegisterRoutes(v2)

	// [WEBHOOK]
	_callbackRepo := webhookRepo.NewCallbackRepository(db)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}

	// [Socialpay API]
	// Initialize payment processors
	telebirrProc := telebirrProcessor.NewProcessor(telebirrProcessor.ProcessorConfig{
		SecurityCredential: os.Getenv("TELEBIRR_SECURITY_CREDENTIAL"),
		Password:           os.Getenv("TELEBIRR_PASSWORD"),
		IsTestMode:         os.Getenv("APP_ENV") != "production",
		ShortCode:          os.Getenv("TELEBIRR_SHORT_CODE"),
		IdentityID:         os.Getenv("TELEBIRR_IDENTITY_ID"),
		BaseURL:            os.Getenv("TELEBIRR_BASE_URL"),
		CallbackURL:        "https://api.Socialpay.co/api/v1/settle/telebirr",
	})

	// Initialize CBE processor
	cbeProc := cbeProcessor.NewProcessor(cbeProcessor.ProcessorConfig{
		MerchantID:    os.Getenv("CBE_MERCHANT_ID"),
		MerchantKey:   os.Getenv("CBE_MERCHANT_KEY"),
		TerminalID:    os.Getenv("CBE_TERMINAL_ID"),
		CredentialKey: os.Getenv("CBE_CREDENTIAL_KEY"),
		IsTestMode:    os.Getenv("APP_ENV") != "production",
		BaseURL:       os.Getenv("CBE_BASE_URL"),
		CallbackURL:   "https://api.Socialpay.co/api/v1/settle/cbe",
	})

	// Initialize Awash processor
	awashProc := awashProcessor.NewProcessor(awashProcessor.ProcessorConfig{
		MerchantID:         os.Getenv("AWASH_TEST_MERCHANT_CODE"), // code
		CredentialKey:      os.Getenv("AWASH_TEST_PASSWORD"),
		IsTestMode:         true,
		BaseURL:            os.Getenv("AWASH_TEST_BASE_URL"),
		CallbackURL:        os.Getenv("AWASH_TEST_CALLBACK_URL"),
		MerchantTillNumber: os.Getenv("AWASH_TEST_TIN_NUMBER"),
	})

	// Initialize Cybersource processor
	cybersourceProc := cybersourceProcessor.NewProcessor(cybersourceProcessor.ProcessorConfig{
		AccessKey:  os.Getenv("CYBERSOURCE_ACCESS_KEY"),
		ProfileID:  os.Getenv("CYBERSOURCE_PROFILE_ID"),
		SecretKey:  os.Getenv("CYBERSOURCE_SECRET_KEY"),
		IsTestMode: os.Getenv("APP_ENV") != "production",
	})

	// Initializing EthSwitch Processor

	ethSwitchProc := etswitch.NewEtSwitchProcessor(etswitch.ProcessorConfig{
		UserName:   os.Getenv("ETHSWITCH_USERNAME"),
		Credential: os.Getenv("ETHSWITCH_PASSWORD"),
		BaseURL:    os.Getenv("ETHSWITCH_BASE_URL"),
		IsTestMode: true,
	})

	// Initialize M-PESA processor
	mpesaProc := mpesaProcessor.NewProcessor(mpesaProcessor.ProcessorConfig{
		Username:    os.Getenv("SAFARICOM_USERNAME"),
		Password:    os.Getenv("SAFARICOM_PASSWORD"),
		IsTestMode:  os.Getenv("APP_ENV") != "production",
		BaseURL:     os.Getenv("MPESA_BASE_URL"),
		CallbackURL: "https://api.Socialpay.co/api/v1/settle/mpesa",
	})

	// Create processors map for settlement handler
	processors := map[txEntity.TransactionMedium]payment.Processor{
		txEntity.TELEBIRR:    telebirrProc,
		txEntity.CBE:         cbeProc,
		txEntity.CYBERSOURCE: cybersourceProc,
		txEntity.MPESA:       mpesaProc,
		txEntity.ETHSWITCH:   ethSwitchProc,
	}

	// Initialize payment service with all processors as variadic arguments
	_paymentService := SocialpayUsecase.NewPaymentService(
		telebirrProc,
		cbeProc,
		cybersourceProc,
		mpesaProc,
		awashProc,
		ethSwitchProc,
	)

	_tipService := tipService.NewTipProcessingService(
		_transactionRepo,
		walletUseCase,
		_paymentService,
	)

	// Initialize notification service with merchant repository
	log.Printf("Initializing notification service...")
	var _transactionNotifier *notificationUsecase.TransactionNotifier = notifications.NewTransactionNotifier(
		_v2MerchantRepo,
		log,
	)

	_webhookUseCase := webhookUsecase.NewWebhookUseCase(
		cfg,
		_transactionRepo,
		_callbackRepo,
		walletUseCase,
		_tipService,
		_transactionNotifier,
	)
	_webhookController := webhookController.NewWebhookController(
		_webhookUseCase,
		middlewareProvider.JWTAuth,
	)
	_webhookController.RegisterRoutes(v2)

	// Initialize settlement handler
	settlementHandler := paymentController.NewSettlementHandler(processors, _webhookUseCase)

	// Register settlement routes (no middleware)
	settlementHandler.RegisterRoutes(v2)

	// [QR]
	_qrRepo := qrRepo.NewQRRepository(db)
	_qrUseCase := qrUsecase.NewQRUseCase(
		_qrRepo,
		_transactionRepo,
		_transactionUseCase,
		_paymentService,
		walletUseCase,
	)
	_qrHandler := qrHandler.NewHandler(_qrUseCase, middlewareProvider.JWTAuth)
	_qrHandler.RegisterRouter(v2)

	_SocialpayAPIUseCase := SocialpayUsecase.NewPaymentUseCase(SocialpayUsecase.UseCaseConfig{
		TransactionRepo:    _transactionRepo,
		HostedPaymentRepo:  _hostedPaymentRepo,
		TransactionUseCase: _transactionUseCase,
		PaymentService:     _paymentService,
		WalletUseCase:      walletUseCase,
		MerchantUseCase:    _v2MerchantUseCase,
	})

	_SocialpayAPIHandler := SocialpayController.NewHandler(
		_SocialpayAPIUseCase,
		middlewareProvider.APIKey,
		_v2MerchantRepo,
		_qrUseCase,
	)
	_SocialpayAPIHandler.RegisterRoutes(v2)
	_SocialpayAPIHandler.RegisterQRRoutes(v2)

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Create a context that can be cancelled
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start webhook consumer in a goroutine
	consumerDone := make(chan struct{})
	go func() {
		defer close(consumerDone)
		log.Printf("Starting webhook consumer")
		_webhookConsumer := webhookConsumer.NewWebhookDispatcherWorker(
			cfg,
			db,
			_webhookUseCase,
			_transactionUseCase,
			_hostedPaymentRepo,
		)
		_webhookConsumer.Start(ctx)
	}()

	// Start webhook sender worker in a goroutine
	senderDone := make(chan struct{})
	go func() {
		defer close(senderDone)
		log.Printf("Starting webhook sender worker")
		_webhookSender := webhookConsumer.NewWebhookSenderWorker(
			cfg,
			_webhookUseCase,
		)
		_webhookSender.Start(ctx)
	}()

	// Start server in a goroutine
	serverDone := make(chan struct{})
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	go func() {
		defer close(serverDone)
		log.Printf("Starting V2 server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal
	<-quit
	log.Println("Shutting down server...")

	// Cancel context to stop consumer and producer
	cancel()

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	// Stop the webhook producer
	if _webhookUseCase, ok := _webhookUseCase.(*webhookUsecase.WebhookUseCaseImpl); ok {
		if _webhookUseCase.GetProducer() != nil {
			_webhookUseCase.GetProducer().Stop()
		}
		if _webhookUseCase.GetSendProducer() != nil {
			_webhookUseCase.GetSendProducer().Stop()
		}
	}

	// Wait for both goroutines to finish
	<-consumerDone
	<-senderDone
	<-serverDone
	log.Println("Server exited properly")
}
