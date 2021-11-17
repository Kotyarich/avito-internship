package usecase

import (
	"avito-intership/exchange"
	"avito-intership/mocks"
	"github.com/stretchr/testify/suite"
	"testing"
)

type exchangeUseCaseSuite struct {
	suite.Suite
	repository *mocks.RateRepository
	useCase    exchange.Exchanger
}

func (suite *exchangeUseCaseSuite) SetupTest() {
	repository := new(mocks.RateRepository)
	useCase := NewExchanger(repository)

	suite.repository = repository
	suite.useCase = useCase
}

func (suite *exchangeUseCaseSuite) TestGetBalance_RUB() {
	var amount float32 = 100
	var rate float32 = 100

	currency := "USD"

	suite.repository.On("GetRubleRate", currency).Return(rate, nil)

	result, err := suite.useCase.ConvertRubles(amount, currency)

	suite.Nil(err, "no error while converting")
	suite.Equal(amount / rate, result, )
}

func TestBalanceUseCase(t *testing.T) {
	suite.Run(t, new(exchangeUseCaseSuite))
}
