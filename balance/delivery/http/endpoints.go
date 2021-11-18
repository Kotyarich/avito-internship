package http

import (
	"avito-intership/balance"
	"avito-intership/exchange"
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
	Error  *string `json:"error"`
}

type StatusMessage struct {
	Success bool    `json:"success"`
	Message *string `json:"message"`
}

func (h Handler) writeStatus(success bool, message *string, w *http.ResponseWriter) {
	status := StatusMessage{
		Success: success,
		Message: message,
	}

	(*w).Header().Add("Content-Type", "application/json")
	_ = json.NewEncoder(*w).Encode(status)
}

func (h Handler) GetBalanceEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad id argument"
		h.writeStatus(false, &message, &w)
		return
	}

	currency := r.FormValue("currency")
	if currency == "" {
		currency = exchange.RUB
	}

	amount, err := h.useCase.GetBalance(id, currency)

	balanceResponse := Balance{id, amount, nil}
	if err == balance.ErrConversion {
		errMessage := err.Error()
		balanceResponse.Error = &errMessage
	} else if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		message := "Server error"
		h.writeStatus(false, &message, &w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(balanceResponse)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		message := "Server error"
		h.writeStatus(false, &message, &w)
	}
}

func (h Handler) ChangeBalanceEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad id argument"
		h.writeStatus(false, &message, &w)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad amount argument"
		h.writeStatus(false, &message, &w)
		return
	}

	product := balance.RefillId
	if amount < 0 {
		product, err = strconv.ParseInt(r.FormValue("product"), 10, 32)
		if err != nil || product < 0 {
			log.Println(err.Error())
			w.WriteHeader(http.StatusBadRequest)
			message := "Bad product argument"
			h.writeStatus(false, &message, &w)
			return
		}
	}

	err = h.useCase.ChangeBalance(id, float32(amount), product)
	if err == balance.ErrTooLowBalance {
		log.Println(err.Error())
		w.WriteHeader(http.StatusConflict)
		message := err.Error()
		h.writeStatus(false, &message, &w)
	} else if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		message := "Server error"
		h.writeStatus(false, &message, &w)
	} else {
		w.WriteHeader(http.StatusOK)
		h.writeStatus(true, nil, &w)
	}
}

func (h Handler) TransferMoneyEndpoint(w http.ResponseWriter, r *http.Request) {
	srcId, err := strconv.ParseInt(r.FormValue("src"), 10, 64)
	if err != nil || srcId <= 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad src argument"
		h.writeStatus(false, &message, &w)
		return
	}

	dstId, err := strconv.ParseInt(r.FormValue("dst"), 10, 64)
	if err != nil || dstId <= 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad dst argument"
		h.writeStatus(false, &message, &w)
		return
	}

	amount, err := strconv.ParseFloat(r.FormValue("amount"), 32)
	if err != nil || amount <= 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad amount argument"
		h.writeStatus(false, &message, &w)
		return
	}

	err = h.useCase.TransferMoney(srcId, dstId, float32(amount))
	if err == balance.ErrTooLowBalance {
		log.Println(err.Error())
		w.WriteHeader(http.StatusConflict)
		message := err.Error()
		h.writeStatus(false, &message, &w)
	} else if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		message := "Server error"
		h.writeStatus(false, &message, &w)
	} else {
		h.writeStatus(true, nil, &w)
	}
}

func (h Handler) GetHistoryEndpoint(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil || id <= 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad id argument"
		h.writeStatus(false, &message, &w)
		return
	}

	page, err := strconv.ParseInt(r.FormValue("page"), 10, 64)
	if err != nil || page <= 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad page argument"
		h.writeStatus(false, &message, &w)
		return
	}

	perPage, err := strconv.ParseInt(r.FormValue("per_page"), 10, 64)
	if err != nil || perPage <= 0 {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		message := "Bad per_page argument"
		h.writeStatus(false, &message, &w)
		return
	}

	sortValue := r.FormValue("sort")
	sort := balance.SortDate
	if sortValue == "amount" {
		sort = balance.SortAmount
	}

	descValue := r.FormValue("desc")
	desc := false
	if descValue == "true" {
		desc = true
	}

	transactions, err := h.useCase.GetHistory(id, page, perPage, sort, desc)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		message := "Server error"
		h.writeStatus(false, &message, &w)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(transactions)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		message := "Server error"
		h.writeStatus(false, &message, &w)
	}
}
