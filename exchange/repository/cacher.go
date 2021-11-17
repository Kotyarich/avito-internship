package repository

type Cacher interface {
	GetRubleRate(currency string) (float32, error)
	SetRate(currency string, rate float32) error
}
