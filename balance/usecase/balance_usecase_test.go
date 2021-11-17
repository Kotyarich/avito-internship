package usecase

import (
	"avito-intership/balance"
	"avito-intership/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type balanceUseCaseSuite struct {
	suite.Suite
	repository *mocks.Repository
	useCase    balance.UseCase
}

func (suite *balanceUseCaseSuite) SetupTest() {
	repository := new(mocks.Repository)
	useCase := NewBalanceUseCase(repository)

	suite.repository = repository
	suite.useCase = useCase
}

func (suite *balanceUseCaseSuite) TestGetBalance_RUB() {
	var id int64 = 1
	var amount float32 = 100
	currency := "RUB"

	suite.repository.On("GetBalance", id).Return(amount, nil)

	result, err := suite.useCase.GetBalance(id, currency)

	suite.Nil(err, "no error when return amount")
	suite.Equal(amount, result, "result and amount should be equal")
}

func (suite *balanceUseCaseSuite) TestGetBalance_USD() {
	var id int64 = 1
	var amount float32 = 100
	currency := "USD"

	suite.repository.On("GetBalance", id).Return(amount, nil)

	result, err := suite.useCase.GetBalance(id, currency)

	suite.Nil(err, "no error when return amount")
	suite.Equal(amount, result, "result and amount should be equal")
}

func (suite *balanceUseCaseSuite) TestChangeBalance_Add() {
	var id int64 = 1
	var amount float32 = 10

	suite.repository.On("ChangeBalance", id, amount).Return(nil)

	err := suite.useCase.ChangeBalance(id, amount)

	suite.Nil(err, "no error when changing balance")
}

func (suite *balanceUseCaseSuite) TestChangeBalance_Withdraw() {
	var id int64 = 1
	var amount float32 = -10

	suite.repository.On("ChangeBalance", id, amount).Return(nil)

	err := suite.useCase.ChangeBalance(id, amount)

	suite.Nil(err, "no error when changing balance")
}

func (suite *balanceUseCaseSuite) TestChangeBalance_TooLowBalance() {
	var id int64 = 1
	var amount float32 = 10

	suite.repository.On("ChangeBalance", id, amount).Return(balance.ErrTooLowBalance)

	err := suite.useCase.ChangeBalance(id, amount)

	suite.Equal(balance.ErrTooLowBalance, err, "too low balance error expected")
}

func (suite *balanceUseCaseSuite) TestTransferBalance_Positive() {
	var src int64 = 1
	var dst int64 = 2
	var amount float32 = 10

	suite.repository.On("TransferMoney", src, dst, amount).Return(nil)

	err := suite.useCase.TransferMoney(src, dst, amount)

	suite.Nil(err, "no error during transfer expected")
}

func (suite *balanceUseCaseSuite) TestTransferBalance_TooLowBalance() {
	var src int64 = 1
	var dst int64 = 2
	var amount float32 = 10

	suite.repository.On("TransferMoney", src, dst, amount).Return(balance.ErrTooLowBalance)

	err := suite.useCase.TransferMoney(src, dst, amount)

	suite.Equal(balance.ErrTooLowBalance, err, "too low balance error expected")
}

func TestBalanceUseCase(t *testing.T) {
	suite.Run(t, new(balanceUseCaseSuite))
}
