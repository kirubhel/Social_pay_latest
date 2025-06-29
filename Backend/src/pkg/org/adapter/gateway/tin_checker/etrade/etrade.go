// etrade/etrade.go
package etrade

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/socialpay/socialpay/src/pkg/org/core/entity"
	"github.com/socialpay/socialpay/src/pkg/org/usecase"
)

type Etrade struct {
	log     *log.Logger
	Referer string
	client  *http.Client
}

func New(log *log.Logger) usecase.TINChecker {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return Etrade{log: log, Referer: "https://etrade.gov.et/business-license-checker", client: client}
}

// Implement both interface methods
func (etrade Etrade) CheckTIN(tin string, uc usecase.Usecase) (*entity.Organization, error) {
	raw, err := etrade.CheckTINRaw(tin)
	if err != nil {
		return nil, err
	}

	// Convert raw response to Organization (simplified example)
	org := &entity.Organization{
		Name:    raw["BusinessName"].(string),
		Capital: raw["PaidUpCapital"].(float64),
		// ... other fields
	}
	return org, nil
}

func (etrade Etrade) CheckTINRaw(tin string) (map[string]interface{}, error) {
	etrade.log.Println("CHECKING TIN RAW")

	req, err := http.NewRequest(http.MethodGet,
		fmt.Sprintf("https://etrade.gov.et/api/Registration/GetRegistrationInfoByTin/%s/en", tin), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Referer", etrade.Referer)

	res, err := etrade.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("failed to check TIN")
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &responseData); err != nil {
		return nil, err
	}

	return responseData, nil
}
