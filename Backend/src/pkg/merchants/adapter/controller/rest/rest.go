package rest

import (
	"log"
	"net/http"
	"time"

	auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
	"github.com/socialpay/socialpay/src/pkg/merchants/usecase"
)

type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	repo       usecase.Repository
	auth       auth.Controller
	sm         *http.ServeMux
	//sms        usecase.SMSSender
}

// Response
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type MerchantsResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

type Merchant struct {
	ID                 string    `json:"id"`                  // UUID
	UserID             string    `json:"user_id"`             // UUID
	LegalName          string    `json:"legal_name"`          // Varying character
	TradingName        string    `json:"trading_name"`        // Varying character
	BusinessRegNumber  string    `json:"business_reg_number"` // Varying character
	TaxIdentifier      string    `json:"tax_identifier"`      // Varying character
	IndustryType       string    `json:"industry_type"`       // Varying character
	BusinessType       string    `json:"business_type"`       // Varying character
	IsBettingClient    bool      `json:"is_betting_client"`   // Boolean
	LoyaltyCertificate string    `json:"loyalty_certificate"` // Varying character
	WebsiteURL         string    `json:"website_url"`         // Varying character
	EstablishedDate    time.Time `json:"established_date"`    // Date
	CreatedAt          time.Time `json:"created_at"`          // Timestamp
	UpdatedAt          time.Time `json:"updated_at"`          // Timestamp
	Status             string    `json:"status"`              // Varying character
}

func New(log *log.Logger, sm *http.ServeMux, interactor usecase.Interactor, repo usecase.Repository, auth auth.Controller) Controller {

	controller := Controller{log: log, interactor: interactor, repo: repo, auth: auth}

	sm.HandleFunc("/api/merchant/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.CreateKey(w, r)
			}
		}
	})
	sm.HandleFunc("/api/merchant-by-user", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetMerchantByUserID(w, r)
			}
		}
	})

	sm.HandleFunc("/api/fetch/all-merchants", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetMerchants(w, r)
			}
		}
	})

	sm.HandleFunc("/api/update/merchant/status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.UpdateMerchantStatus(w, r)
			}
		}
	})

	sm.HandleFunc("/api/merchant/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.UpdateFullMerchant(w, r)
			}
		}
	})

	sm.HandleFunc("/api/get/merchant/details", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetMerchantBusinessInformations(w, r)
			}
		}
	})

	sm.HandleFunc("/api/merchant/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.DeleteMerchant(w, r)
			}
		}
	})

	sm.HandleFunc("/api/merchant/additional-info", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.AddMerchantInfo(w, r)
			}
		}
	})

	sm.HandleFunc("/api/merchant/upload", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.UploadDocument(w, r)
			}
		}
	})

	// sm.HandleFunc("/key/enable", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		{
	// 			controller.Enable(w, r)
	// 		}
	// 	}
	// })

	// sm.HandleFunc("/key/disable", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		{
	// 			controller.Disable(w, r)
	// 		}
	// 	}
	// })

	controller.sm = sm
	return controller
}
