package main

import (
	"context"
	"log"
	"os"
	"time"

	// HTTP Server
	"github.com/socialpay/socialpay/src/pkg/auth/infra/network/http"
	"github.com/socialpay/socialpay/src/pkg/bank/adapter/controller"

	// Shared Database
	sharedDB "github.com/socialpay/socialpay/src/pkg/shared/database"
	// SMS
	"github.com/socialpay/socialpay/src/pkg/auth/adapter/gateway/sms"
	// [AUTH]
	"github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
	auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/rest"
	authRepo "github.com/socialpay/socialpay/src/pkg/auth/adapter/gateway/repo/psql"
	authUsecase "github.com/socialpay/socialpay/src/pkg/auth/usecase"

	// [HELP]
	help "github.com/socialpay/socialpay/src/pkg/help/adapter/controller/rest"
	helpRepo "github.com/socialpay/socialpay/src/pkg/help/adapter/gateway/repo/psql"
	helpUsecase "github.com/socialpay/socialpay/src/pkg/help/usecase"

	// [ORG]
	org "github.com/socialpay/socialpay/src/pkg/org/adapter/controller/rest"

	orgRepo "github.com/socialpay/socialpay/src/pkg/org/adapter/gateway/repo/psql"
	"github.com/socialpay/socialpay/src/pkg/org/adapter/gateway/tin_checker/etrade"
	orgUsecase "github.com/socialpay/socialpay/src/pkg/org/usecase"

	// [ACCOUNT]
	acc "github.com/socialpay/socialpay/src/pkg/account/adapter/controller/rest"
	accRepo "github.com/socialpay/socialpay/src/pkg/account/adapter/gateway/repo/psql"
	accUsecase "github.com/socialpay/socialpay/src/pkg/account/usecase"

	// [ACCOUNT]
	gateway "github.com/socialpay/socialpay/src/pkg/gateways/adapter/controller/rest"
	gatewayRepo "github.com/socialpay/socialpay/src/pkg/gateways/adapter/gateway/repo/psql"
	gatewayUsecase "github.com/socialpay/socialpay/src/pkg/gateways/usecase"

	//[Key]
	key "github.com/socialpay/socialpay/src/pkg/key/adapter/controller/rest"
	keyRepo "github.com/socialpay/socialpay/src/pkg/key/adapter/gateway/repo/psql"
	keyUsecase "github.com/socialpay/socialpay/src/pkg/key/usecase"

	//[Merchants]
	merchant "github.com/socialpay/socialpay/src/pkg/merchants/adapter/controller/rest"
	merchantRepo "github.com/socialpay/socialpay/src/pkg/merchants/adapter/gateway/repo/psql"
	merchantUsecase "github.com/socialpay/socialpay/src/pkg/merchants/usecase"

	// [ERP]
	erp "github.com/socialpay/socialpay/src/pkg/erp/adapter/controller/rest"
	erpRepo "github.com/socialpay/socialpay/src/pkg/erp/adapter/gateway/repo/psql"
	erpUsecase "github.com/socialpay/socialpay/src/pkg/erp/usecase"

	checkout "github.com/socialpay/socialpay/src/pkg/checkout/adapter/controller/rest"
	checkoutRepo "github.com/socialpay/socialpay/src/pkg/checkout/adapter/gateway/repo/psql"
	checkoutUsecase "github.com/socialpay/socialpay/src/pkg/checkout/service/basic"

	// [ACCESS CONTROL]
	access_control "github.com/socialpay/socialpay/src/pkg/access_control/adapter/controller/rest"
	accessRepo "github.com/socialpay/socialpay/src/pkg/access_control/adapter/gateway/repo/psql"
	accessUsecase "github.com/socialpay/socialpay/src/pkg/access_control/usecase"

	// [STORAGE]
	storage "github.com/socialpay/socialpay/src/pkg/storage/adapter/controller/rest"
	storageRepo "github.com/socialpay/socialpay/src/pkg/storage/adapter/gateway/repo/psql"
	storageUsecase "github.com/socialpay/socialpay/src/pkg/storage/usecase"
)

// @title           SocialPay API V1
// @version         1.0
// @description     SocialPay API V1 documentation - Legacy API with authentication, merchant management, ERP, and payment processing
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.socialpay.com/support
// @contact.email  support@socialpay.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      196.190.251.194:8082
// @BasePath  /api/v1/

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer token authentication for protected endpoints. Use format: "Bearer {token}"
// @bearerFormat JWT

