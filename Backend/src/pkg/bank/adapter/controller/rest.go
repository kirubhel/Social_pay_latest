package controller

import (
	"log"
	"net/http"

	auth "github.com/socialpay/socialpay/src/pkg/auth/adapter/controller/procedure"
)

type Controller struct {
	log        *log.Logger
	sm         *http.ServeMux
	auth       auth.Controller
	
}

func NewAwashTest(log *log.Logger, sm *http.ServeMux,auth auth.Controller) Controller {

	con:=Controller{log:log,auth:auth}
	
	sm.HandleFunc("/api/v1/bank/awash-debit",func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodPost :
			con.TestDebitHandler(w,r)
		}
	
	})


	sm.HandleFunc("/api/v1/bank/awash/rs-status",func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodPost :
			con.TestDebitStatus(w,r)
		}
	
	})


	sm.HandleFunc("/api/v1/bank/awash/success",func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodPost :
			con.SuccessCallBack(w,r)
		}
	
	})




	con.sm=sm
	return con
}