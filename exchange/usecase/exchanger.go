package usecase

import (
	"avito-intership/exchange"
	"math"
)

type Exchanger struct {
	repository exchange.RateRepository
}

func NewExchanger(netRepository exchange.RateRepository) *Exchanger {
	return &Exchanger{
		repository: netRepository,
	}
}

func (e *Exchanger) ConvertRubles(amount float32, currency string) (float32, error) {
	rate, err := e.repository.GetRubleRate(currency)
	if err != nil {
		return 0, err
	}
	// Округляем до 2 знаков после запятой
	converted := float32(math.Round(float64(amount/rate*100)) / 100)

	return converted, nil
}
