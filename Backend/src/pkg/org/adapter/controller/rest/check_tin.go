package rest

import (
	"log"
	"net/http"
)

func (controller Controller) CheckTIN(w http.ResponseWriter, r *http.Request) {
	type Request struct {
		TIN string
	}

	type SimplifiedResponse struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data,omitempty"`
		Error   string      `json:"error,omitempty"`
	}

	type Organization struct {
		Name          string  `json:"name"`
		NameAmharic   string  `json:"name_amharic,omitempty"`
		TIN           string  `json:"tin"`
		RegNo         string  `json:"registration_number"`
		LegalType     string  `json:"legal_type"`
		Capital       float64 `json:"capital"`
		RegDate       string  `json:"registration_date"`
		Manager       string  `json:"manager,omitempty"`
		BusinessLines []struct {
			Code        int    `json:"code"`
			Description string `json:"description"`
		} `json:"business_lines,omitempty"`
	}

	var req Request
	req.TIN = r.URL.Query().Get("tin")

	rawResponse, err := controller.interactor.CheckTINRaw(req.TIN)
	if err != nil {
		log.Printf("Error checking TIN: %v", err)
		SendJSONResponse(w, SimplifiedResponse{
			Success: false,
			Error:   err.Error(),
		}, http.StatusBadRequest)
		return
	}

	// Transform the raw response into simplified structure
	simplified := Organization{
		Name:        rawResponse["BusinessName"].(string),
		NameAmharic: rawResponse["BusinessNameAmh"].(string),
		TIN:         rawResponse["Tin"].(string),
		RegNo:       rawResponse["RegNo"].(string),
		LegalType:   rawResponse["LegalCondtion"].(string),
		Capital:     rawResponse["PaidUpCapital"].(float64),
		RegDate:     rawResponse["RegDate"].(string),
	}

	// Add manager info if available
	if associates, ok := rawResponse["AssociateShortInfos"].([]interface{}); ok && len(associates) > 0 {
		if manager, ok := associates[0].(map[string]interface{}); ok {
			simplified.Manager = manager["ManagerNameEng"].(string)
		}
	}

	// Add business lines
	if businesses, ok := rawResponse["Businesses"].([]interface{}); ok && len(businesses) > 0 {
		if business, ok := businesses[0].(map[string]interface{}); ok {
			if subgroups, ok := business["SubGroups"].([]interface{}); ok {
				for _, sg := range subgroups {
					if subgroup, ok := sg.(map[string]interface{}); ok {
						simplified.BusinessLines = append(simplified.BusinessLines, struct {
							Code        int    `json:"code"`
							Description string `json:"description"`
						}{
							Code:        int(subgroup["Code"].(float64)),
							Description: subgroup["Description"].(string),
						})
					}
				}
			}
		}
	}

	SendJSONResponse(w, SimplifiedResponse{
		Success: true,
		Data:    simplified,
	}, http.StatusOK)
}
