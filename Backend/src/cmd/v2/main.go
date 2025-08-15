package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	// Gin
	"github.com/gin-gonic/gin"
	"github.com/socialpay/socialpay/src/pkg/env"

	// Repositories
	transactionRepo "github.com/socialpay/socialpay/src/pkg/transaction/core/repository"

	// [SocialPay API]
	"github.com/socialpay/socialpay/src/pkg/shared/payment"
	paymentController "github.com/socialpay/socialpay/src/pkg/shared/payment/controller/gin"
	"github.com/socialpay/socialpay/src/pkg/shared/utils"
	socialpayController "github.com/socialpay/socialpay/src/pkg/socialpayapi/adapter/controller/gin"
	socialpayUsecase "github.com/socialpay/socialpay/src/pkg/socialpayapi/usecase"

	// Payment Processors
	awashProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/awash"
	cbeProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/cbe"
	cybersourceProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/cybersource"
	etswitch "github.com/socialpay/socialpay/src/pkg/shared/payment/etswitch"
	kachaProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/kacha"
	mpesaProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/mpesa"
	telebirrProcessor "github.com/socialpay/socialpay/src/pkg/shared/payment/telebirr"

	// Transaction Types
	txEntity "github.com/socialpay/socialpay/src/pkg/transaction/core/entity"

	// Swagger
	_ "github.com/socialpay/socialpay/docs"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	// [API KEY]
	apikeyHandler "github.com/socialpay/socialpay/src/pkg/apikey_mgmt/adapter/controller/gin"
	apikeyRepo "github.com/socialpay/socialpay/src/pkg/apikey_mgmt/adapter/gateway/repo/sqlc"
	apikeyUsecase "github.com/socialpay/socialpay/src/pkg/apikey_mgmt/usecase"

	// [V2 MERCHANT]
	v2MerchantHandler "github.com/socialpay/socialpay/src/pkg/v2_merchant/adapter/controller/gin"
	v2MerchantRepo "github.com/socialpay/socialpay/src/pkg/v2_merchant/adapter/gateway/repo/sqlc"

	// v2MerchantSeeder "github.com/socialpay/socialpay/src/pkg/v2_merchant/seeder"
	v2MerchantUsecase "github.com/socialpay/socialpay/src/pkg/v2_merchant/usecase"

	// [TRANSACTION]
	transactionHandler "github.com/socialpay/socialpay/src/pkg/transaction/adapter/controller/gin"
	transactionUsecase "github.com/socialpay/socialpay/src/pkg/transaction/usecase"

	// [MIDDLEWARE]
	"github.com/socialpay/socialpay/src/pkg/shared/middleware"

	// [OLD AUTH] - Keep for backward compatibility during transition
	authRepo "github.com/socialpay/socialpay/src/pkg/auth/adapter/gateway/repo/psql"
	authUsecase "github.com/socialpay/socialpay/src/pkg/auth/usecase"

	// [SHARED DATABASE]
	sharedDB "github.com/socialpay/socialpay/src/pkg/shared/database"

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
	//"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// [WALLET]
	walletController "github.com/socialpay/socialpay/src/pkg/wallet/adapter/controller"
	walletRepo "github.com/socialpay/socialpay/src/pkg/wallet/adapter/gateway/repository"
	walletUsecase "github.com/socialpay/socialpay/src/pkg/wallet/usecase"

	// [COMMISSION]
	commissionController "github.com/socialpay/socialpay/src/pkg/commission/adapter/controller"
	commissionRepo "github.com/socialpay/socialpay/src/pkg/commission/adapter/gateway/repository"
	commission_usecase "github.com/socialpay/socialpay/src/pkg/commission/usecase"

	// [ERP_V2]
	erpHandler "github.com/socialpay/socialpay/src/pkg/erp_v2/adapter/controller/gin"
	erpRepo "github.com/socialpay/socialpay/src/pkg/erp_v2/adapter/gateway/repo/sqlc"
	erpUsecase "github.com/socialpay/socialpay/src/pkg/erp_v2/usecase"

	// [CONFIG]
	config "github.com/socialpay/socialpay/src/pkg/config"

	// [NOTIFICATIONS]
	notifications "github.com/socialpay/socialpay/src/pkg/notifications"
	notificationUsecase "github.com/socialpay/socialpay/src/pkg/notifications/usecase"

	// Auth v2 imports
	authv2Handler "github.com/socialpay/socialpay/src/pkg/authv2/adapter/controller/gin"
	"github.com/socialpay/socialpay/src/pkg/authv2/adapter/repository/postgres"
	authv2Seeder "github.com/socialpay/socialpay/src/pkg/authv2/seeder"
	authv2Service "github.com/socialpay/socialpay/src/pkg/authv2/service"

	// Team Member imports
	teamMemberHandler "github.com/socialpay/socialpay/src/pkg/team_member/adapter/controller/gin"
	teamMemberRepository "github.com/socialpay/socialpay/src/pkg/team_member/adapter/repository/postgres"
	teamMemberService "github.com/socialpay/socialpay/src/pkg/team_member/service"

	// RBAC imports
	rbacHandler "github.com/socialpay/socialpay/src/pkg/rbac/adapter/controller/gin"
	rbacRepository "github.com/socialpay/socialpay/src/pkg/rbac/adapter/repository/postgres"
	rbacService "github.com/socialpay/socialpay/src/pkg/rbac/service"

	// File service imports
	fileHandler "github.com/socialpay/socialpay/src/pkg/file/adapter/controller/gin"
	fileService "github.com/socialpay/socialpay/src/pkg/file/service"

	// Logging
	"github.com/socialpay/socialpay/src/pkg/shared/logging"

	// IP Whitelisting
	ipWhitelistHandler "github.com/socialpay/socialpay/src/pkg/ip_whitelist/adapter/controller"
	ipWhitelistRepo "github.com/socialpay/socialpay/src/pkg/ip_whitelist/adapter/gateway/repository"
	ipWhitelistUsecase "github.com/socialpay/socialpay/src/pkg/ip_whitelist/usecase"
)

