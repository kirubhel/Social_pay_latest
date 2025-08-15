package rest

import (
	"encoding/json"
	"log"
	"net/http"

	auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
	"github.com/socialpay/socialpay/src/pkg/key/usecase"

	"github.com/google/uuid"
)

type Error struct {
	Type    string `json:"type"`
	Message string `json:"message"`
}

func (err Error) Error() string {
	return err.Message
}

type Controller struct {
	log        *log.Logger
	interactor usecase.Interactor
	repo       usecase.KeyRepository
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

type User struct {
	Id        uuid.UUID `json:"id"`
	SirName   string    `json:"sir_name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	UserType  string    `json:"user_type"`
}

func New(log *log.Logger, sm *http.ServeMux, interactor usecase.Interactor, repo usecase.KeyRepository, auth auth.Controller) Controller {

	controller := Controller{log: log, interactor: interactor, repo: repo, auth: auth}

	sm.HandleFunc("/key/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			{
				controller.CreateKey(w, r)
			}
		}
	})
	sm.HandleFunc("/key/get-key-by-token", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			{
				controller.GetApiKeyByToken(w, r)
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
