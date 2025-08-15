package usecase

import (
	"log"
)

// SMS Sender
type SMSSender interface {
	SendSMS(string, string) error
}

type Usecase struct {
	log  *log.Logger
	repo Repository
	sms  SMSSender
}

func New(log *log.Logger, repo Repository, sms SMSSender) Usecase {
	return Usecase{log, repo, sms}
}
