package validation

//TODO: define the valid status values for webhooks
// This is a placeholder for the actual implementation
// ValidWebhookStatuses defines the valid status values for webhooks
var ValidWebhookStatuses = map[string]bool{
	"success": true,
	"failed":  true,
	"pending": true,
}

// IsValidStatus checks if the given status is valid
func IsValidStatus(status string) bool {
	return ValidWebhookStatuses[status]
}
