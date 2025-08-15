package processors

import (
	"github.com/socialpay/socialpay/src/pkg/account/adapter/gateway/mpesa"
)

func ProcessMPesa(id string, amount float64, phone string) error {
	response, err := mpesa.HandleSTKPushRequest(mpesa.USSDPushRequest{
		BusinessShortCode: "3eda08ae-56c3-42b9-b195-d32be8eb5aca",
		Password:          "35234b15fb27cdcaff832ab20bbc937921d9e0eea6b4c538a7c90d5a6971926c",
		Timestamp:         "20240216165627",
		TransactionType:   "PayProximityMerchant",
		Amount:            amount,
		PartyA:            "50729086-28c8-4030-b4e4-7c8f1992d495",
		PartyB:            "3eda08ae-56c3-42b9-b195-d32be8eb5aca",
		PhoneNumber:       phone,
		TransactionDesc:   "",
		CallBackURL:       "https://api.socialpay.co/api/v1/checkout/transactions/notify",
		AccountReference:  id,
		MerchantName:      "C2B_STK_BUYGOODS_CLIENT_FASTPAY",
	})
	if err != nil {
		return err
	}
	// Handle the response if needed
	_ = response
	return nil
}
