package balance

import "errors"

var (
	ErrTooLowBalance = errors.New("balance can't be lower than 0")
	ErrConversion    = errors.New("conversion wasn't completed, amount returned in RUB")
)
