package usecase

import (
	"log"
)

type Error struct {
	Type    string
	Message string
}

func (err Error) Error() string {
	return err.Message
}

// SMS Sender
type SMSSender interface {
	SendSMS(string, string) error
}

type Usecase struct {
	log  *log.Logger
	repo Repo
	sms  SMSSender
}

func New(log *log.Logger, repo Repo, sms SMSSender) Interactor {
	return Usecase{log: log, repo: repo, sms: sms}
}
