package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/socialpay/socialpay/src/pkg/gateways/usecase"
)

// Error struct for handling error responses
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Details string `json:"details,omitempty"`
}

func (err Error) Error() string {
	return err.Message
}

// Controller struct
type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	sm         *http.ServeMux
}

// Response struct for API responses
type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   error       `json:"error,omitempty"`
}

// User struct
type User struct {
	Id        uuid.UUID `json:"id"`
	SirName   string    `json:"sir_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UserType  string    `json:"user_type"`
}

// New function to initialize the Controller
func New(log *log.Logger, sm *http.ServeMux, interactor usecase.Interactor) (*Controller, error) {

	// Initialize the Controller with RabbitMQ config and connection
	controller := &Controller{
		log:        log,
		interactor: interactor,
		sm:         sm,
	}

	// Define HTTP routes
	sm.HandleFunc("/api/v1/gateways/list/all", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.ListGateways(w, r)
		}
	})

	sm.HandleFunc("/api/v1/gateways/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.CreateGateway(w, r)
		}
	})

	sm.HandleFunc("/api/v1/gateways/update", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UpdateGateway(w, r)
		}
	})

	sm.HandleFunc("/api/v1/gateways/delete", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodDelete:
			controller.DeleteGateway(w, r)
		}
	})

	sm.HandleFunc("/api/v1/merchant/gateways/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.ListMerchantGateways(w, r)
		}
	})
	sm.HandleFunc("/api/v1/gateways/merchants/list", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.ListGatewayMerchants(w, r)
		}
	})

	sm.HandleFunc("/api/v1/gateways/merchants/link", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.LinkGatewayToMerchant(w, r)
		}
	})

	sm.HandleFunc("/api/v1/gateways/merchants/unlink", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.UnlinkGatewayFromMerchant(w, r)
		}
	})

	sm.HandleFunc("/api/v1/gateways/merchants/disable", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			controller.DisableMerchantGateway(w, r)
		}
	})
	sm.HandleFunc("/api/v1/gateways/merchants/switch", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			controller.EnableMerchantGateway(w, r)
		}
	})

	return controller, nil
}

// SendJSONResponse function to send JSON responses
func SendJSONResponse(w http.ResponseWriter, data interface{}, status int) {
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
