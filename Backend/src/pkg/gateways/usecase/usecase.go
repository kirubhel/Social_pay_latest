package usecase

import (
	"log"
)

// Errors

type Error struct {
	Type    string
	Message string
	Details string
}

func (err Error) Error() string {
	return err.Message
}

type Usecase struct {
	log  *log.Logger
	repo Repo
}

func New(log *log.Logger, repo Repo) Interactor {
	return Usecase{log, repo}
}
