package etswitch

var CurrencyToISOCode = map[string]string{
	"ETB": "230",
	// UNCOMMENT THE FOLLOWING CURRECY IF YOU NEED OTHERS CURRENCY THAN ETB
	// "USD": "840",
	// "EUR": "978",
	// "AED": "784",
	// "GBP": "826",
	// "CNY": "156",
}

// Mapping EthSwitch response status
var ETHStatusToConstant = map[string]string{
	"declined":  "FAILED",
	"completed": "SUCCESS",
	"canceled":  "CANCELED",
}

type EthSwitchResponse struct {
	Code         int    `json:"code,omitempty"`
	ErrorCode    int    `json:"errorCode,omitempty"`
	ErrorMessage string `json:"errorMessage,omitempty"`
	OrderId      string `json:"orderId"`
	FormUrl      string `json:"formUrl"`
	MdOrder      string `json:"mdOrder"`
}
