package balance

type UseCase interface {
	ChangeBalance(userId int64, amount float32) error
	GetBalance(userId int64, currency string) (float32, error)
	TransferMoney(srcUserId int64, dstUserId int64, amount float32) error
}