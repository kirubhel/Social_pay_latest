package rest

import (
	// auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
	"log"
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/storage/usecase"
)

type Controller struct {
	log        *log.Logger
	sm         *http.ServeMux
	interactor usecase.Interactor
	// auth       auth.Controller
}

func New(log *log.Logger, sm *http.ServeMux, interactor usecase.Interactor) Controller {
	controller := Controller{log: log, interactor: interactor}

	// Route
	sm.HandleFunc("/storage", func(w http.ResponseWriter, r *http.Request) {
		controller.log.Println("Storage")
		switch r.Method {
		case http.MethodPost:
			{
				controller.Upload(w, r)
			}
		}
	})

	controller.sm = sm

	return controller
}
