package http

import (
	"avito-intership/balance"
	"github.com/gorilla/mux"
	"net/http"
)

func RegisterEndpoints(router *mux.Router, uc balance.UseCase) {
	handler := NewHandler(uc)

	router.HandleFunc("/api/v1/balance/{id:[0-9]+}", handler.GetBalanceEndpoint).
		Methods(http.MethodOptions, http.MethodGet)
	router.HandleFunc("/api/v1/balance/{id:[0-9]+}", handler.ChangeBalanceEndpoint).
		Methods(http.MethodOptions, http.MethodPost)
	router.HandleFunc("/api/v1/transfer", handler.TransferMoneyEndpoint).
		Methods(http.MethodOptions, http.MethodPost)
}

