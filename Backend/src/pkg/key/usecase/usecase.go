package usecase

import (
	"log"
)

// Errors

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
	repo KeyRepository
	sms  SMSSender
}

func New(log *log.Logger, repo KeyRepository, sms SMSSender) Usecase {
	return Usecase{log, repo, sms}
}
