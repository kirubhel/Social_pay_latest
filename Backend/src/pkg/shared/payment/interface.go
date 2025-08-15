package payment

// STDCallbackRequest represents the callback request structure from CBE
type STDCallbackRequest struct {
	ReferenceId  string `json:"referenceId"`
	Status       string `json:"status"`
	Message      string `json:"message"`
	ProviderTxId string `json:"providerTxId"`
	ProviderData string `json:"providerData"`
	Timestamp    string `json:"timestamp"`
	Type         string `json:"type,omitempty"`
}
