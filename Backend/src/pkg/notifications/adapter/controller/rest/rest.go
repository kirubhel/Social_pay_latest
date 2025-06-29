package rest

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
	"github.com/socialpay/socialpay/src/pkg/notifications/usecase"
)

// Controller struct
type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	auth       auth.Controller
	sm         *http.ServeMux
	sms        usecase.SMSSender
}

// New function initializes the Controller and sets up routes
func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux, auth auth.Controller, sms usecase.SMSSender) Controller {
	controller := Controller{log: log, interactor: interactor, auth: auth, sms: sms, sm: sm}

	// Handle routing
	sm.HandleFunc("/api/v1/send/transaction/sms", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.SendTransactionSMS(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return controller
}

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *Error      `json:"error,omitempty"`
}

// Define a generic response structure
type YimuluResponse struct {
	StatusCode    int    `json:"status_code,omitempty"`
	Message       string `json:"message,omitempty"`
	WalletBalance string `json:"wallet_balance,omitempty"`
	Error         *Error `json:"error,omitempty"`
}

func SendJSONResponse(w http.ResponseWriter, data Response, status int) {
	serData, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(serData)
}

/*


	// // Bank Accounts
	// sm.HandleFunc("/accounts/bank-accounts", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPost:
	// 		{
	// 			controller.GetAddBankAccount(w, r)
	// 		}
	// 	}
	// })

	// // Verify account
	// sm.HandleFunc("/accounts/bank-accounts/verify", func(w http.ResponseWriter, r *http.Request) {
	// 	switch r.Method {
	// 	case http.MethodPatch:
	// 		{
	// 			controller.GetVerifyBankAccount(w, r)
	// 		}
	// 	}
	// })



*/
