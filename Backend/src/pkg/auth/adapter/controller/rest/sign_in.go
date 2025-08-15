package rest

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/auth/usecase"
	"github.com/socialpay/socialpay/src/pkg/merchants/core/entity"

	"github.com/google/uuid"
)

// Response
type AuthResponse struct {
	NextStep string `json:"next_step,omitempty"`
	Message  string `json:"message,omitempty"`
	Token    *struct {
		Active  string `json:"active"`
		Refresh string `json:"refresh"`
	} `json:"token,omitempty"`
	User *struct {
		Id        uuid.UUID `json:"id"`
		SirName   string    `json:"sir_name,omitempty"`
		FirstName string    `json:"first_name"`
		LastName  string    `json:"last_name,omitempty"`
		UserType  string    `json:"user_type,omitempty"`
	} `json:"user,omitempty"`
	Merchants []entity.Merchant `json:"merchants"`
}

func (controller Controller) GetSignIn(w http.ResponseWriter, r *http.Request) {
	controller.log.SetPrefix("[AUTH] [ADAPTER] [CONTROLLER] [REST] [GetSignIn] ")
	// Request
	type Request struct {
		Token string `json:"token"`
		Code  string `json:"code"`
		Phone struct {
			Prefix string `json:"prefix"`
			Number string `json:"number"`
		} `json:"phone"`
	}

	// Parse request
	var req Request

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		controller.log.Println(err)
		// Send error response
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "INVALID_REQUEST",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}

	// Auth phone
	log.Println("requesting for phone authentication")
	err = controller.interactor.AuthPhone(req.Token, req.Phone.Prefix, req.Phone.Number, req.Code)
	if err != nil {
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    err.(usecase.Error).Type,
				Message: err.(usecase.Error).Message,
			},
		}, http.StatusBadRequest)
		return
	}

	// Create session
	log.Println("requesting for session creation")
	session, at, err := controller.interactor.CreateSession(req.Token)
	if err != nil {
		log.Printf("[GetSignIn] CreateSession error: %+v\n", err)
		if err, ok := err.(usecase.Error); ok {
			switch err.Type {
			case "SET_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SET_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "CHECK_PASSWORD":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "CHECK_PASSWORD",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			case "SIGN_UP":
				{
					SendJSONResponse(w, Response{
						Success: true,
						Data: AuthResponse{
							NextStep: "SIGN_UP",
							Message:  err.Message,
						},
					}, http.StatusAccepted)
					return
				}
			}
		}
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNSPECIFIED",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	log.Printf("[GetSignIn] Session created for user: %s\n", session.User.Id.String())
	merchants, err := controller.merchantRepo.GetMerchantsByUserID(session.User.Id)
	if err != nil {
		log.Printf("[GetSignIn] GetMerchantsByUserID error: %s\n", err.Error())
		SendJSONResponse(w, Response{
			Success: false,
			Error: Error{
				Type:    "UNSPECIFIED",
				Message: err.Error(),
			},
		}, http.StatusBadRequest)
		return
	}
	log.Printf("[GetSignIn] Returning response for user: %s with %d merchants\n", session.User.Id.String(), len(merchants))
	log.Printf("[GetSignIn] Merchants: %+v\n", merchants)
	mappedMerchants := make([]entity.Merchant, 0)
	for _, merchant := range merchants {
		mappedMerchants = append(mappedMerchants, entity.Merchant{
			MerchantID:         merchant.MerchantID,
			UserID:             merchant.UserID,
			LegalName:          merchant.LegalName,
			TradingName:        merchant.TradingName,
			BusinessRegNumber:  merchant.BusinessRegNumber,
			TaxIdentifier:      merchant.TaxIdentifier,
			IndustryType:       merchant.IndustryType,
			BusinessType:       merchant.BusinessType,
			IsBettingClient:    merchant.IsBettingClient,
			LoyaltyCertificate: merchant.LoyaltyCertificate,
			WebsiteURL:         merchant.WebsiteURL,
			EstablishedDate:    merchant.EstablishedDate,
			CreatedAt:          merchant.CreatedAt,
			UpdatedAt:          merchant.UpdatedAt,
			Status:             merchant.Status,
		})
	}
	// // Return Response
	SendJSONResponse(w, Response{
		Success: true,
		Data: AuthResponse{
			Token: &struct {
				Active  string "json:\"active\""
				Refresh string "json:\"refresh\""
			}{
				Active:  at,
				Refresh: session.Token,
			},
			User: &struct {
				Id        uuid.UUID "json:\"id\""
				SirName   string    "json:\"sir_name,omitempty\""
				FirstName string    "json:\"first_name\""
				LastName  string    "json:\"last_name,omitempty\""
				UserType  string    "json:\"user_type,omitempty\""
			}{
				Id:        session.User.Id,
				SirName:   session.User.SirName,
				FirstName: session.User.FirstName,
				LastName:  session.User.LastName,
				UserType:  session.User.UserType,
			},
			Merchants: mappedMerchants,
		},
	}, http.StatusOK)
}
