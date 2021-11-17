package exchangerates

import (
	"avito-intership/exchange"
	"avito-intership/exchange/repository"
	"log"
)

type Repository struct {
	netRepo exchange.RateRepository
	cache   repository.Cacher
}

func NewExchangeRepository(netRepo exchange.RateRepository, cache repository.Cacher) *Repository {
	return &Repository{
		netRepo: netRepo,
		cache:   cache,
	}
}

func (r *Repository) GetRubleRate(currency string) (float32, error) {
	rate, err := r.cache.GetRubleRate(currency)
	if err == nil {
		return rate, nil
	}

	rate, err = r.netRepo.GetRubleRate(currency)
	if err != nil {
		return 0, err
	}

	err = r.cache.SetRate(currency, rate)
	if err != nil {
		log.Println(err)
	}

	return rate, nil
}