// @title           SocialPay API V2
// @version         2.0
// @description     SocialPay API V2 documentation - Complete payment gateway API with authentication, merchant management, transactions, QR payments, and ERP functionality
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.socialpay.com/support
// @contact.email  support@socialpay.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      196.190.251.194:8082:8080
// @BasePath  /api/v2/

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

// @tag.name Authentication
// @tag.description Authentication and authorization endpoints

// @tag.name Merchants
// @tag.description Merchant management and operations

// @tag.name Payments
// @tag.description Payment processing and transaction management

// @tag.name QR
// @tag.description QR code payment operations

// @tag.name Transactions
// @tag.description Transaction history and analytics

// @tag.name API Keys
// @tag.description API key management for merchants

// @tag.name ERP
// @tag.description Enterprise Resource Planning operations

// @tag.name Webhooks
// @tag.description Webhook management and callbacks

// @tag.name Wallet
// @tag.description Wallet operations and management

// @tag.name Commission
// @tag.description Commission management and calculations

// @tag.name IP Whitelist
// @tag.description IP whitelist management for security

// @tag.name Team Management
// @tag.description Team member and role management

// @tag.name RBAC
// @tag.description Role-Based Access Control operations

// @security BearerAuth

func main() {
	// Check environment variables first
	env.LoadEnv()
	env.CheckEnv()

	log := log.New(os.Stdout, "[SOCIALPAY-V2]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Output Adapters]
	// [DB] Postgres - Use shared connection for all modules
	db, err := sharedDB.GetSharedConnection()
	if err != nil {
		log.Fatal("Failed to get shared database connection:", err.Error())
	}

	log.Println("Shared database connection established for all modules")

	// Run database migrations
	log.Println("Running database migrations...")

	// Get current working directory
	workDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get working directory:", err)
	}
	log.Printf("Current working directory: %s", workDir)

	// Construct absolute path for migrations
	migrationsPath := fmt.Sprintf("file://%s/db/migrations", workDir)
	log.Printf("Migrations path: %s", migrationsPath)

	// Construct database URL
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)
	log.Printf("Database URL: postgres://%s:***@%s:%s/%s?sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
		os.Getenv("SSL_MODE"),
	)
	fmt.Println("Database URL:", dbURL)

	/* // Check if migrations directory exists
	if _, err := os.Stat(fmt.Sprintf("%s/db/migrations", workDir)); os.IsNotExist(err) {
		log.Fatal("Migrations directory does not exist at:", fmt.Sprintf("%s/db/migrations", workDir))
	}

	// Initialize migrate instance
	log.Println("Initializing migrate...")
	m, err := migrate.New(
		migrationsPath,
		dbURL,
	)
	if err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to initialize migrations:", err)
	}

	// Run migrations
	log.Println("Running migrations...")
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("Failed to run migrations:", err)
	}

	log.Println("Database migrations completed successfully") */

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Initialize Auth v2 components
	log.Println("Initializing Auth v2 system...")

	// Initialize Auth v2 repository using existing db connection
	authv2RepoInstance := postgres.NewAuthRepository(db, log)

	// Initialize Auth v2 service
	authv2ServiceInstance := authv2Service.NewAuthService(
		authv2RepoInstance,
		os.Getenv("JWT_SECRET"), // JWT secret from environment
		log,                     // standard logger
	)

	// Initialize Auth v2 handler
	authv2HandlerInstance := authv2Handler.NewAuthHandler(authv2ServiceInstance)

	// Initialize Auth v2 seeder
	authv2SeederInstance := authv2Seeder.NewAuthSeeder(authv2ServiceInstance, authv2RepoInstance, db)

	// Seed Auth v2 data
	if err := authv2SeederInstance.SeedAll(context.Background()); err != nil {
		log.Printf("Warning: Failed to seed Auth v2 data: %v", err)
	}

	// Recalculate wallet amounts after migrations
	log.Println("Starting wallet amount recalculation...")
	// if err := recalculateWalletBalances(db); err != nil {
	// 	log.Printf("Warning: Failed to recalculate wallet amounts: %v", err)
	// } else {
	// 	log.Println("Wallet amount recalculation completed successfully")
	// }

	// Keep old auth for backward compatibility during transition
	_authRepo, err := authRepo.NewPsqlRepo(log, db)
	if err != nil {
		log.Fatal("Failed to initialize auth repository:", err)
	}
	_authUseCase := authUsecase.New(log, _authRepo, sms.New(log))

	// Beign [API KEY]
	_apikeyRepo := apikeyRepo.NewRepository(db)
	_apikeyUseCase := apikeyUsecase.NewAPIKeyUseCase(_apikeyRepo)

	_ipWhitelistRepo := ipWhitelistRepo.NewIPWhitelistRepository(db)
	_ipWhitelistUsecase := ipWhitelistUsecase.NewIPWhitelistUseCase(_ipWhitelistRepo)

	// Initialize middleware provider
	middlewareProvider := middleware.NewMiddlewareProvider(_authUseCase, _apikeyUseCase, authv2ServiceInstance, _ipWhitelistUsecase)

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
	log.Printf("Created v2 group with base path: %s", v2.BasePath())

	v2.Static("/static", "./public")

	// Register Auth v2 routes
	authv2HandlerInstance.RegisterRoutes(v2)
	authv2HandlerInstance.RegisterAdminRoutes(v2)

	// Initialize rbac repository with existing db connection
	rbacRepositoryInstance := rbacRepository.NewRBACRepository(db, log)

	// Initialize rbac service
	rbacServiceInstance := rbacService.NewRBACService(
		rbacRepositoryInstance,
		os.Getenv("JWT_SECRET"),
		log,
	)

	// Initialize rbac
	rbacHandlerInstance := rbacHandler.NewRBACHandler(authv2ServiceInstance, rbacServiceInstance, middlewareProvider.RBAC)

	// Register rbac routes
	rbacHandlerInstance.RegisterRoutes(v2)

	// Initialize Team member management
	log.Println("Initializing team member management system...")

	// Initialize Team memeber management repository using existing db connection
	teamMemberRepoInstance := teamMemberRepository.NewTeamMemberRepository(db, log, rbacRepositoryInstance)

	// Initialize Team memeber management service
	teamMemberServiceInstance := teamMemberService.NewTeamMemberService(
		authv2RepoInstance,
		teamMemberRepoInstance,
		os.Getenv("JWT_SECRET"), // JWT secret from environment
		log,                     // standard logger
	)

	// Initialize Team memeber management handler
	teamMemberHandlerInstance := teamMemberHandler.NewTeamMemberHandler(authv2ServiceInstance, teamMemberServiceInstance, middlewareProvider.RBAC)

	// Register Team memeber management routes
	teamMemberHandlerInstance.RegisterRoutes(v2)

	// [ERP_V2]
	_erpRepo := erpRepo.NewSQLCRepository(db)
	_erpUsecase := erpUsecase.NewERPUseCase(_erpRepo)
	_erpHandler := erpHandler.NewERPHandler(authv2ServiceInstance, _erpUsecase, middlewareProvider.RBAC)
	_erpHandler.RegisterRoutes(v2)
	// Initialize cloudinary
	_cloudinaryInstance, err := utils.InitCloudinary()
	if err != nil {
		log.Fatal("Failed to initialize cloudinary:", err)
	}

	// Initialize file service
	_fileServiceInstance := fileService.NewFileService(_cloudinaryInstance, log)

	// Initialize file service handler
	_fileServiceHandler := fileHandler.NewFileHandler(authv2ServiceInstance, _fileServiceInstance)

	// Register File service routes
	_fileServiceHandler.RegisterRoutes(v2)

	// [TRANSACTION]
	_transactionRepo := transactionRepo.NewTransactionRepository(db)
	_hostedPaymentRepo := transactionRepo.NewHostedPaymentRepository(db)
	_transactionUseCase := transactionUsecase.NewTransactionUsecase(
		_transactionRepo,
	)
	// Handler at the bottom

	// [API KEY]
	_apikeyHandler := apikeyHandler.NewHandler(
		_apikeyUseCase,
		_apikeyRepo,
		middlewareProvider.JWTAuth,
		middlewareProvider.Public,
		middlewareProvider.RBAC,
	)
	_apikeyHandler.RegisterRoutes(v2)

	// Merchant seeder
	// merchantSeeder := v2MerchantSeeder.NewMerchantSeeder(db, authv2RepoInstance)
	// merchantSeeder.SeedMerchant(context.Background(), 10)

	// [V2 MERCHANT]
	_v2MerchantRepo := v2MerchantRepo.NewMerchantRepository(db)
	_v2MerchantUseCase := v2MerchantUsecase.NewMerchantUseCase(authv2ServiceInstance, _v2MerchantRepo)
	_v2MerchantHandler := v2MerchantHandler.NewHandler(
		authv2ServiceInstance,
		_v2MerchantUseCase,
		middlewareProvider.RBAC,
		_v2MerchantRepo,
	)
	_v2MerchantHandler.RegisterRoutes(v2)

	// IP Whitelist
	_ipWhitelistHandler := ipWhitelistHandler.NewIPWhitelistController(authv2ServiceInstance, middlewareProvider.RBAC, _ipWhitelistUsecase)
	_ipWhitelistHandler.RegisterRoutes(v2)

	// [ WALLET]
	fmt.Println("[WALLET] Initializing Wallet Repository")
	_walletRepo := walletRepo.NewWalletRepository(db)
	_walletUseCase := walletUsecase.NewMerchantWalletUsecase(_walletRepo, logging.NewStdLogger("[WALLET]"))
	_walletController := walletController.NewWalletController(_walletUseCase, middlewareProvider.JWTAuth, middlewareProvider.RBAC)
	_walletController.RegisterRoutes(v2)

	// [ADMIN WALLET CONTROLLER]
	_adminwalleRepo := walletRepo.NewWalletRepository(db)
	_adminWalletUseCase := walletUsecase.NewAdminWalletUsecase(_adminwalleRepo)
	_adminWalletController := walletController.NewAdminWalletController(_adminWalletUseCase, middlewareProvider)
	_adminWalletController.RegisterRoutes(v2)

	// [COMMISSION]
	_commissionRepo := commissionRepo.NewCommissionRepository(db)
	_commissionUseCase := commission_usecase.NewCommissionUseCase(_commissionRepo)
	_commissionController := commissionController.NewCommissionController(
		_commissionUseCase,
		middlewareProvider,
	)
	log.Printf("Registering commission controller routes under v2 group: %s", v2.BasePath())
	_commissionController.RegisterRoutes(v2)

	// [SocialPay API]
	// Initialize payment processors
	telebirrProc := telebirrProcessor.NewProcessor(telebirrProcessor.ProcessorConfig{
		SecurityCredential: os.Getenv("TELEBIRR_SECURITY_CREDENTIAL"),
		Password:           os.Getenv("TELEBIRR_PASSWORD"),
		IsTestMode:         os.Getenv("APP_ENV") != "production",
		ShortCode:          os.Getenv("TELEBIRR_SHORT_CODE"),
		IdentityID:         os.Getenv("TELEBIRR_IDENTITY_ID"),
		BaseURL:            os.Getenv("TELEBIRR_BASE_URL"),
		CallbackURL:        "https://api.socialpay.co/api/v1/settle/telebirr",
	})

	// Initialize CBE processor
	cbeProc := cbeProcessor.NewProcessor(cbeProcessor.ProcessorConfig{
		MerchantID:    os.Getenv("CBE_MERCHANT_ID"),
		MerchantKey:   os.Getenv("CBE_MERCHANT_KEY"),
		TerminalID:    os.Getenv("CBE_TERMINAL_ID"),
		CredentialKey: os.Getenv("CBE_CREDENTIAL_KEY"),
		IsTestMode:    os.Getenv("APP_ENV") != "production",
		BaseURL:       os.Getenv("CBE_BASE_URL"),
		CallbackURL:   "https://api.socialpay.co/api/v1/settle/cbe",
	})

	// Initialize Awash processor
	awashProc := awashProcessor.NewProcessor(awashProcessor.ProcessorConfig{
		MerchantID:         os.Getenv("AWASH_TEST_MERCHANT_CODE"), // code
		CredentialKey:      os.Getenv("AWASH_TEST_PASSWORD"),
		IsTestMode:         true,
		BaseURL:            os.Getenv("AWASH_TEST_BASE_URL"),
		CallbackURL:        os.Getenv("AWASH_TEST_CALLBACK_URL"),
		MerchantTillNumber: os.Getenv("AWASH_TEST_TIN_NUMBER"),
		TxnRepository:      _transactionRepo,
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
		RetunUrl:   os.Getenv("APP_CHECKOUT_URL"),
		IsTestMode: true,
	})

	// Initialize M-PESA processor
	mpesaProc := mpesaProcessor.NewProcessor(mpesaProcessor.ProcessorConfig{
		Username:    os.Getenv("SAFARICOM_USERNAME"),
		Password:    os.Getenv("SAFARICOM_PASSWORD"),
		IsTestMode:  os.Getenv("APP_ENV") != "production",
		BaseURL:     os.Getenv("MPESA_BASE_URL"),
		CallbackURL: "https://api.socialpay.co/api/v1/settle/mpesa",
	})

	// Initialize Kacha processor
	kachaProc := kachaProcessor.NewProcessor(kachaProcessor.ProcessorConfig{
		IsTestMode:  os.Getenv("APP_ENV") != "production",
		BaseURL:     os.Getenv("KACHA_BASE_URL"),
		CallbackURL: "https://api.socialpay.co/api/v1/settle/std",
	})

	// Create processors map for settlement handler
	processors := map[txEntity.TransactionMedium]payment.Processor{
		txEntity.TELEBIRR:    telebirrProc,
		txEntity.CBE:         cbeProc,
		txEntity.CYBERSOURCE: cybersourceProc,
		txEntity.MPESA:       mpesaProc,
		txEntity.ETHSWITCH:   ethSwitchProc,
		txEntity.KACHA:       kachaProc,
	}

	// Initialize payment service with all processors as variadic arguments
	_paymentService := socialpayUsecase.NewPaymentService(
		telebirrProc,
		cbeProc,
		cybersourceProc,
		mpesaProc,
		awashProc,
		ethSwitchProc,
		kachaProc,
	)

	_tipService := socialpayUsecase.NewTipProcessingService(
		_transactionRepo,
		_walletUseCase,
		_paymentService,
	)

	// Initialize notification service with merchant repository
	log.Printf("Initializing notification service...")
	var _transactionNotifier *notificationUsecase.TransactionNotifier = notifications.NewTransactionNotifier(
		_v2MerchantRepo,
		log,
	)

	// [WEBHOOK]
	_callbackRepo := webhookRepo.NewCallbackRepository(db)
	_cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config: " + err.Error())
	}
	_webhookUseCase := webhookUsecase.NewWebhookUseCase(
		_cfg,
		_transactionRepo,
		_callbackRepo,
		_walletUseCase,
		_adminWalletUseCase,
		_commissionUseCase,
		_tipService,
		_transactionNotifier,
	)
	_webhookController := webhookController.NewWebhookController(
		_webhookUseCase,
		middlewareProvider.JWTAuth,
		middlewareProvider.RBAC,
	)
	_webhookController.RegisterRoutes(v2)

	// Because of the transactionHandler is depending on the webhookUseCase, we need to initialize it here

	_transactionHandler := transactionHandler.NewTransactionHistoryHandler(
		_transactionUseCase,
		middlewareProvider,
		_webhookUseCase,
	)

	// Register transaction routes
	_transactionHandler.RegisterRoutes(v2)
	_transactionHandler.RegisterAdminRoutes(v2)

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
		_walletUseCase,
		_commissionUseCase,
	)
	_qrHandler := qrHandler.NewHandler(_qrUseCase, middlewareProvider.JWTAuth, middlewareProvider.RBAC)
	_qrHandler.RegisterRouter(v2)

	_socialpayAPIUseCase := socialpayUsecase.NewPaymentUseCase(socialpayUsecase.UseCaseConfig{
		TransactionRepo:    _transactionRepo,
		HostedPaymentRepo:  _hostedPaymentRepo,
		TransactionUseCase: _transactionUseCase,
		PaymentService:     _paymentService,
		WalletUseCase:      _walletUseCase,
		MerchantUseCase:    _v2MerchantUseCase,
		CommissionUseCase:  _commissionUseCase,
	})

	_socialpayAPIHandler := socialpayController.NewHandler(
		_socialpayAPIUseCase,
		middlewareProvider.APIKey,
		_v2MerchantRepo,
		*middlewareProvider.IPChecker,
		_qrUseCase,
		_webhookUseCase,
	)
	_socialpayAPIHandler.RegisterRoutes(v2)
	_socialpayAPIHandler.RegisterQRRoutes(v2)

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
			_cfg,
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
			_cfg,
			_webhookUseCase,
		)
		_webhookSender.Start(ctx)
	}()

	// Initialize and start cron service
	_transactionStatusChecker := socialpayUsecase.NewTransactionStatusChecker(
		_transactionRepo,
		_paymentService,
		_webhookUseCase,
	)

	_cronService := socialpayUsecase.NewCronService(_transactionStatusChecker, ctx)

	if err := _cronService.Start(); err != nil {
		log.Fatalf("Failed to start cron service: %v", err)
	}

	log.Printf("Cron service started with %d jobs", _cronService.GetJobCount())

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

	// Stop cron service
	_cronService.Stop()

	// Close shared database connection
	if err := sharedDB.CloseSharedConnection(); err != nil {
		log.Printf("Error closing shared database connection: %v", err)
	}

	// Wait for all goroutines to finish
	<-consumerDone
	<-senderDone
	<-serverDone
	log.Println("Server exited properly")
}

