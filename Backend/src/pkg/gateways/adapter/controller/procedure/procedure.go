package procedure

import (
	"log"

	"github.com/socialpay/socialpay/src/pkg/auth/usecase"
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
}

func New(log *log.Logger, interactor usecase.Interactor) Controller {
	return Controller{log: log, interactor: interactor}
}
