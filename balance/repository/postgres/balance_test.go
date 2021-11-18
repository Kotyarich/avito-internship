package postgres

import (
	"avito-intership/balance"
	"avito-intership/models"
	"avito-intership/utils"
	"database/sql"
	"github.com/ory/dockertest"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

const (
	smallAmount float32 = 100
	halfAmount float32 = 50
	bigAmount float32 = 1000
)

type balanceRepositorySuite struct {
	suite.Suite

	db *sql.DB
	pool *dockertest.Pool
	resource *dockertest.Resource

	repository balance.Repository
	curId int64
}


func (suite *balanceRepositorySuite) SetupSuite() {
	db, pool, resource := utils.DockerDBUp()
	err := utils.InitTable(db, "../../../init.sql")
	if err != nil {
		log.Fatal(err.Error())
	}

	repository := NewBalanceRepository(db)

	suite.repository = repository
	suite.db = db
	suite.pool = pool
	suite.resource = resource
	suite.curId = 0
}

func (suite *balanceRepositorySuite) TestGetHistory() {
	suite.curId += 1
	id := suite.curId
	var page int64 = 1
	var perPage int64 = 5
	sort := balance.SortDate
	desc := false
	var amount float32 = 1

	transactions := []*models.Transaction{
		{UserId:id, Amount:amount, TargetId:balance.RefillId, Type:"fill"},
	}

	err := suite.repository.ChangeBalance(id, amount, balance.RefillId)
	suite.NoError(err, "changing balance should not produce error")

	res, err := suite.repository.GetHistory(id, page, perPage, sort, desc)

	suite.NoError(err, "getting history should not produce error")
	suite.Equal(transactions[0].Type, res[0].Type)
	suite.Equal(transactions[0].TargetId, res[0].TargetId)
	suite.Equal(transactions[0].Amount, res[0].Amount)
	suite.Equal(transactions[0].UserId, res[0].UserId)
	suite.Equal(len(transactions), len(res))
}

func (suite *balanceRepositorySuite) TestGetBalance_Zero() {
	var expectedAmount float32 = 0
	suite.curId += 1
	id := suite.curId

	amount, err := suite.repository.GetBalance(id)

	suite.NoError(err, "getting balance should not produce error")
	suite.Equal(expectedAmount, amount)
}

func (suite *balanceRepositorySuite) TestChangeBalance_NewBalance() {
	amount := smallAmount
	suite.curId += 1
	id := suite.curId
	var product int64 = 1

	err := suite.repository.ChangeBalance(id, amount, product)

	suite.NoError(err, "changing balance should not produce error")
}

func (suite *balanceRepositorySuite) TestChangeBalance_TooLow() {
	amount  := -bigAmount
	suite.curId += 1
	id := suite.curId
	var product int64 = 1

	err := suite.repository.ChangeBalance(id, smallAmount, product)
	suite.NoError(err, "positive changing balance should not produce error")

	err = suite.repository.ChangeBalance(id, amount, product)
	suite.EqualError(balance.ErrTooLowBalance, err.Error())
}

func (suite *balanceRepositorySuite) TestGetBalance_NonZero() {
	expectedAmount := smallAmount
	suite.curId += 1
	id := suite.curId
	var product int64 = 1

	err := suite.repository.ChangeBalance(id, smallAmount, product)
	suite.NoError(err, "positive changing balance should not produce error")

	amount, err := suite.repository.GetBalance(id)

	suite.NoError(err, "getting balance should not produce error")
	suite.Equal(expectedAmount, amount)
}

func (suite *balanceRepositorySuite) TestChangeBalance_Withdraw() {
	amount := halfAmount
	suite.curId += 1
	srcId := suite.curId
	suite.curId += 1
	dstId := suite.curId
	var product int64 = 1

	err := suite.repository.ChangeBalance(srcId, smallAmount, product)
	suite.NoError(err, "positive changing balance should not produce error")

	err = suite.repository.TransferMoney(srcId, dstId, amount)
	suite.NoError(err, "changing balance should not produce error")

	newAmount, err := suite.repository.GetBalance(srcId)

	suite.NoError(err, "getting balance should not produce error")
	suite.Equal(amount, newAmount)
}

func (suite *balanceRepositorySuite) TestTransferMoney_TooLowBalance() {
	amount := halfAmount
	suite.curId += 1
	srcId := suite.curId
	suite.curId += 1
	dstId := suite.curId

	err := suite.repository.TransferMoney(srcId, dstId, amount)
	suite.Equal(balance.ErrTooLowBalance, err)
}


func (suite *balanceRepositorySuite) TearDownSuite() {
	err := utils.DropTable(suite.db, []string{"balances"})
	if err != nil {
		log.Println(err)
	}

	if err := suite.pool.Purge(suite.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}

func TestBalanceUseCase(t *testing.T) {
	suite.Run(t, new(balanceRepositorySuite))
}
