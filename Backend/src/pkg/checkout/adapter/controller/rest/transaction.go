package rest

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/socialpay/socialpay/src/pkg/checkout/adapter/controller/procedure"
	"github.com/socialpay/socialpay/src/pkg/checkout/core/entity"
	key "github.com/socialpay/socialpay/src/pkg/key/core/entity"
)

type Transaction struct {
	Id      string `json:"id"`
	To      string `json:"to"`
	Ttl     int    `json:"ttl"`
	Pricing struct {
		Amount float64              `json:"amount"`
		Fees   []map[string]float64 `json:"fees"`
	} `json:"pricing"`
	Status struct {
		Value   entity.TransactionStatus `json:"value"`
		Message string                   `json:"message"`
	} `json:"status"`
	GateWay string `json:"gateway"`
	Type    string `json:"type"`

	Details map[string]interface{} `json:"details"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func encodeTransaction(v entity.Transaction) Transaction {
	return Transaction{
		Id:  v.Id,
		To:  v.For,
		Ttl: v.Ttl,
		Pricing: struct {
			Amount float64              "json:\"amount\""
			Fees   []map[string]float64 "json:\"fees\""
		}(v.Pricing),
		Status: struct {
			Value   entity.TransactionStatus "json:\"value\""
			Message string                   "json:\"message\""
		}(v.Status),
		GateWay:   v.GateWay,
		Type:      v.Type,
		Details:   v.Details,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func (controller Controller) GetInitTransaction(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		To        string                 `json:"to"`
		Medium    string                 `json:"medium"`
		Amount    float64                `json:"amount"`
		Details   map[string]interface{} `json:"details"`
		Redirects struct {
			Success  string `json:"success"`
			Cancel   string `json:"cancel"`
			Declined string `json:"declined"`
		} `json:"redirects"`
	}

	var req Request

	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.Error(),
				Message: err.Error(),
			},
		}, 400)
		return
	}
	host := os.Getenv("HOST")
	if host == "" {
		host = "http://196.190.251.194:8082" // Default value if the environment variable is not set
	}
	api_key := r.Header.Get("api-key")

	if api_key == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}
	fmt.Println("api_key", api_key)
	response, err := http.Get(host + "/key/get-key-by-token?token=" + api_key)
	if err != nil {
		controller.log.Println("Error creating request:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to create request to merchant service",
			},
		}, http.StatusInternalServerError)
		return
	}

	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		controller.log.Println("Error: received non-200 response code")
		fmt.Println("response", response.Body, "status", response.StatusCode)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to retrieve merchant information",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Read raw body first
	bodyBytes, _ := io.ReadAll(response.Body)
	fmt.Println("Raw response body:", string(bodyBytes))
	type apiKeyResponse struct {
		APIKey struct {
			Token      string `json:"token"`
			PrivateKey string `json:"private_key"`
			PublicKey  string `json:"public_key"`
		} `json:"apikey"`
	}

	// Replace your current decoding code with:
	var keyResp apiKeyResponse
	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&keyResp)
	if err != nil {
		controller.log.Println("Error decoding response:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INTERNAL_SERVER_ERROR",
				Message: "Failed to decode merchant information",
			},
		}, http.StatusInternalServerError)
		return
	}

	// Convert to your APIKey struct
	key := key.APIKey{
		PrivateKey: keyResp.APIKey.PrivateKey,
		PublicKey:  keyResp.APIKey.PublicKey,
		// Set other fields if available
	}
	fmt.Println(req)
	requestJSON, err := json.Marshal(req)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Private Key:", key.PrivateKey)
	fmt.Println("Public Key:", key.PublicKey)
	encrypted, err := procedure.Encrypt(key.PublicKey, requestJSON)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Encrypted request:", encrypted)

	// 5. Decrypt the request
	decrypted, err := procedure.Decrypt(key.PrivateKey, encrypted)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Decrypted request:", string(decrypted))
	err = json.Unmarshal(decrypted, &req)
	if err != nil {
		log.Fatalf("Failed to unmarshal decrypted request: %v", err)
	}

	// 4. Now you can use the `req` object normally
	fmt.Printf("Parsed Request: %+v\n", req)
	var decryptedRequest Request
	err = json.Unmarshal(decrypted, &decryptedRequest)
	if err != nil {
		log.Fatalf("Failed to unmarshal decrypted request: %v", err)
	}
	// Service Call
	txn, err := controller.interactor.InitTransaction(decryptedRequest.To, decryptedRequest.Medium, decryptedRequest.Amount, struct {
		Success string
		Cancel  string
		Decline string
	}{
		Success: req.Redirects.Success,
		Cancel:  req.Redirects.Cancel,
		Decline: req.Redirects.Declined,
	}, req.Details)

	if err != nil {
		fmt.Println(err)
		controller.log.Println(err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.Error(),
				Message: "",
			},
		}, 400)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: map[string]interface{}{
			"status": map[string]interface{}{
				"current": "TXN_INIT",
				"next":    "TXN_PROCESS",
			},
			"transaction": encodeTransaction(*txn),
		},
	}, 201)
}

func (controller Controller) GetConfirmTransaction(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")

	// Service Call
	res, err := controller.interactor.ConfirmTransaction(id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.Error(),
				Message: "",
			},
		}, 400)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: map[string]interface{}{
			"status": map[string]interface{}{
				"current": "TXN_PROCESS",
				"next":    "TXN_CHECKOUT",
			},
			"checkout": res,
		},
	}, 201)
}

func (controller Controller) GetTransaction(w http.ResponseWriter, r *http.Request) {
	var id = r.URL.Query().Get("id")

	txn, err := controller.interactor.GetTransaction(id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.Error(),
				Message: "",
			},
		}, 400)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    encodeTransaction(txn),
	}, 201)
}

func (controller Controller) UpdateCybersourceStatus(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	var id = r.Form.Get("req_transaction_uuid")
	var reasonCode = r.Form.Get("reason_code")
	var decision = r.Form.Get("decision")

	switch decision {
	case "ACCEPT":
		{
			if reasonCode == "100" {
				controller.log.Println("Verify Payment")
				controller.interactor.UpdatePaymentStatus(id, struct {
					Value   entity.TransactionStatus
					Message string
				}{
					Value:   entity.TxnCompleted,
					Message: "Transaction Completed",
				})
			}
			return
		}
	case "CANCEL":
		{
			controller.interactor.UpdatePaymentStatus(id, struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value: entity.TxnCanceled,
			})
			return
		}
	case "DECLINE":
		{
			controller.interactor.UpdatePaymentStatus(id, struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value: entity.TxnDeclined,
			})
			return
		}
	}
}

func (controller Controller) UpdateCBEBirrStatus(w http.ResponseWriter, r *http.Request) {
	type Header struct {
		XMLName                  xml.Name `xml:"Header"`
		Version                  string   `xml:"Version"`
		OriginatorConversationID string   `xml:"OriginatorConversationID"`
		ConversationID           string   `xml:"ConversationID"`
	}

	type Body struct {
		XMLName           xml.Name `xml:"Body"`
		ResultType        string   `xml:"ResultType"`
		ResultCode        string   `xml:"ResultCode"`
		ResultDesc        string   `xml:"ResultDesc"`
		TransactionResult struct {
			TransactionID string `xml:"TransactionID"`
		} `xml:"TransactionResult"`
	}

	type Result struct {
		XMLName xml.Name `xml:"Result"`
		Header  Header   `xml:"Header"`
		Body    Body     `xml:"Body"`
	}

	type Body1 struct {
		XMLName xml.Name `xml:"Body"`
		Result  Result   `xml:"Result"`
	}

	type Envelope struct {
		XMLName xml.Name `xml:"Envelope"`
		Body    Body1    `xml:"Body"`
	}

	var parsedResult Envelope

	decoder := xml.NewDecoder(r.Body)
	err := decoder.Decode(&parsedResult)
	controller.log.Println(err)
	controller.log.Println(parsedResult)

	switch fmt.Sprintf("%s%s", parsedResult.Body.Result.Body.ResultType, parsedResult.Body.Result.Body.ResultCode) {
	case "00":
		{
			controller.interactor.UpdatePaymentStatus(parsedResult.Body.Result.Header.OriginatorConversationID, struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnCompleted,
				Message: "Transaction Completed",
			})
			return
		}
	case "01":
		{
			controller.interactor.UpdatePaymentStatus(parsedResult.Body.Result.Header.OriginatorConversationID, struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnCanceled,
				Message: "Transaction Cancelled",
			})
			return
		}
	default:
		{
			controller.interactor.UpdatePaymentStatus(parsedResult.Body.Result.Header.OriginatorConversationID, struct {
				Value   entity.TransactionStatus
				Message string
			}{
				Value:   entity.TxnDeclined,
				Message: "Transaction Declined",
			})
			return
		}
	}
}

func (controller Controller) GetHandleTransactionNotification(w http.ResponseWriter, r *http.Request) {
	switch r.Host {
	case "https://secureacceptance.cybersource.com":
		{
			controller.UpdateCybersourceStatus(w, r)
			return
		}
	case "http://196.190.251.169:33180":
		{
			controller.UpdateCybersourceStatus(w, r)
			return
		}
	}
}
