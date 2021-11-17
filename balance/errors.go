package balance

import "errors"

var (
	ErrTooLowBalance = errors.New("balance can't be lower than 0")
)
