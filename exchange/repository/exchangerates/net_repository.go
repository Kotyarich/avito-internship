package exchangerates

import (
	"avito-intership/exchange"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	baseUrl        string = "http://api.exchangeratesapi.io/v1/"
	latestEndpoint string = "latest"
	apiKeyEnd      string = "EXCHANGE_KEY"
)

type netRepository struct {
	apiKey string
}

func NewNetRepository() *netRepository {
	apiKey := os.Getenv(apiKeyEnd)

	return &netRepository{
		apiKey: apiKey,
	}
}

type apiError struct {
	Code int16  `json:"code"`
	Info string `json:"info"`
}

type apiResponse struct {
	Success   bool               `json:"success"`
	Timestamp uint64             `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float32 `json:"rates"`
	Error     apiError           `json:"error"`
}

func (r *netRepository) GetRubleRate(currency string) (float32, error) {
	queryString := fmt.Sprintf("?access_key=%s&symbols=%s,%s", r.apiKey, exchange.RUB, currency)
	response, err := http.Get(baseUrl + latestEndpoint + queryString)
	if err != nil {
		return 0, err
	}
	defer response.Body.Close()

	var responseBody apiResponse
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	if err != nil {
		return 0, err
	}

	if !responseBody.Success {
		return 0, fmt.Errorf(responseBody.Error.Info)
	}
	// Базовый тариф exhangerateapi не позволяет указать базовую валюту для получения курса,
	// поэтому для получения курса получаются курс евро к рублю и курс требуемой валюты к евро
	rublesInEur := responseBody.Rates[exchange.RUB]
	eurCurrencyRate := responseBody.Rates[currency]

	rubblesCurrencyRate := rublesInEur / eurCurrencyRate

	return rubblesCurrencyRate, nil
}
