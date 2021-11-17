package usecase

import (
	"avito-intership/balance"
	"avito-intership/exchange"
	"log"
)

type BalanceUseCase struct {
	balanceRepo balance.Repository
	exchanger   exchange.Exchanger
}

func NewBalanceUseCase(repo balance.Repository, exchanger exchange.Exchanger) *BalanceUseCase {
	return &BalanceUseCase{
		balanceRepo: repo,
		exchanger:   exchanger,
	}
}

func (u BalanceUseCase) GetBalance(userId int64, currency string) (float32, error) {
	amount, err := u.balanceRepo.GetBalance(userId)
	if err != nil {
		return 0, err
	}

	if currency != exchange.RUB {
		converted, err := u.exchanger.ConvertRubles(amount, currency)
		if err != nil {
			// В случае ошибки конвертации возращаем пользователю баланс в рублях
			log.Println(err)
			return amount, balance.ErrConversion
		}

		return converted, nil
	}

	return amount, nil
}

func (u BalanceUseCase) ChangeBalance(userId int64, amount float32) error {
	err := u.balanceRepo.ChangeBalance(userId, amount)
	return err
}

func (u BalanceUseCase) TransferMoney(srcUserId int64, dstUserId int64, amount float32) error {
	err := u.balanceRepo.TransferMoney(srcUserId, dstUserId, amount)
	return err
}