// @tag.name Authentication
// @tag.description Authentication and authorization endpoints

// @tag.name Merchants
// @tag.description Merchant management and operations

// @tag.name ERP
// @tag.description Enterprise Resource Planning operations

// @tag.name Account
// @tag.description User account management

// @tag.name Key
// @tag.description API key management

// @tag.name Gateway
// @tag.description Payment gateway management

// @tag.name Help
// @tag.description Help and support endpoints

// @tag.name Organization
// @tag.description Organization management

// @tag.name Storage
// @tag.description File storage operations

// @tag.name Access Control
// @tag.description Access control and permissions

// @tag.name Checkout
// @tag.description Checkout and payment processing

func main() {
	log := log.New(os.Stdout, "[SOCIALPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Output Adapters]
	// [DB] Postgres - Use shared connection
	db, err := sharedDB.GetSharedConnection()
	if err != nil {
		log.Fatal("Failed to get shared database connection:", err.Error())
	}

	log.Println("Using shared database connection for v1 API")

	// Test the database connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// [Input Adapters]
	// [SERVER] HTTP
	s := http.New(log)

	// [Module Adapters]
	_merchantRepo, err := merchantRepo.NewPsqlRepo(log, db)

	// [AUTH]
	_authRepo, err := authRepo.NewPsqlRepo(log, db)
	if err != nil {
		log.Println(err)
	} else {
		_authSMS := sms.New(log)
		_authInteractor := authUsecase.New(log, _authRepo, _authSMS)

		auth.New(log, s.ServeMux, _authInteractor, _authRepo, _authSMS, _merchantRepo)
	}

	// [KEY]
	_keyRepo, err := keyRepo.NewPsqlRepo(log, db)
	if err != nil {
		log.Println(err)
	} else {
		_keyInteractor := keyUsecase.New(log, _keyRepo, sms.New(log))
		key.New(log, s.ServeMux, &_keyInteractor, _keyRepo, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

	// [MERCHANTS]
	if err != nil {
		log.Println(err)
	} else {
		_merchantInteractor := merchantUsecase.New(log, _merchantRepo, sms.New(log))
		merchant.New(log, s.ServeMux, &_merchantInteractor, _merchantRepo, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

	// Awash bank

	controller.NewAwashTest(log, s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))

	// [HELP]
	helpRepo, err := helpRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		help.New(log, helpUsecase.New(log, helpRepo), s.ServeMux)
	}

	// [ORG]
	_orgRepo, err := orgRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		org.New(log, orgUsecase.New(log, _orgRepo, etrade.New(log)), s.ServeMux)
	}

	// [STORAGE]
	_storageRepo, err := storageRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		storage.New(log, s.ServeMux, storageUsecase.New(log, _storageRepo))
	}

	// [ACCOUNT]
	_accRepo, err := accRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		acc.New(log, accUsecase.New(log, _accRepo, sms.New(log)), s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

	// [GATEWAY]
	_gatewayRepo, err := gatewayRepo.NewPsqlRepo(log, db)
	if err != nil {
		log.Println(err)
	} else {
		gatewayInteractor := gatewayUsecase.New(log, _gatewayRepo)
		gateway.New(log, s.ServeMux, gatewayInteractor)
	}

	// ERP
	_erpRepo, err := erpRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		erp.New(log, erpUsecase.New(log, _erpRepo), s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

	// Checkout
	_checkoutRepo, err := checkoutRepo.New(log, db)
	if err != nil {
		log.Printf("[CHECKOUT] checkout not initlized %v", err)
	} else {
		checkout.New(log, s.ServeMux, checkoutUsecase.New(log, _checkoutRepo, accUsecase.New(log, _accRepo, sms.New(log))))
		log.Print("[CHECKOUT] checkout initlized")
	}

	// ACCESS CONTROL
	_accessRepo, err := accessRepo.New(log, db)
	if err != nil {
		log.Println(err)
	} else {
		access_control.New(log, accessUsecase.New(log, _accessRepo), s.ServeMux, procedure.New(log, authUsecase.New(log, _authRepo, sms.New(log))))
	}

	// Use the existing Server.Serve() method which already implements graceful shutdown
	log.Println("Starting server on :8004 with graceful shutdown")
	s.Serve()

	log.Println("Server exited properly")
}
