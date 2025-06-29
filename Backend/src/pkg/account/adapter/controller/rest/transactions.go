package rest

import (
	"bytes"
	"crypto/rand"
	"crypto/tls"
	"database/sql"
	"encoding/base32"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/socialpay/socialpay/src/pkg/account/adapter/gateway/mpesa"
	"github.com/socialpay/socialpay/src/pkg/account/core/entity"
	"github.com/socialpay/socialpay/src/pkg/account/usecase"
	"github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
	auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
	"github.com/socialpay/socialpay/src/pkg/auth/infra/storage/psql"
	"github.com/socialpay/socialpay/src/pkg/utils"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

type Transaction struct {
	Id     uuid.UUID `json:"id"`
	From   Account   `json:"from"`
	To     Account   `json:"to"`
	Amount float64   `json:"amount"`
	Type   string    `json:"type"`
	Date   time.Time `json:"date"`
}

func NewTransactionFromEntity(i entity.Transaction) Transaction {
	return Transaction{
		Id:     i.Id,
		From:   NewAccountFromEntity(i.From),
		To:     NewAccountFromEntity(i.To),
		Type:   string(i.Type),
		Date:   i.CreatedAt,
		Amount: i.Amount,
	}
}

func (controller Controller) UpdateUser(w http.ResponseWriter, r *http.Request) {

	var user entity.User2
	decoder := json.NewDecoder((r.Body))
	err := decoder.Decode(&user)

	if err != nil {

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return

	}

	users, err := controller.interactor.UpdateUserUsecase(user)
	if err != nil {

		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return

	}
	SendJSONResponse(w, Response{
		Success: true,
		Data:    users,
	}, http.StatusOK)

}

func (controller Controller) GetVerifyTransactionHosted(w http.ResponseWriter, r *http.Request) {
	// Authenticate request
	log := log.New(os.Stdout, "[SOCIALPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	db, err := psql.New(log)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	// Parse request
	type Request struct {
		Phone         string `json:"phone"`
		Token         string `json:"token"`
		TwoFA         string `json:"2fa"`
		Challenge     string `json:"challenge"`
		Signature     string `json:"signature"`
		OTP           string `json:"otp"`
		ChallengeType string `json:"challenge_type"`
		// Amount      float64   `json:"amount"`
	}

	var req Request
	var transactionChallenge entity.TransactionChallange

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	transactionChallenge.Signature = req.Signature
	transactionChallenge.TwoFA = req.TwoFA
	transactionChallenge.Challenge = req.Challenge
	transactionChallenge.OTP = req.OTP

	var sender_id uuid.UUID
	sqlStmt := `select pi.user_id from auth.phones as p
join auth.phone_identities as pi on p.id = pi.phone_id
WHERE p.number = $1
;`
	err = db.QueryRow(sqlStmt, req.Phone).Scan(&sender_id)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	_, err = controller.interactor.VerifyTransaction(sender_id, req.Token, transactionChallenge, req.ChallengeType)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "success",
	}, http.StatusOK)

	// Usecase
}

func (controller Controller) GetVerifyTransaction(w http.ResponseWriter, r *http.Request) {
	// Authenticate request

	controller.log.Println("Adding Transaction")

	// Check header
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	// Parse request
	type Request struct {
		Token         string `json:"token"`
		TwoFA         string `json:"2fa"`
		Challenge     string `json:"challenge"`
		Signature     string `json:"signature"`
		OTP           string `json:"otp"`
		ChallengeType string `json:"challenge_type"`
		// Amount      float64   `json:"amount"`
	}

	var req Request
	var transactionChallenge entity.TransactionChallange

	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	transactionChallenge.Signature = req.Signature
	transactionChallenge.TwoFA = req.TwoFA
	transactionChallenge.Challenge = req.Challenge
	transactionChallenge.OTP = req.OTP

	_, err = controller.interactor.VerifyTransaction(session.User.Id, req.Token, transactionChallenge, req.ChallengeType)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "success",
	}, http.StatusOK)

	// Usecase
}

func (controller Controller) GetApiKeys(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	private_key, err := controller.interactor.GetApiKeysUsecase(session.User.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	type Response2 struct {
		PrivateKey string `json:"private_key"`
	}

	var res Response2
	res.PrivateKey = private_key
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetApplyForToken(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Username and password are required",
			},
		}, http.StatusBadRequest)
		return
	}

	private_key, err := controller.interactor.ApplyForTokenUsecase(req.Username, req.Password)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	type Response2 struct {
		Token  string `json:"token"`
		Detail struct {
			Type string `json:"type"`
			Info string `json:"info"`
		} `json:"detail"`
	}

	var res Response2
	res.Token = private_key
	res.Detail.Type = "Bearer Token"
	res.Detail.Info = "Add the Bearer token to the header to authorize."
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)
}

func (controller Controller) GetCheckBalanceApi(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| sesstion id", session.User.Id, session.Id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}
	type Request struct {
		From uuid.UUID `json:"id"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	// if req.Username == "" || req.Password == "" {
	//     SendJSONResponse(w, Response{
	//         Success: false,
	//         Error: &Error{
	//             Type:    "INVALID_REQUEST",
	//             Message: "Username and password are required",
	//         },
	//     }, http.StatusBadRequest)
	//     return
	// }

	Balance, err := controller.interactor.CheckBalance(req.From)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	type Response2 struct {
		Balance float64 `json:"balance"`
	}

	var res Response2
	res.Balance = Balance
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetRegisterKeys(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	if req.Username == "" || req.Password == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "Username and password are required",
			},
		}, http.StatusBadRequest)
		return
	}

	private_key, err := controller.interactor.CreateRegisterKeys(session.User.Id, req.Username, req.Password)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	type Response2 struct {
		PrivateKey string `json:"private_key"`
	}

	var res Response2
	res.PrivateKey = private_key
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) CompleteTransaction(w http.ResponseWriter, r *http.Request) {

	transactionID := r.FormValue("transaction_id")
	log := log.New(os.Stdout, "[SOCIALPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Ouput Adapters]
	// [DB] Postgres
	db, err := psql.New(log)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	type Request struct {
		Phone string `json:"phone"`
	}
	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST222",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var sender_id uuid.UUID
	sqlStmt := `select pi.user_id from auth.phones as p
join auth.phone_identities as pi on p.id = pi.phone_id
WHERE p.number = $1
;`

	err = db.QueryRow(sqlStmt, req.Phone).Scan(&sender_id)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// chek the balance here

	_, err = db.Exec(`UPDATE accounts.transactions SET verified=true, "from"=$1 WHERE id=$2`, sender_id, transactionID)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Transaction completed successfully",
	}, http.StatusOK)
}
func (controller Controller) GetTransactionDetails(w http.ResponseWriter, r *http.Request) {

	transactionID := r.URL.Query().Get("transaction_id")
	log := log.New(os.Stdout, "[SOCIALPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	db, err := psql.New(log)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	var transaction entity.Transaction
	var amount sql.NullString

	err = db.QueryRow("SELECT id, amount, currency FROM accounts.transactions WHERE id=$1", transactionID).Scan(&transaction.Id, &amount, &transaction.Currency)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	if amount.Valid {
		amount2, err := utils.AesDecription(amount.String)
		if err != nil {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "INVALID_REQUEST",
					Message: err.Error(),
				},
			}, http.StatusBadRequest)
			return

		}
		amount_float, err := strconv.ParseFloat(amount2, 64)
		if err != nil {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "INVALID_REQUEST",
					Message: err.Error(),
				},
			}, http.StatusBadRequest)
			return
		}
		transaction.Amount = amount_float

	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    transaction,
	}, http.StatusOK)
}

// **************************** hosted checkout initate

func (controller Controller) GetRequestHostedTransactionInitiate(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	type Request struct {
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
		Signature    string  `json:"signature"`
		Callback_url string  `json:"callback_url"`
	}

	type Payment struct {
		Amount       float64 `json:"amount"`
		Currency     string  `json:"currency"`
		Callback_url string  `json:"callback_url"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	payment := Payment{
		Amount:       req.Amount,
		Currency:     req.Currency,
		Callback_url: req.Callback_url,
	}

	jsonOutput, err := json.Marshal(payment)
	if err != nil {
		log.Fatalf("Error marshaling JSON: %v", err)
	}

	// Convert the JSON byte slice to a string
	jsonString := string(jsonOutput)

	txn, err := controller.interactor.CreateHostedTransactionInitiate(req.Amount, payment.Currency, payment.Callback_url, req.Signature, jsonString, token)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    txn,
	}, http.StatusOK)

}

