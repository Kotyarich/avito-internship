package usecase

import "avito-intership/balance"

type BalanceUseCase struct {
	balanceRepo balance.Repository
}

func NewBalanceUseCase(repo balance.Repository) *BalanceUseCase {
	return &BalanceUseCase{balanceRepo:repo}
}

func (u BalanceUseCase) GetBalance(userId int64, currency string) (float32, error) {
	amount, err := u.balanceRepo.GetBalance(userId)

	return amount, err
}

func (u BalanceUseCase) ChangeBalance(userId int64, amount float32) error {
	err := u.balanceRepo.ChangeBalance(userId, amount)
	return err
}

func (u BalanceUseCase) TransferMoney(srcUserId int64, dstUserId int64, amount float32) error {
	err := u.balanceRepo.TransferMoney(srcUserId, dstUserId, amount)
	return err
}
