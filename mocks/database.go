package mocks

import (
	"github.com/pashagolub/pgxmock"
	"log"
)

func NewDBMock() pgxmock.PgxPoolIface {
	mock, err := pgxmock.NewPool()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	return mock
}
