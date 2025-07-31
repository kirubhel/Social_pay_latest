package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/account/usecase"
	auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
)

// Controller struct
type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	auth       auth.Controller
	sm         *http.ServeMux
}

// New function initializes the Controller and sets up routes
func New(log *log.Logger, interactor usecase.Interactor, sm *http.ServeMux, auth auth.Controller) Controller {
	controller := Controller{log: log, interactor: interactor, auth: auth}

	// Handle routing
	// User Update
	sm.HandleFunc("/api/v1/user-update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UpdateUser(w, r)
		}
	})

	// Get User Profile
	sm.HandleFunc("/api/v1/user-profile", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUserProfile(w, r)
		}
	})

	// Send OTP
	sm.HandleFunc("/send-otp", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetSendOtp(w, r)
		}
	})

	// Finger Print
	sm.HandleFunc("/finger-print", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			controller.GetSetFingerPrint(w, r)
		}
	})

	// Generate Challenge
	sm.HandleFunc("/generate-challenge", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetGenerateChallenge(w, r)
		}
	})

	// Verify Signature
	sm.HandleFunc("/verify-signature", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetVerifySignatureHandler(w, r)
		}
	})

	// Store Public Key
	sm.HandleFunc("/store-public-key", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetstorePublicKeyHandler(w, r)
		}
	})

	// Accounts
	sm.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUserAccounts(w, r)
		case http.MethodPost:
			switch r.URL.Query().Get("type") {
			case "bank":
				controller.GetAddBankAccount(w, r)
			}
		case http.MethodDelete:
			controller.GetDeleteAccount(w, r)
		}
	})

	sm.HandleFunc("/accounts/verify", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			controller.GetVerifyAccount(w, r)
		}
	})

	// Banks
	sm.HandleFunc("/accounts/banks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetBanks(w, r)
		case http.MethodPost:
			controller.GetAddBank(w, r)
		}
	})

	// Transactions
	sm.HandleFunc("/accounts/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetUserTransactions(w, r)
		case http.MethodPost:
			controller.GetRequestTransaction(w, r)
		}
	})
	// Transactions
	sm.HandleFunc("/api/accounts/airtime/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.GetAirtimeTransactions(w, r)
		}
	})

	sm.HandleFunc("/accounts/balance", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetCheckBalanceApi(w, r)
		}
	})

	sm.HandleFunc("/accounts/apply-for-token", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetApplyForToken(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-keys", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRegisterKeys(w, r)
		case http.MethodGet:
			controller.GetApiKeys(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-hosted-intiate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRequestHostedTransactionInitiate(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-hosted", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetVerifyTransactionHosted(w, r)
		case http.MethodGet:
			controller.GetTransactionDetails(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-intiate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRequestTransactionInitiate(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-intiate-hosted", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.GetRequestTransactionInitiateForHosted(w, r)
		}
	})

	sm.HandleFunc("/accounts/transactions-varify", func(w http.ResponseWriter, r *http.Request) {
		controller.GetVerifyTransaction(w, r)
	})

	sm.HandleFunc("/accounts/transactions-dashboard", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.TransactionsDashboard(w, r)
		}
	})

	sm.HandleFunc("/account/mpesa/ussd-push", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.MpesaUssdPush(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})
	sm.HandleFunc("/account/mpesa/transaction-status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.MpesaUssdTransactionStatus(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})
	// MPESA B2C

	sm.HandleFunc("/account/mpesa/payment/b2c", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.MPesaB2C(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})
	sm.HandleFunc("/account/mpesa/b2c/transaction-status", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.MPesaB2CTransactionStatus(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})
	sm.HandleFunc("/api/account/telebirr/ussd-push", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.TelebirrUssdPush(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	sm.HandleFunc("/api/accounts/telebirr/transactionstatus", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.TelebirrUssdTransactionStatus(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	// MPESA B2C

	sm.HandleFunc("/api/accounts/telebirr/payment/b2c", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.TelebirrB2C(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})
	sm.HandleFunc("/api/accounts/telebirr/b2c/transactionstatus", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.TelebirrB2CTransactionStatus(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	// TELEPORT YIMULU AIRTIME TOPUP

	sm.HandleFunc("/account/teleport/yimulu/initiate", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CheckBalance(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
		}
	})

	sm.HandleFunc("/account/teleport/yimulu/topup", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.TopupClientBalance(w, r)
		default:
			http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
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