func (controller Controller) GetRequestTransactionInitiateForHosted(w http.ResponseWriter, r *http.Request) {
	log := log.New(os.Stdout, "[SOCIALPAY1]", log.Lmsgprefix|log.Ldate|log.Ltime|log.Lshortfile)

	// [Ouput Adapters]
	// [DB] Postgres
	db, err := psql.New(log)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	type Request struct {
		Phone string `json:"phone"`
		Data  string `json:"data"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	var sender_id uuid.UUID
	sqlStmt := `select pi.user_id from auth.phones as p
join auth.phone_identities as pi on p.id = pi.phone_id
WHERE p.number = $1
;`

	err = db.QueryRow(sqlStmt, req.Phone).Scan(&sender_id)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	var sender_acount uuid.UUID
	sqlStmt2 := `select a.id from accounts.accounts as a 
where a.user_id= $1`
	err = db.QueryRow(sqlStmt2, sender_id).Scan(&sender_acount)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	fmt.Println("GetRequestTransactionInitiateForHosted:there ||||||||||||||||||||||||||||||||||||||||", sender_id)

	decryptedData, err := utils.AesDecription(req.Data)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return

	}

	parts := strings.Split(decryptedData, ",")

	amount, err := strconv.ParseFloat(parts[1], 64)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	currency := parts[2]
	url := parts[3]
	merchantIdString := parts[0]

	parsedUUID, err := uuid.Parse(merchantIdString)
	fmt.Println("GetRequestTransactionInitiateForHosted:new ------ ||||||||||||||||||||||||||||||||||||||||", parsedUUID)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	var reciver_acount uuid.UUID
	sqlStmt3 := `select a.id from accounts.accounts as a 
	where a.user_id= $1`
	err = db.QueryRow(sqlStmt3, parsedUUID).Scan(&reciver_acount)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	fmt.Println("$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$$wwwww", sender_acount)

	txn, err := controller.interactor.CreateTransactionInitiate(sender_id, sender_acount, reciver_acount, amount, "SocialPAY", "", "", "")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	type Response2 struct {
		Result       interface{} `json:"result"`
		Currency     string      `json:"currency"`
		Callback_url string      `json:"callback_url"`
	}

	var res Response2

	res.Callback_url = url
	res.Currency = currency
	res.Result = txn

	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetRequestTransactionInitiate(w http.ResponseWriter, r *http.Request) {

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}
	type Request struct {
		From   uuid.UUID                `json:"from"`
		To     uuid.UUID                `json:"to"`
		Type   string                   `json:"type"`
		Amount float64                  `json:"amount"`
		Medium entity.TransactionMedium `json:"medium"`
		Detail string                   `json:"details"`
	}

	var req Request
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	txn, err := controller.interactor.CreateTransactionInitiate(session.User.Id, req.From, req.To, req.Amount, req.Medium, req.Type, token, req.Detail)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    txn,
	}, http.StatusOK)

}

type USSDPushRequest struct {
	BusinessShortCode string  `json:"BusinessShortCode"`
	Password          string  `json:"Password"`
	Timestamp         string  `json:"Timestamp"`
	TransactionType   string  `json:"TransactionType"`
	Amount            float64 `json:"Amount"`
	PartyA            string  `json:"PartyA"`
	PartyB            string  `json:"PartyB"`
	PhoneNumber       string  `json:"PhoneNumber"`
	CallBackURL       string  `json:"CallBackURL"`
	AccountReference  string  `json:"AccountReference"`
	TransactionDesc   string  `json:"TransactionDesc"`
	MerchantName      string  `json:"MerchantName"`
	Signature         string  `json:"signature"`
	TwoFA             string  `json:"2fa"`
	Challenge         string  `json:"challenge"`
	OTP               string  `json:"otp"`
	ChallengeType     string  `json:"challenge_type"`
}

func (controller Controller) MpesaUssdPush(w http.ResponseWriter, r *http.Request) {
	controller.log.Println("PROCESSING M-PESA USSD Push Request")

	var req mpesa.USSDPushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	mpesaResponse, err := mpesa.HandleSTKPushRequest(req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Data: map[string]interface{}{
				"error": err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	// We Determine success based on ResponseCode
	success := false
	if responseCode, exists := mpesaResponse["ResponseCode"]; exists {
		if code := fmt.Sprintf("%v", responseCode); code == "0" {
			success = true
		}
	}

	SendJSONResponse(w, Response{
		Success: success,
		Data:    mpesaResponse,
	}, http.StatusOK)
}

func (controller Controller) MPesaB2C(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || M-PESA B2C USSD Push Request ||||||||")
	// Authenticate request
	var token string
	controller.log.Println("verify request")

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}
	fmt.Println("||||||| || M=PESA USSD Push Request ||||||||", session)
	var req mpesa.B2CPaymentRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	controller.log.Printf("log USSD Push Request ... %+v", req)

	if err := mpesa.HandleB2CPaymentRequest(req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "MPESA_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send success response
	SendJSONResponse(w, Response{Success: true, Data: "USSD Push request sent successfully"}, http.StatusOK)
}

func (controller Controller) MpesaUssdTransactionStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("||||||| Received USSD Transaction Status Callback ||||||||")
		controller.log.Println("Processing USSD Transaction Status Callback")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Log the raw callback data
		controller.log.Printf("Raw Callback Data: %s", string(body))

		envFilePath := ".env"

		if err := godotenv.Load(envFilePath); err != nil {
			controller.log.Println("Error loading .env file, proceeding without it:", err)
		}

		forwardURL := os.Getenv("SOCIALPAY_API_URL")
		if forwardURL == "" {
			http.Error(w, "SOCIALPAY_API_URL is not configured", http.StatusInternalServerError)
			return
		}

		// Parse the callback data
		var callbackData map[string]interface{}
		if err := json.Unmarshal(body, &callbackData); err != nil {
			controller.log.Printf("Error unmarshalling callback data: %v", err)
			http.Error(w, "Invalid callback data", http.StatusBadRequest)
			return
		}

		// Safely extract relevant information with type checking
		stkCallback, ok := callbackData["Envelope"].(map[string]interface{})["Body"].(map[string]interface{})["stkCallback"].(map[string]interface{})
		if !ok {
			controller.log.Printf("Invalid callback structure: stkCallback not found")
			http.Error(w, "Invalid callback structure", http.StatusBadRequest)
			return
		}

		// Handle MerchantRequestID
		referenceId, ok := stkCallback["MerchantRequestID"].(string)
		if !ok {
			controller.log.Printf("Invalid MerchantRequestID type")
			http.Error(w, "Invalid MerchantRequestID", http.StatusBadRequest)
			return
		}

		// Handle ResultCode (which can be either string or float64)
		var resultCode float64
		switch v := stkCallback["ResultCode"].(type) {
		case float64:
			resultCode = v
		case string:
			// If it's a string error code (like "TP40113"), we treat it as a failure
			resultCode = 1 // Non-zero indicates failure
		default:
			controller.log.Printf("Invalid ResultCode type: %T", stkCallback["ResultCode"])
			http.Error(w, "Invalid ResultCode", http.StatusBadRequest)
			return
		}

		// Handle ResultDesc
		message, ok := stkCallback["ResultDesc"].(string)
		if !ok {
			controller.log.Printf("Invalid ResultDesc type")
			http.Error(w, "Invalid ResultDesc", http.StatusBadRequest)
			return
		}

		// Handle CheckoutRequestID
		providerTxId, ok := stkCallback["CheckoutRequestID"].(string)
		if !ok {
			controller.log.Printf("Invalid CheckoutRequestID type")
			http.Error(w, "Invalid CheckoutRequestID", http.StatusBadRequest)
			return
		}

		timestamp := time.Now().Format(time.RFC3339)
		typeOfService := "MPESA"

		// Determine status based on ResultCode
		status := "SUCCESS"
		if resultCode != 0 {
			status = "FAILURE"

			// Special handling for password errors
			if resultDesc, ok := stkCallback["ResultDesc"].(string); ok &&
				(strings.Contains(resultDesc, "security credential") ||
					strings.Contains(resultDesc, "password")) {
				message = "Incorrect PIN entered"
			}
		}

		// Convert the callbackData to a JSON string for providerData
		providerDataBytes, err := json.Marshal(callbackData)
		if err != nil {
			controller.log.Printf("Error marshalling provider data: %v", err)
			http.Error(w, "Failed to prepare provider data", http.StatusInternalServerError)
			return
		}
		providerDataString := string(providerDataBytes)

		// Prepare the new payload with providerData as string
		newPayload := map[string]interface{}{
			"referenceId":  referenceId,
			"status":       status,
			"message":      message,
			"providerTxId": providerTxId,
			"providerData": providerDataString,
			"timestamp":    timestamp,
			"type":         typeOfService,
		}

		payloadBytes, err := json.Marshal(newPayload)
		if err != nil {
			controller.log.Printf("Error marshalling new payload: %v", err)
			http.Error(w, "Failed to prepare forwarding payload", http.StatusInternalServerError)
			return
		}

		// Log the prepared payload
		controller.log.Printf("Forwarding Payload: %s", string(payloadBytes))

		// Forward the payload
		forwardResponse, err := http.Post(forwardURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			controller.log.Printf("Failed to forward callback: %v", err)
			http.Error(w, "Failed to forward callback", http.StatusInternalServerError)
			return
		}
		defer forwardResponse.Body.Close()

		if forwardResponse.StatusCode != http.StatusOK {
			responseBody, _ := io.ReadAll(forwardResponse.Body)
			controller.log.Printf("Failed to forward callback, received status: %s, Response Body: %s", forwardResponse.Status, string(responseBody))
			http.Error(w, "Failed to forward callback", http.StatusInternalServerError)
			return
		}

		// Respond to the original request
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "Callback received and forwarded successfully"}`)

	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func (controller Controller) MPesaB2CTransactionStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Println("||||||| Received USSD Transaction Status Callback ||||||||")
		controller.log.Println("Processing USSD Transaction Status Callback")
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		// Log the raw callback data
		controller.log.Printf("Raw Callback Data: %s", string(body))

		// Send a success response
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"message": "Callback received successfully"}`)

	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

// Simplified request struct for the API endpoint
type TelebirrPaymentRequest struct {
	CommandID                string  `json:"command_id"`
	OriginatorConversationID string  `json:"conversation_id"`
	PrimaryParty             string  `json:"primary_party"`
	IdentifierType           int     `json:"identifier_type"`
	Identifier               string  `json:"identifier"`
	ReceiverParty            string  `json:"receiver_party"`
	Amount                   float64 `json:"amount"`
	Currency                 string  `json:"currency"`
}

// Full SOAP request struct
type TelebirrUSSDPushRequest struct {
	CommandID                string  `json:"CommandID"`
	OriginatorConversationID string  `json:"OriginatorConversationID"`
	ThirdPartyID             string  `json:"ThirdPartyID"`
	Password                 string  `json:"Password"`
	ResultURL                string  `json:"ResultURL"`
	Timestamp                string  `json:"Timestamp"`
	IdentifierType           int     `json:"IdentifierType"`
	Identifier               string  `json:"Identifier"`
	SecurityCredential       string  `json:"SecurityCredential"`
	ShortCode                string  `json:"ShortCode"`
	PrimaryParty             string  `json:"PrimaryParty"`
	ReceiverParty            string  `json:"ReceiverParty"`
	Amount                   float64 `json:"Amount"`
	Currency                 string  `json:"Currency"`
}

type TelebirrSoapResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName  xml.Name `xml:"Body"`
		Response struct {
			XMLName                  xml.Name `xml:"Response"`
			Version                  string   `xml:"Header>Version"`
			OriginatorConversationID string   `xml:"Header>OriginatorConversationID"`
			ConversationID           string   `xml:"Header>ConversationID"`
			ResponseCode             int      `xml:"Body>ResponseCode"`
			ResponseDesc             string   `xml:"Body>ResponseDesc"`
			ServiceStatus            int      `xml:"Body>ServiceStatus"`
		} `xml:"Response"`
	} `xml:"Body"`
}

func (controller Controller) TelebirrUssdPush(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || PROCESSING TELEBIRR USSD PUSH REQUEST ||||||||")

	fmt.Println("||||||||||||| TELEBIRR USSD Push Request ||||||||||||||||||")

	// Decode the JSON request
	var req TelebirrPaymentRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	controller.log.Printf("Log USSD Push Request ... %+v", req)
	envFilePath := ".env"
	if err := godotenv.Load(envFilePath); err != nil {
		controller.log.Println("Error loading .env file:", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "CONFIG_ERROR",
				Message: "Failed to load configuration",
			},
		}, http.StatusInternalServerError)
		return
	}

	if req.PrimaryParty == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: "primary_party (payer phone number) is required",
			},
		}, http.StatusBadRequest)
		return
	}

	Request := TelebirrUSSDPushRequest{
		CommandID:                req.CommandID,
		OriginatorConversationID: req.OriginatorConversationID,
		ThirdPartyID:             os.Getenv("TELEBIRR_THIRD_PARTY_ID"),
		Password:                 os.Getenv("TELEBIRR_PASSWORD"),
		ResultURL:                os.Getenv("TELEBIRR_RESULT_URL"),
		Timestamp:                time.Now().Format("20060102030405"),
		IdentifierType:           req.IdentifierType,
		Identifier:               req.Identifier,
		SecurityCredential:       os.Getenv("TELEBIRR_SECURITY_CREDENTIAL"),
		ShortCode:                os.Getenv("TELEBIRR_SHORT_CODE"),
		PrimaryParty:             req.PrimaryParty,
		ReceiverParty:            req.ReceiverParty,
		Amount:                   req.Amount,
		Currency:                 req.Currency,
	}

	// Send the SOAP request and get the response
	soapRequest := generateSoapRequest(Request)
	soapResponse, err := sendSoapRequest(soapRequest)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "TELEBIRR_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	// Log the SOAP response
	controller.log.Printf("SOAP Response: %+v", soapResponse)

	// Check the response code
	if soapResponse.Body.Response.ResponseCode != 0 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "TELEBIRR_ERROR",
				Message: soapResponse.Body.Response.ResponseDesc,
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send success response
	SendJSONResponse(w, Response{
		Success: true,
		Data: map[string]interface{}{
			"message":         "Telebirr USSD Push request sent successfully",
			"conversation_id": soapResponse.Body.Response.ConversationID,
			"response_code":   soapResponse.Body.Response.ResponseCode,
			"response_desc":   soapResponse.Body.Response.ResponseDesc,
			"service_status":  soapResponse.Body.Response.ServiceStatus,
		},
	}, http.StatusOK)
}

// constructs the SOAP request from the TelebirrUSSDPushRequest struct
func generateSoapRequest(req TelebirrUSSDPushRequest) string {
	soapRequest := fmt.Sprintf(`
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:api="http://cps.huawei.com/cpsinterface/api_requestmgr" xmlns:req="http://cps.huawei.com/cpsinterface/request" xmlns:com="http://cps.huawei.com/cpsinterface/common">
   <soapenv:Header/>
   <soapenv:Body>
      <api:Request>
         <req:Header>
            <req:Version>1.0</req:Version>
            <req:CommandID>%s</req:CommandID>
            <req:OriginatorConversationID>%s</req:OriginatorConversationID>
            <req:Caller>
               <req:CallerType>2</req:CallerType>
               <req:ThirdPartyID>%s</req:ThirdPartyID>
               <req:Password>%s</req:Password>
               <req:ResultURL>%s</req:ResultURL>
            </req:Caller>
            <req:KeyOwner>1</req:KeyOwner>
            <req:Timestamp>%s</req:Timestamp>
         </req:Header>
         <req:Body>
            <req:Identity>
               <req:Initiator>
                  <req:IdentifierType>%d</req:IdentifierType>
                  <req:Identifier>%s</req:Identifier>
                  <req:SecurityCredential>%s</req:SecurityCredential>
                  <req:ShortCode>%s</req:ShortCode>
               </req:Initiator>
               <req:PrimaryParty>
                  <req:IdentifierType>1</req:IdentifierType>
                  <req:Identifier>%s</req:Identifier>
               </req:PrimaryParty>
               <req:ReceiverParty>
                  <req:IdentifierType>4</req:IdentifierType>
                  <req:Identifier>%s</req:Identifier>
               </req:ReceiverParty>
            </req:Identity>
            <req:TransactionRequest>
               <req:Parameters>
                  <req:Amount>%.2f</req:Amount>
                  <req:Currency>%s</req:Currency>
               </req:Parameters>
            </req:TransactionRequest>
         </req:Body>
      </api:Request>
   </soapenv:Body>
