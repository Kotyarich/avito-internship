package http

import (
	"avito-intership/balance"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

type Handler struct {
	useCase balance.UseCase
}

func NewHandler(useCase balance.UseCase) *Handler {
	return &Handler{
		useCase: useCase,
	}
}

type Balance struct {
	Id     int64   `json:"id"`
	Amount float32 `json:"amount"`
}

func (h Handler) GetBalanceEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id < 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	currency := r.FormValue("currency")
	if currency == "" {
		currency = "RUB"
	}

	amount, err := h.useCase.GetBalance(id, currency)

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Balance{id, amount})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h Handler) ChangeBalanceEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id < 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.ChangeBalance(id, float32(amount))
	if err == balance.ErrTooLowBalance {
		log.Println(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}

func (h Handler) TransferMoneyEndpoint(w http.ResponseWriter, r *http.Request) {
	srcId, err := strconv.ParseInt(r.FormValue("src"), 10, 64)
	if err != nil || srcId < 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	dstId, err := strconv.ParseInt(r.FormValue("dst"), 10, 64)
	if err != nil || dstId < 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
	if err != nil || amount < 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = h.useCase.TransferMoney(srcId, dstId, float32(amount))
	if err == balance.ErrTooLowBalance {
		log.Println(err.Error())
		w.WriteHeader(http.StatusConflict)
	} else if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}
