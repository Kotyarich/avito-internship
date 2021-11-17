package postgres

import (
	"avito-intership/balance"
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
	err := utils.InitTable(db, "../../../scripts/init.sql")
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

	err := suite.repository.ChangeBalance(id, amount)

	suite.NoError(err, "changing balance should not produce error")
}

func (suite *balanceRepositorySuite) TestChangeBalance_TooLow() {
	amount  := -bigAmount
	suite.curId += 1
	id := suite.curId

	err := suite.repository.ChangeBalance(id, smallAmount)
	suite.NoError(err, "positive changing balance should not produce error")

	err = suite.repository.ChangeBalance(id, amount)
	suite.EqualError(balance.ErrTooLowBalance, err.Error())
}

func (suite *balanceRepositorySuite) TestGetBalance_NonZero() {
	expectedAmount := smallAmount
	suite.curId += 1
	id := suite.curId

	err := suite.repository.ChangeBalance(id, smallAmount)
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

	err := suite.repository.ChangeBalance(srcId, smallAmount)
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
