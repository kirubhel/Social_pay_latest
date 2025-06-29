package core

type Authorization struct {
	MerchantCode       string `json:"merchantCode"`
	MerchantTillNumber string `json:"merchantTillNumber"`
	RequestID          string `json:"requestId"`
	RequestSignature   string `json:"requestSignature"`
}

type PaymentRequest struct {
	PayerPhone       string `json:"payerPhone,omitempty"`
	Reason           string `json:"reason,omitempty"`
	Amount           string `json:"amount,omitempty"`
	ExternalReference string `json:"externalReference"`
	CallbackURL      string `json:"callbackUrl,omitempty"`
}

type DebitRequest struct {
	Authorization  Authorization  `json:"authorization,omitempty"`
	PaymentRequest PaymentRequest `json:"paymentRequest"`
}



// callback request body 
type AwashPayload struct {
	Amount           float64 `json:"amount"`
	DateRequested    string  `json:"dateRequested"`
	ExternalReference string `json:"externalReference"`
	PayerPhone       string  `json:"payerPhone"`
	ReturnCode       int     `json:"returnCode"`
	ReturnMessage    string  `json:"returnMessage"`
	Status           string  `json:"status"`

}