</soapenv:Envelope>`,
		req.CommandID,
		req.OriginatorConversationID,
		req.ThirdPartyID,
		req.Password,
		req.ResultURL,
		req.Timestamp,
		req.IdentifierType,
		req.Identifier,
		req.SecurityCredential,
		req.ShortCode,
		req.PrimaryParty,
		req.ReceiverParty,
		req.Amount,
		req.Currency,
	)

	return soapRequest
}

func sendSoapRequest(soapRequest string) (*TelebirrSoapResponse, error) {
	url := "http://10.180.70.177:30001/payment/services/APIRequestMgrService"
	fmt.Println("TELEBIRR REQUEST URL IS ", url)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(soapRequest)))
	if err != nil {
		return nil, fmt.Errorf("failed to create SOAP request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	fmt.Println("Sending SOAP Request:")
	fmt.Println(soapRequest)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send SOAP request: %w", err)
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("SOAP request failed: Status Code: %d, Response: %s", resp.StatusCode, string(body))
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse the SOAP response
	var soapResponse TelebirrSoapResponse
	err = xml.Unmarshal(body, &soapResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to parse SOAP response: %w", err)
	}

	// Return the parsed response
	return &soapResponse, nil
}

type TelebirrCallbackResponse struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    struct {
		XMLName xml.Name `xml:"Body"`
		Result  struct {
			XMLName                  xml.Name `xml:"Result"`
			Version                  string   `xml:"Header>Version"`
			OriginatorConversationID string   `xml:"Header>OriginatorConversationID"`
			ConversationID           string   `xml:"Header>ConversationID"`
			ResultType               string   `xml:"Body>ResultType"` // Changed to string
			ResultCode               string   `xml:"Body>ResultCode"` // Changed to string
			ResultDesc               string   `xml:"Body>ResultDesc"`
			TransactionResult        struct {
				TransactionID    string `xml:"TransactionID"`
				ResultParameters struct {
					ResultParameter []struct {
						Key   string `xml:"Key"`
						Value string `xml:"Value"`
					} `xml:"ResultParameter"`
				} `xml:"ResultParameters"`
			} `xml:"Body>TransactionResult"`
		} `xml:"Result"`
	} `xml:"Body"`
}

func (controller Controller) TelebirrUssdTransactionStatus(w http.ResponseWriter, r *http.Request) {
	envFilePath := ".env"

	if err := godotenv.Load(envFilePath); err != nil {
		controller.log.Println("Warning: .env file not found, proceeding without it.")
	}

	forwardURL := os.Getenv("SOCIALPAY_API_URL")
	if forwardURL == "" {
		controller.log.Println("SOCIALPAY_API_URL is not set in .env file")
		http.Error(w, "SOCIALPAY_API_URL is not configured", http.StatusInternalServerError)
		return
	}

	fmt.Println("||||||| Received Telebirr USSD Transaction Status Callback ||||||||")
	controller.log.Println("Processing USSD Transaction Status Callback")

	body, err := io.ReadAll(r.Body)
	if err != nil {
		controller.log.Printf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	controller.log.Println("Raw Callback Data:")
	controller.log.Println(string(body))

	var callbackResponse TelebirrCallbackResponse
	if err := xml.Unmarshal(body, &callbackResponse); err != nil {
		controller.log.Printf("Failed to parse SOAP response: %v", err)
		http.Error(w, "Failed to parse SOAP response", http.StatusBadRequest)
		return
	}

	// Extract transaction details
	transactionID := callbackResponse.Body.Result.TransactionResult.TransactionID
	resultDesc := callbackResponse.Body.Result.ResultDesc

	// Determine status based on ResultCode
	status := "SUCCESS"
	resultCode := callbackResponse.Body.Result.ResultCode
	if resultCode != "0" && resultCode != "8" {
		status = "FAILURE"
	}

	// Create the new payload
	newPayload := map[string]interface{}{
		"referenceId":  callbackResponse.Body.Result.OriginatorConversationID,
		"status":       status,
		"message":      resultDesc,
		"providerTxId": transactionID,
		"providerData": string(body),
		"timestamp":    time.Now().Format(time.RFC3339),
		"type":         "TELEBIRR",
	}

	// Marshal the new payload to JSON
	payloadBytes, err := json.Marshal(newPayload)
	if err != nil {
		controller.log.Printf("Error marshalling new payload: %v", err)
		http.Error(w, "Failed to prepare forwarding payload", http.StatusInternalServerError)
		return
	}

	// Forward the new payload to the specified URL
	forwardResponse, err := http.Post(forwardURL, "application/json", bytes.NewBuffer(payloadBytes))
	if err != nil {
		controller.log.Printf("Failed to forward callback: %v", err)
		http.Error(w, "Failed to forward callback", http.StatusInternalServerError)
		return
	}
	defer forwardResponse.Body.Close()

	if forwardResponse.StatusCode != http.StatusOK {
		controller.log.Printf("Failed to forward callback, received status: %s", forwardResponse.Status)
		http.Error(w, "Failed to forward callback", http.StatusInternalServerError)
		return
	}

	// Create the JSON response
	response := map[string]interface{}{
		"success":        status == "SUCCESS",
		"message":        resultDesc,
		"transaction_id": transactionID,
	}

	jsonResponse, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		controller.log.Printf("Failed to format JSON response: %v", err)
		http.Error(w, "Failed to format JSON response", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)
}

type B2CPaymentRequest struct {
	OriginatorConversationID string  `json:"OriginatorConversationID"`
	InitiatorName            string  `json:"InitiatorName"`
	SecurityCredential       string  `json:"SecurityCredential"`
	CommandID                string  `json:"CommandID"`
	PartyA                   string  `json:"PartyA"`
	PartyB                   string  `json:"PartyB"`
	Amount                   float64 `json:"Amount"`
	Currency                 string  `json:"Currency"`
	Remarks                  string  `json:"Remarks"`
	Occasion                 string  `json:"Occasion"`
	QueueTimeOutURL          string  `json:"QueueTimeOutURL"`
	ResultURL                string  `json:"ResultURL"`
	ThirdPartyID             string  `json:"ThirdPartyID"`
	Password                 string  `json:"Password"`
}

func (controller Controller) TelebirrB2C(w http.ResponseWriter, r *http.Request) {
	controller.log.Println("Processing Telebirr B2C Payment Request")

	var req B2CPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if err := validateB2CPaymentRequest(req); err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "VALIDATION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	if err := controller.processB2CPayment(req); err != nil {
		controller.log.Printf("Payment processing failed: %v", err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "PAYMENT_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    "Payment request processed successfully",
	}, http.StatusOK)
}

func validateB2CPaymentRequest(req B2CPaymentRequest) error {
	requiredFields := map[string]string{
		"OriginatorConversationID": req.OriginatorConversationID,
		"InitiatorName":            req.InitiatorName,
		"SecurityCredential":       req.SecurityCredential,
		"CommandID":                req.CommandID,
		"PartyA":                   req.PartyA,
		"PartyB":                   req.PartyB,
		"Currency":                 req.Currency,
		"Remarks":                  req.Remarks,
		"Occasion":                 req.Occasion,
		"QueueTimeOutURL":          req.QueueTimeOutURL,
		"ResultURL":                req.ResultURL,
		"ThirdPartyID":             req.ThirdPartyID,
		"Password":                 req.Password,
	}

	for field, value := range requiredFields {
		if value == "" {
			return fmt.Errorf("%s is required", field)
		}
	}

	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	if len(req.Currency) != 3 {
		return fmt.Errorf("currency must be a 3-letter ISO code")
	}

	return nil
}

func (controller Controller) processB2CPayment(req B2CPaymentRequest) error {
	timestamp := time.Now().Format("20060102150405")
	soapRequest := buildSOAPRequest(req, timestamp)

	controller.log.Printf("Generated SOAP Request:\n%s", soapRequest)

	httpReq, err := http.NewRequest(
		"POST",
		"http://10.180.70.177:30001/payment/services/APIRequestMgrService",
		bytes.NewBuffer([]byte(soapRequest)),
	)

	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "text/xml; charset=utf-8")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	controller.log.Printf("API Response: %s", string(body))

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("payment API returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func buildSOAPRequest(req B2CPaymentRequest, timestamp string) string {
	return fmt.Sprintf(`
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
                  xmlns:com="http://cps.huawei.com/cpsinterface/common" 
                  xmlns:api="http://cps.huawei.com/cpsinterface/api_requestmgr" 
                  xmlns:req="http://cps.huawei.com/cpsinterface/request">
   <soapenv:Header/>
   <soapenv:Body>
      <api:Request>
         <req:Header>
            <req:Version>1.0</req:Version>
            <req:CommandID>%s</req:CommandID>
            <req:OriginatorConversationID>%s</req:OriginatorConversationID>
            <req:Caller>
               <req:CallerType>2</req:CallerType>
               <req:ThirdPartyID>%s</req:ThirdPartyID>
               <req:Password>%s</req:Password>
               <req:ResultURL>%s</req:ResultURL>
            </req:Caller>
            <req:KeyOwner>1</req:KeyOwner>
            <req:Timestamp>%s</req:Timestamp>
         </req:Header>
         <req:Body>
            <req:Identity>
               <req:Initiator>
                  <req:IdentifierType>12</req:IdentifierType>
                  <req:Identifier>%s</req:Identifier>
                  <req:SecurityCredential>%s</req:SecurityCredential>
                  <req:ShortCode>%s</req:ShortCode>
               </req:Initiator>
               <req:ReceiverParty>
                  <req:IdentifierType>1</req:IdentifierType>
                  <req:Identifier>%s</req:Identifier>
               </req:ReceiverParty>
            </req:Identity>
            <req:TransactionRequest>
               <req:Parameters>
                  <req:Amount>%.2f</req:Amount>
                  <req:Currency>%s</req:Currency>
               </req:Parameters>
            </req:TransactionRequest>
            <req:ReferenceData>
               <req:ReferenceItem>
                  <com:Key>Remarks</com:Key>
                  <com:Value>%s</com:Value>
               </req:ReferenceItem>
               <req:ReferenceItem>
                  <com:Key>Occasion</com:Key>
                  <com:Value>%s</com:Value>
               </req:ReferenceItem>
            </req:ReferenceData>
         </req:Body>
      </api:Request>
   </soapenv:Body>
