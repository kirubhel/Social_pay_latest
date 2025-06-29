package awash

import (
	"strings"
	"time"
)

type AwashResponse struct {
	Amount        float64    `json:"amount,omitempty"`
	Expires       CustomTime `json:"expires,omitempty"` // use custom time
	Instructions  string     `json:"instructions,omitempty"` // fixed key name
	RequestId     string     `json:"requestId"`
	ReturnCode    int        `json:"returnCode"`
	ReturnMessage string     `json:"returnMessage"`
}


type CustomTime struct {
	time.Time
}

func (ct *CustomTime) UnmarshalJSON(b []byte) error {
	str := strings.Trim(string(b), `"`)
	if str == "" || str == "null" {
		return nil
	}
	t, err := time.Parse("2006-01-02 15:04:05", str)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

