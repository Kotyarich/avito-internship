package http

import (
	"avito-intership/balance"
	"avito-intership/mocks"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type balanceHandlerSuite struct {
	suite.Suite

	useCase       *mocks.UseCase
	handler       *Handler
	testingServer *httptest.Server
}

func (suite *balanceHandlerSuite) SetupSuite() {
	useCase := new(mocks.UseCase)
	handler := NewHandler(useCase)

	router := mux.NewRouter()
	RegisterEndpoints(router, useCase)
	testingServer := httptest.NewServer(router)

	suite.testingServer = testingServer
	suite.useCase = useCase
	suite.handler = handler
}

func (suite *balanceHandlerSuite) TearDownSuite() {
	defer suite.testingServer.Close()
}

func (suite *balanceHandlerSuite) TestGetBalanceHandler_Ok() {
	var id int64 = 1
	var amount float32 = 100
	currency := "RUB"

	suite.useCase.On("GetBalance", id, currency).Return(amount, nil)

	response, err := http.Get(fmt.Sprintf("%s/api/v1/balance/%d?currency=%s",
		suite.testingServer.URL, id, currency))
	suite.NoError(err, "request should not produce error")
	defer response.Body.Close()

	var responseBody Balance
	err = json.NewDecoder(response.Body).Decode(&responseBody)
	suite.NoError(err, "decoding should not produce error")

	suite.Equal(http.StatusOK, response.StatusCode)
	suite.Equal(responseBody.Amount, amount)
	suite.Equal(responseBody.Id, id)
}

func (suite *balanceHandlerSuite) TestChangeBalanceHandler_Ok() {
	var id int64 = 1
	var amount float32 = 100
	var product int64 = 0

	suite.useCase.On("ChangeBalance", id, amount, product).Return(nil)

	response, err := http.Post(fmt.Sprintf("%s/api/v1/balance/%d?amount=%f",
		suite.testingServer.URL, id, amount), "", bytes.NewBuffer([]byte{}))
	suite.NoError(err, "request should not produce error")
	defer response.Body.Close()

	suite.Equal(http.StatusOK, response.StatusCode)
}

func (suite *balanceHandlerSuite) TestChangeBalanceHandler_LowBalance() {
	var id int64 = 1
	var amount float32 = -100
	var product int64 = 1

	suite.useCase.On("ChangeBalance", id, amount, product).Return(balance.ErrTooLowBalance)

	response, err := http.Post(fmt.Sprintf("%s/api/v1/balance/%d?amount=%f&product=%d",
		suite.testingServer.URL, id, amount, product), "", bytes.NewBuffer([]byte{}))
	suite.NoError(err, "request should not produce error")
	defer response.Body.Close()

	suite.Equal(http.StatusConflict, response.StatusCode)
}

func (suite *balanceHandlerSuite) TestTransferMoneyHandler_Ok() {
	var src int64 = 1
	var dst int64 = 2
	var amount float32 = 10

	suite.useCase.On("TransferMoney", src, dst, amount).Return(nil)

	response, err := http.Post(fmt.Sprintf("%s/api/v1/transfer?src=%d&dst=%d&amount=%f",
		suite.testingServer.URL, src, dst, amount), "", bytes.NewBuffer([]byte{}))
	suite.NoError(err, "request should not produce error")
	defer response.Body.Close()

	suite.Equal(http.StatusOK, response.StatusCode)
}

func (suite *balanceHandlerSuite) TestTransferMoneyHandler_LowBalance() {
	var src int64 = 1
	var dst int64 = 2
	var amount float32 = 100

	suite.useCase.On("TransferMoney", src, dst, amount).Return(balance.ErrTooLowBalance)

	response, err := http.Post(fmt.Sprintf("%s/api/v1/transfer?src=%d&dst=%d&amount=%f",
		suite.testingServer.URL, src, dst, amount), "", bytes.NewBuffer([]byte{}))
	suite.NoError(err, "request should not produce error")
	defer response.Body.Close()

	suite.Equal(http.StatusConflict, response.StatusCode)
}

func TestBalanceHandler(t *testing.T) {
	suite.Run(t, new(balanceHandlerSuite))
}