</soapenv:Envelope>`,
		req.CommandID,
		req.OriginatorConversationID,
		req.ThirdPartyID,
		req.Password,
		req.ResultURL,
		timestamp,
		req.InitiatorName,
		req.SecurityCredential,
		req.PartyA,
		req.PartyB,
		req.Amount,
		req.Currency,
		req.Remarks,
		req.Occasion)
}

func (controller Controller) TelebirrB2CTransactionStatus(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Load environment variables
		envFilePath := ".env"

		if err := godotenv.Load(envFilePath); err != nil {
			controller.log.Println("Warning: .env file not found, proceeding without it.")
		}

		forwardURL := os.Getenv("SOCIALPAY_API_URL")
		if forwardURL == "" {
			controller.log.Println("SOCIALPAY_API_URL is not set in .env file")
			http.Error(w, "SOCIALPAY_API_URL is not configured", http.StatusInternalServerError)
			return
		}

		fmt.Println("||||||| Received Telebirr B2C Transaction Status Callback ||||||||")
		controller.log.Println("Processing Telebirr B2C Transaction Status Callback")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			controller.log.Printf("Failed to read request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusInternalServerError)
			return
		}
		defer r.Body.Close()

		controller.log.Println("Raw Callback Data:")
		controller.log.Println(string(body))

		var callbackResponse TelebirrCallbackResponse
		if err := xml.Unmarshal(body, &callbackResponse); err != nil {
			controller.log.Printf("Failed to parse SOAP response: %v", err)
			http.Error(w, "Failed to parse SOAP response", http.StatusBadRequest)
			return
		}

		// Extract transaction details
		transactionID := callbackResponse.Body.Result.TransactionResult.TransactionID
		resultDesc := callbackResponse.Body.Result.ResultDesc
		status := "SUCCESS"
		if callbackResponse.Body.Result.ResultCode != "0" && callbackResponse.Body.Result.ResultCode != "8" {
			status = "FAILURE"
		}

		// Create the new payload
		newPayload := map[string]interface{}{
			"referenceId":  callbackResponse.Body.Result.OriginatorConversationID,
			"status":       status,
			"message":      resultDesc,
			"providerTxId": transactionID,
			"providerData": string(body),
			"timestamp":    time.Now().Format(time.RFC3339),
			"type":         "TELEBIRR",
		}

		// Marshal the new payload to JSON
		payloadBytes, err := json.Marshal(newPayload)
		if err != nil {
			controller.log.Printf("Error marshalling new payload: %v", err)
			http.Error(w, "Failed to prepare forwarding payload", http.StatusInternalServerError)
			return
		}

		// Forward the new payload to the specified URL
		forwardResponse, err := http.Post(forwardURL, "application/json", bytes.NewBuffer(payloadBytes))
		if err != nil {
			controller.log.Printf("Failed to forward callback: %v", err)
			http.Error(w, "Failed to forward callback", http.StatusInternalServerError)
			return
		}
		defer forwardResponse.Body.Close()

		if forwardResponse.StatusCode != http.StatusOK {
			controller.log.Printf("Failed to forward callback, received status: %s", forwardResponse.Status)
			http.Error(w, "Failed to forward callback", http.StatusInternalServerError)
			return
		}

		// Create the JSON response
		response := map[string]interface{}{
			"success":        true,
			"message":        "B2C callback received and forwarded successfully",
			"transaction_id": transactionID,
		}

		jsonResponse, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			controller.log.Printf("Failed to format JSON response: %v", err)
			http.Error(w, "Failed to format JSON response", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonResponse)

	default:
		http.Error(w, "Unsupported method", http.StatusMethodNotAllowed)
	}
}

func init() {
	envFilePath := ".env"

	err := godotenv.Load(envFilePath)
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
}

func createClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}

	client := &http.Client{
		Transport: transport,
	}

	return client
}

// Define the response structure for the /auth/login endpoint
type LoginResponse struct {
	AccessToken string `json:"accessToken"`
	Success     bool   `json:"success"`
	StatusCode  int    `json:"statusCode"`
	Message     string `json:"message"`
}

type BalanceResponse struct {
	StatusCode    int    `json:"status_code"`
	Success       bool   `json:"success"`
	Message       string `json:"message"`
	WalletBalance string `json:"wallet_balance"`
}

type AirtimeTopupRequest struct {
	ProductID      int    `json:"productId"`
	Amount         int    `json:"amount"`
	MSISDN         string `json:"msisdn"`
	TransactionRef string `json:"-"`
}

type AirtimeTopupSuccessResponse struct {
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
	Data       struct {
		Message     string `json:"message"`
		TelebirrRef string `json:"telebirrRef"`
	} `json:"data"`
}

// Error response schema
type AirtimeTopupErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Success    bool   `json:"success"`
	Message    string `json:"message"`
	Data       any    `json:"data"`
}

func generateShortTransactionRef() string {
	var b [10]byte
	_, err := rand.Read(b[:])
	if err != nil {
		panic(err)
	}
	return "Social_" + base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(b[:])
}

func (controller Controller) CheckBalance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Check Balance Request ||||||||")
	var token string
	controller.log.Println("verify request")

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)

	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}

	fmt.Println("||||||| || M-PESA USSD Push Request ||||||||", session)
	apiKey := os.Getenv("YIMULU_API_KEY")
	if apiKey == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "MISSING_API_KEY",
				Message: "API key not found in environment variables",
			},
		}, http.StatusInternalServerError)
		return
	}

	accessToken, err := Authenticate(apiKey)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "AUTHENTICATION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	balance, err := GetCurrentBalance(accessToken)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "BALANCE_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data: YimuluResponse{
			StatusCode:    balance.StatusCode,
			Message:       balance.Message,
			WalletBalance: balance.WalletBalance,
		},
	}, http.StatusOK)
}

func (controller Controller) TopupClientBalance(w http.ResponseWriter, r *http.Request) {
	fmt.Println("||||||| || Airtime Topup Request ||||||||")
	controller.log.Println("verify request")

	var token string
	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide a valid header token",
			},
		}, http.StatusUnauthorized)
		return
	}

	token = strings.Split(r.Header.Get("Authorization"), " ")[1]
	session, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: err.Error(),
			},
		}, http.StatusUnauthorized)
		return
	}

	fmt.Println("||||||| || TOPUP Request ||||||||", session)
	apiKey := os.Getenv("YIMULU_API_KEY")
	if apiKey == "" {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "MISSING_API_KEY",
				Message: "API key not found in environment variables",
			},
		}, http.StatusInternalServerError)
		return
	}
	accessToken, err := Authenticate(apiKey)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "AUTHENTICATION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	var req AirtimeTopupRequest
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	req.TransactionRef = generateShortTransactionRef()
	merchantID := session.Id

	fmt.Println("MERCHANT ID IS ", merchantID)

	AirtimeResponse, err := controller.interactor.StoreAirtimeTransaction(
		merchantID,
		req.Amount,
		req.MSISDN,
		req.TransactionRef,
	)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "STORE_AIRTIME_TRANSACTION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	fmt.Println(AirtimeResponse)

	// Proceed with the top-up request
	response, err := AirtimeTopup(accessToken, req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "TOPUP_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	var transactionRef, telebirrRef string
	transactionRef = req.TransactionRef
	telebirrRef = response.Data.TelebirrRef

	// Prepare response object
	responseWithRef := map[string]interface{}{
		"success":         response.Success,
		"status_code":     response.StatusCode,
		"message":         response.Message,
		"transaction_ref": transactionRef,
		"data":            response.Data,
	}

	err = controller.interactor.UpdateAirtimeSuccessTransaction(
		transactionRef,
		telebirrRef,
		responseWithRef,
	)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UPDATE_AIRTIME_TRANSACTION_ERROR",
				Message: err.Error(),
			},
		}, http.StatusInternalServerError)
		return
	}

	// Send the final response
	SendJSONResponse(w, Response{
		Success: true,
		Data:    responseWithRef,
	}, http.StatusOK)
}

func Authenticate(apiKey string) (string, error) {
	loginRequest := map[string]string{
		"key": apiKey,
	}

	jsonBody, err := json.Marshal(loginRequest)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request: %v", err)
	}
	client := createClient() // Use the custom client
	resp, err := client.Post("https://api.teletop.et/api/v1/client/auth/login", "application/json", bytes.NewBuffer(jsonBody))
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()
	var loginResponse LoginResponse
	err = json.NewDecoder(resp.Body).Decode(&loginResponse)
	if err != nil {
		return "", fmt.Errorf("failed to decode login response: %v", err)
	}

	if !loginResponse.Success {
		return "", fmt.Errorf("login failed: %s", loginResponse.Message)
	}

	return loginResponse.AccessToken, nil
}

func GetCurrentBalance(accessToken string) (*BalanceResponse, error) {
	client := createClient() // Use the custom client
	req, err := http.NewRequest("GET", "https://api.teletop.et/api/v1/client/info/get-current-balance", nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create balance request: %v", err)
	}
	req.Header.Set("x-api-token", accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send balance request: %v", err)
	}
	defer resp.Body.Close()
	var apiResponse struct {
		StatusCode int    `json:"statusCode"`
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		Data       struct {
			WalletBalance string `json:"walletBalance"`
		} `json:"data"`
	}

	err = json.NewDecoder(resp.Body).Decode(&apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode balance response: %v", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("balance request failed: %s", apiResponse.Message)
	}

	balanceResponse := &BalanceResponse{
		StatusCode:    apiResponse.StatusCode,
		Success:       apiResponse.Success,
		Message:       apiResponse.Message,
		WalletBalance: apiResponse.Data.WalletBalance,
	}

	return balanceResponse, nil
}

func AirtimeTopup(accessToken string, req AirtimeTopupRequest) (*AirtimeTopupSuccessResponse, error) {
	client := createClient() // Use the custom client
	payload := map[string]interface{}{
		"productId": req.ProductID,
		"amount":    req.Amount,
		"msisdn":    req.MSISDN,
	}

	jsonBody, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %v", err)
	}

	request, err := http.NewRequest("POST", "https://api.teletop.et/api/v1/client/topup-transactions/airtime-topup", bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create top-up request: %v", err)
	}

	request.Header.Set("x-api-token", accessToken)
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed to send top-up request: %v", err)
	}
	defer resp.Body.Close()
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	fmt.Println("API Response:", string(responseBody))
	var apiResponse struct {
		Success    bool   `json:"success"`
		Message    string `json:"message"`
		StatusCode int    `json:"statusCode"`
		Data       struct {
			Message     string `json:"message"`
			TelebirrRef string `json:"telebirrRef"`
		} `json:"data"`
	}

	err = json.Unmarshal(responseBody, &apiResponse)
	if err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	if !apiResponse.Success {
		return nil, fmt.Errorf("top-up failed: %s", apiResponse.Message)
	}

	return &AirtimeTopupSuccessResponse{
		Success:    apiResponse.Success,
		Message:    apiResponse.Message,
		StatusCode: apiResponse.StatusCode,
		Data: struct {
			Message     string `json:"message"`
			TelebirrRef string `json:"telebirrRef"`
		}{
			Message:     apiResponse.Data.Message,
			TelebirrRef: apiResponse.Data.TelebirrRef,
		},
	}, nil
}

func (controller Controller) GetRequestTransaction(w http.ResponseWriter, r *http.Request) {

	fmt.Println("||||||| GetRequestTransaction")
	controller.log.Println("Adding Transaction")
	// Authenticate (AuthN)

	if len(strings.Split(r.Header.Get("Authorization"), " ")) != 2 {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Please provide an authentication token in header",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Validate the token
	token := strings.Split(r.Header.Get("Authorization"), " ")[1]

	session, err := controller.auth.GetCheckAuth(token)
	fmt.Println("||| start auth")
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(procedure.Error).Type,
				Message: err.(procedure.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	fmt.Println("||| pass auth")
	// Parse the request
	type Request struct {
		From          uuid.UUID `json:"from"`
		To            uuid.UUID `json:"to"`
		Type          string    `json:"type"`
		Amount        float64   `json:"amount"`
		Medium        string    `json:"medium"`
		TwoFA         string    `json:"2fa"`
		Challenge     string    `json:"challenge"`
		Signature     string    `json:"signature"`
		OTP           string    `json:"otp"`
		ChallengeType string    `json:"challenge_type"`
	}

	var req Request
	var transactionChallenge entity.TransactionChallange
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&req)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	transactionChallenge.Signature = req.Signature
	transactionChallenge.TwoFA = req.TwoFA
	transactionChallenge.Challenge = req.Challenge

	transactionChallenge.OTP = req.OTP

	// decode the transaction challenge
	// decoder2 := json.NewDecoder(r.Body)
	// err = decoder2.Decode(&transactionChallenge)
	// if err != nil {
	// 	SendJSONResponse(w, Response{
	// 		Success: false,
	// 		Error: &Error{
	// 			Type:    "INVALID_REQUEST22",
	// 			Message: err.Error(),
	// 		},
	// 	}, http.StatusBadRequest)
	// 	return
	// }

	defer r.Body.Close()

	// Usecase

	// Response

	// Usecase [CREATE TRANSACTION]
	txn, err := controller.interactor.CreateTransaction(session.User.Id, req.From, req.To, req.Amount, req.Type, token, req.ChallengeType, transactionChallenge)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    NewTransactionFromEntity(*txn),
	}, http.StatusOK)

	w.Write([]byte("Received Request"))
}

// func (controller Controller) GetTransactions(w http.ResponseWriter, r *http.Request) {
// 	// Request
// 	type Request struct {
// 		Token string
// 	}

// 	req := Request{}

// 	token := r.Header.Get("Authorization")

// 	if len(strings.Split(token, " ")) == 2 {
// 		req.Token = strings.Split(token, " ")[1]
// 	}

// 	// Authenticate user
// 	session, err := controller.auth.GetCheckAuth(req.Token)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(auth.Error).Type,
// 				Message: err.(auth.Error).Message,
// 			},
// 		}, http.StatusUnauthorized)
// 		return
// 	}

// 	// Get user id
// 	userId := session.User.Id

// 	// Get Transactions
// 	txns, err := controller.interactor.GetUserTransactions(userId)
// 	if err != nil {
// 		SendJSONResponse(w, Response{
// 			Success: false,
// 			Error: &Error{
// 				Type:    err.(usecase.Error).Type,
// 				Message: err.(usecase.Error).Message,
// 			},
// 		}, http.StatusBadRequest)
// 		return
// 	}

// 	// Map transactions to present
// 	var _txns []Transaction = make([]Transaction, 0)
// 	for i := 0; i < len(txns); i++ {
// 		_txns = append(_txns, NewTransactionFromEntity(txns[i]))
// 	}

// 	SendJSONResponse(w, Response{
// 		Success: true,
// 		Data:    _txns,
// 	}, http.StatusOK)
// }

func (controller Controller) GetUserTransactions(w http.ResponseWriter, r *http.Request) {
	// Parse req
	type Request struct {
		Token string
	}
	var req Request

	fmt.Println("************************************** one")
	id_string := r.URL.Query().Get("id")
	fmt.Println("************************************** two")

	fmt.Println("************************************** five")

	token := strings.Split(r.Header.Get("Authorization"), " ")

	if len(token) == 2 {
		req.Token = token[1]
	}

	// Authenticate user
	controller.log.Println("PASSED -1")
	_, err := controller.auth.GetCheckAuth(req.Token)
	controller.log.Println("PASSED 0")

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	var trx interface{}
	if id_string == "" {
		fmt.Println("||||||||||||||||||| one")
		trx, err = controller.interactor.GetAllTransactions()

	} else {
		fmt.Println("||||||||||||||||||| two")

		id, err := uuid.Parse(id_string)

		if err != nil {
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "BAD_REQUEST",
					Message: "bad request",
				},
			}, http.StatusBadRequest)
			return
		}
		trx, err = controller.interactor.GetUserTransactions(id)
		if err != nil {
			var status int
			if err.(usecase.Error).Type == "UNAUTHORIZED" {
				status = http.StatusUnauthorized
			} else {
				status = http.StatusBadRequest
			}
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    err.(usecase.Error).Type,
					Message: err.(usecase.Error).Message,
				},
			}, status)
			return
		}
	}

	if err != nil {
		var status int
		if err.(usecase.Error).Type == "UNAUTHORIZED" {
			status = http.StatusUnauthorized
		} else {
			status = http.StatusBadRequest
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, status)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    trx,
	}, http.StatusOK)

}

func (controller Controller) GetAirtimeTransactions(w http.ResponseWriter, r *http.Request) {
	// Setup logging
	requestID := uuid.New().String()
	controller.log.Printf("[%s] AIRTIME TRANSACTIONS REQUEST STARTED", requestID)
	defer controller.log.Printf("[%s] AIRTIME TRANSACTIONS REQUEST COMPLETED", requestID)

	// Parse request parameters
	idString := r.URL.Query().Get("id")
	authHeader := r.Header.Get("Authorization")

	// Validate authorization header
	if authHeader == "" {
		controller.log.Printf("[%s] Missing authorization header", requestID)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Authorization header is required",
			},
		}, http.StatusUnauthorized)
		return
	}

	// Extract token
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 {
		controller.log.Printf("[%s] Invalid authorization header format", requestID)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "UNAUTHORIZED",
				Message: "Invalid authorization header format",
			},
		}, http.StatusUnauthorized)
		return
	}
	token := tokenParts[1]

	// Authenticate user
	_, err := controller.auth.GetCheckAuth(token)
	if err != nil {
		authErr, ok := err.(auth.Error)
		if !ok {
			authErr = auth.Error{
				Type:    "AUTH_ERROR",
				Message: "Authentication failed",
			}
		}

		controller.log.Printf("[%s] Authentication failed: %v", requestID, err)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    authErr.Type,
				Message: authErr.Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	var trx interface{}
	var trxErr error

	// Handle different request types
	if idString == "" {
		// Get all transactions
		controller.log.Printf("[%s] Fetching all airtime transactions", requestID)
		trx, trxErr = controller.interactor.GetAirtimeTransactions()
	} else {
		// Get transactions for specific user
		controller.log.Printf("[%s] Fetching transactions for user ID: %s", requestID, idString)

		id, err := uuid.Parse(idString)
		if err != nil {
			controller.log.Printf("[%s] Invalid user ID format: %v", requestID, err)
			SendJSONResponse(w, Response{
				Success: false,
				Error: &Error{
					Type:    "BAD_REQUEST",
					Message: "Invalid user ID format",
				},
			}, http.StatusBadRequest)
			return
		}

		trx, trxErr = controller.interactor.GetUserTransactions(id)
	}

	// Handle transaction errors
	if trxErr != nil {
		usecaseErr, ok := trxErr.(usecase.Error)
		if !ok {
			usecaseErr = usecase.Error{
				Type:    "INTERNAL_ERROR",
				Message: "Failed to retrieve transactions",
			}
		}

		status := http.StatusBadRequest
		if usecaseErr.Type == "UNAUTHORIZED" {
			status = http.StatusUnauthorized
		} else if strings.Contains(usecaseErr.Type, "NOT_FOUND") {
			status = http.StatusNotFound
		}

		controller.log.Printf("[%s] Transaction error: %v", requestID, trxErr)
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    usecaseErr.Type,
				Message: usecaseErr.Message,
			},
		}, status)
		return
	}

	// Success response
	SendJSONResponse(w, Response{
		Success: true,
		Data:    trx,
	}, http.StatusOK)
}

func (controller Controller) TransactionsDashboard(w http.ResponseWriter, r *http.Request) {
	// Parse req
	type Request struct {
		Token string
	}
	var req Request

	token := strings.Split(r.Header.Get("Authorization"), " ")

	if len(token) == 2 {
		req.Token = token[1]
	}

	// Authenticate user
	controller.log.Println("PASSED -1")
	_, err := controller.auth.GetCheckAuth(req.Token)
	controller.log.Println("PASSED 0")

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	req_type := r.URL.Query().Get("type")
	year_str := r.URL.Query().Get("year")
	year, err := strconv.Atoi(year_str)
	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "error",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	var data2 interface{}
	switch req_type {
	case "month":
		{
			trx, err := controller.interactor.TransactionsDashboardUsecase(year)
			data2 = trx
			if err != nil {

				SendJSONResponse(w, Response{
					Success: false,
					Error: &Error{
						Type:    "error",
						Message: err.Error(),
					},
				}, http.StatusBadRequest)
				return
			}
		}
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    data2,
	}, http.StatusOK)

}

func (controller Controller) GetSendOtp(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var req Request
	fmt.Println("|||||||||||||||||||||||| one")

	token := strings.Split(r.Header.Get("Authorization"), " ")

	if len(token) == 2 {
		req.Token = token[1]
	}

	fmt.Println("|||||||||||||||||||||||| two", req.Token)

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)
	fmt.Println("|||||||||||||||||||||||| 3")

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.SendOtpUsecase(session.User.Id)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetSetFingerPrint(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	var data interface{}
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&data)

	// fmt.Println("||||||||||||||||||||||||||||||||||| ", data)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error()},
		}, http.StatusBadRequest)
		return
	}

	token2, err := controller.interactor.SendSetFIngerPrintUsecase(session.User.Id, data)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetGenerateChallenge(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		// Username  string `json:"username"`
		// Challenge string `json:"challenge"`
		DeviceID string `json:"device_id"`
	}
	var res struct {
		Challenge string `json:"challenge"`
	}
	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.SendGenerateChallenge(session.User.Id, requestOne.DeviceID)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	res.Challenge = token2
	SendJSONResponse(w, Response{
		Success: true,
		Data:    res,
	}, http.StatusOK)

}

func (controller Controller) GetVerifySignatureHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		// Username  string `json:"username"`
		Challenge string `json:"challenge"`
		Signature string `json:"signature"`
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.GetverifySignature(session.User.Id, requestOne.Challenge, requestOne.Signature)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: false,
		Data:    token2,
	}, http.StatusOK)

}

func (controller Controller) GetstorePublicKeyHandler(w http.ResponseWriter, r *http.Request) {

	type Request struct {
		Token string
	}

	var requestOne struct {
		DeviceID  string `json:"device_id"`
		PublicKey string `json:"public_key"`
	}

	var req Request
	token := strings.Split(r.Header.Get("Authorization"), " ")
	if len(token) == 2 {
		req.Token = token[1]
	}

	controller.log.Println("PASSED -1")
	session, err := controller.auth.GetCheckAuth(req.Token)

	if err != nil {
		controller.log.Println("PASSED 1")
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&requestOne)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    err.(auth.Error).Type,
				Message: err.(auth.Error).Message,
			},
		}, http.StatusUnauthorized)
		return
	}

	token2, err := controller.interactor.GetstorePublicKeyHandler(requestOne.PublicKey, session.User.Id, requestOne.DeviceID)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: &Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	SendJSONResponse(w, Response{
		Success: true,
		Data:    token2,
	}, http.StatusOK)

}
