package balance

import "avito-intership/models"

const RefillId int64 = 0

const (
	WithdrawType string = "product"
	TransferType string = "transfer"
	RefillType   string = "fill"
)

const (
	SortAmount = iota
	SortDate   = iota
)

type Repository interface {
	ChangeBalance(userId int64, amount float32, productId int64) error
	GetBalance(userId int64) (float32, error)
	TransferMoney(srcUserId int64, dstUserId int64, amount float32) error
	GetHistory(userId int64, page int64, perPage int64, sort int, desc bool) ([]*models.Transaction, error)
}
