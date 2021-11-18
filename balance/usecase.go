package balance

import "avito-intership/models"

type UseCase interface {
	ChangeBalance(userId int64, amount float32, productId int64) error
	GetBalance(userId int64, currency string) (float32, error)
	TransferMoney(srcUserId int64, dstUserId int64, amount float32) error
	GetHistory(userId int64, page int64, perPage int64, sort int, desc bool) ([]*models.Transaction, error)
}