package exchange

const RUB string = "RUB"

type Exchanger interface {
	ConvertRubles(amount float32, currency string) (float32, error)
}
