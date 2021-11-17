package exchange

type RateRepository interface {
	GetRubleRate(currency string) (float32, error)
}