// recalculateWalletBalances recalculates wallet amounts based on transaction history
func recalculateWalletBalances(db *sql.DB) error {
	ctx := context.Background()

	// Execute the SQL script to recalculate amounts
	sqlScript := `
DO $$
DECLARE
    wallet_record RECORD;
    calculated_amount DECIMAL(20,2);
    current_amount DECIMAL(20,2);
    affected_rows INTEGER := 0;
    deposit_sum DECIMAL(20,2);
    withdrawal_sum DECIMAL(20,2);
    merchant_wallet_count INTEGER := 0;
    admin_wallet_count INTEGER := 0;
    total_transactions INTEGER := 0;
BEGIN
    RAISE NOTICE 'Starting wallet amount recalculation...';
    
    -- Debug: Check total transactions count
    SELECT COUNT(*) INTO total_transactions FROM public.transactions;
    RAISE NOTICE 'DEBUG: Total transactions in database: %', total_transactions;
    
    -- Debug: Check merchant wallet count
    SELECT COUNT(*) INTO merchant_wallet_count 
    FROM merchant.wallet w 
    WHERE w.wallet_type = 'merchant';
    RAISE NOTICE 'DEBUG: Found % merchant wallets', merchant_wallet_count;
    
    -- Debug: Check admin wallet count
    SELECT COUNT(*) INTO admin_wallet_count 
    FROM merchant.wallet w 
    WHERE w.wallet_type = 'super_admin';
    RAISE NOTICE 'DEBUG: Found % admin wallets', admin_wallet_count;
    
    -- Debug: Show all wallet types that exist
    FOR wallet_record IN 
        SELECT wallet_type, COUNT(*) as count
        FROM merchant.wallet 
        GROUP BY wallet_type
    LOOP
        RAISE NOTICE 'DEBUG: Wallet type "%" has % wallets', wallet_record.wallet_type, wallet_record.count;
    END LOOP;
    
    -- Recalculate merchant wallet amounts
    FOR wallet_record IN 
        SELECT w.id, w.user_id, w.merchant_id, w.amount as current_amount
        FROM merchant.wallet w 
        WHERE w.wallet_type = 'merchant'
    LOOP
        RAISE NOTICE 'DEBUG: Processing merchant wallet - ID: %, User: %, Merchant: %, Current Amount: %', 
            wallet_record.id, wallet_record.user_id, wallet_record.merchant_id, wallet_record.current_amount;
        
        -- Calculate deposits (positive merchant_net) from successful transactions
        SELECT COALESCE(SUM(merchant_net), 0)
        INTO deposit_sum
        FROM public.transactions 
        WHERE merchant_id = wallet_record.merchant_id 
        AND status = 'SUCCESS'
        AND type IN ('DEPOSIT')
        AND merchant_net IS NOT NULL
        AND merchant_net > 0;
        
        RAISE NOTICE 'DEBUG: Merchant % - Deposit sum: %', wallet_record.merchant_id, deposit_sum;
        
        -- Calculate withdrawals (positive merchant_net) from successful transactions
        SELECT COALESCE(SUM(ABS(merchant_net)), 0)
        INTO withdrawal_sum
        FROM public.transactions 
        WHERE merchant_id = wallet_record.merchant_id 
        AND status = 'SUCCESS'
        AND type = 'WITHDRAWAL'
        AND merchant_net IS NOT NULL;
        
        RAISE NOTICE 'DEBUG: Merchant % - Withdrawal sum: %', wallet_record.merchant_id, withdrawal_sum;
        
        -- Balance = DEPOSITS - WITHDRAWALS
        calculated_amount := deposit_sum - withdrawal_sum;
        
        RAISE NOTICE 'DEBUG: Merchant % - Calculated amount: % (Current: %)', 
            wallet_record.merchant_id, calculated_amount, wallet_record.current_amount;
        
                -- Update if amount is different
        IF wallet_record.current_amount != calculated_amount THEN
            RAISE NOTICE 'DEBUG: Updating merchant wallet % - Amount differs!', wallet_record.merchant_id;
            
            UPDATE merchant.wallet 
            SET amount = calculated_amount,
			locked_amount = 0,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = wallet_record.id;
            
            affected_rows := affected_rows + 1;
            
            RAISE NOTICE 'Updated merchant wallet % from % to % (Deposits: %, Withdrawals: %)', 
                wallet_record.merchant_id, wallet_record.current_amount, calculated_amount, deposit_sum, withdrawal_sum;
        ELSE
            RAISE NOTICE 'DEBUG: Merchant wallet % already has correct amount (%)', 
                wallet_record.merchant_id, wallet_record.current_amount;
        END IF;
    END LOOP;
    
        RAISE NOTICE 'DEBUG: Starting admin wallet recalculation...';
    
    -- Recalculate admin wallet amounts (single admin wallet)
    FOR wallet_record IN 
        SELECT w.id, w.user_id, w.amount as current_amount
        FROM merchant.wallet w 
        WHERE w.wallet_type = 'super_admin'
        LIMIT 1
    LOOP
        RAISE NOTICE 'DEBUG: Processing admin wallet - ID: %, User: %, Current Amount: %', 
            wallet_record.id, wallet_record.user_id, wallet_record.current_amount;
        
        -- Calculate total admin net (commission) from successful transactions
        SELECT COALESCE(SUM(admin_net), 0)
        INTO calculated_amount
        FROM public.transactions 
        WHERE status = 'SUCCESS'
        AND admin_net IS NOT NULL;
        
        RAISE NOTICE 'DEBUG: Admin wallet - Calculated amount: % (Current: %)', 
            calculated_amount, wallet_record.current_amount;
        
        -- Update if amount is different
        IF wallet_record.current_amount != calculated_amount THEN
            RAISE NOTICE 'DEBUG: Updating admin wallet - Amount differs!';
            
            UPDATE merchant.wallet 
            SET amount = calculated_amount,
			locked_amount = 0,
                updated_at = CURRENT_TIMESTAMP
            WHERE id = wallet_record.id;
            
            affected_rows := affected_rows + 1;
            
            RAISE NOTICE 'Updated admin wallet % from % to %', 
                wallet_record.user_id, wallet_record.current_amount, calculated_amount;
        ELSE
            RAISE NOTICE 'DEBUG: Admin wallet already has correct amount (%))', wallet_record.current_amount;
        END IF;
    END LOOP;
    
    -- Final debug summary
    RAISE NOTICE 'DEBUG: Summary - Total merchant wallets: %, Total admin wallets: %, Total transactions: %', 
        merchant_wallet_count, admin_wallet_count, total_transactions;
    RAISE NOTICE 'Wallet amount recalculation completed. Updated % wallets.', affected_rows;
END $$;
	`

	_, err := db.ExecContext(ctx, sqlScript)
	if err != nil {
		return fmt.Errorf("failed to execute wallet amount recalculation: %w", err)
	}

	return nil
}
