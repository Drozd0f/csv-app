package service

import (
	"context"
	"encoding/csv"
	"io"
	"log"
	"testing"

	"github.com/Drozd0f/csv-app/pkg/migrator"
	"github.com/Drozd0f/csv-app/test/fixtures"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"

	"github.com/Drozd0f/csv-app/conf"
	"github.com/Drozd0f/csv-app/repository"
	"github.com/Drozd0f/csv-app/schemes"
	"github.com/Drozd0f/csv-app/test"
	"github.com/Drozd0f/csv-app/test/containers"
	"github.com/stretchr/testify/suite"
)

const (
	countTestTransactions = 50
	countFromThirteenth   = 48 // count transactions in file test_transactions.csv by DatePost
	countToThirteenth     = 2
)

type ServiceTestSuite struct {
	suite.Suite
	ctx         context.Context
	dbContainer *containers.TestDatabase
	service     *Service
	repository  *repository.Repository
	conf        *conf.Config
}

func (ts *ServiceTestSuite) SetupSuite() {
	log.Println("setup test suite...")
	ts.ctx = context.Background()
	ts.dbContainer = containers.NewTestDatabase(ts.T())

	rep, err := repository.New(ts.ctx, ts.dbContainer.ConnectionString(ts.T()))
	ts.Require().NoError(err)

	ts.repository = rep

	testConf := test.NewConfig(ts.T(), ts.dbContainer)
	ts.conf = testConf

	ts.service = New(ts.repository, ts.conf)

	err = migrator.MakeMigrate(ts.conf)
	ts.Require().NoError(err)
}

func (ts *ServiceTestSuite) TearDownSuite() {
	log.Println("tear down test suite...")
	ts.dbContainer.Close(ts.T())
}

func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (ts *ServiceTestSuite) SeedDB() {
	file, err := fixtures.Fixtures.Open("transactions_test.csv")
	ts.Require().NoError(err)

	parser := csv.NewReader(file)
	headers, err := parser.Read()
	ts.Require().NoError(err)

	tranS := make([]schemes.Transaction, 0, countTestTransactions)

	for {
		row, err := parser.Read()
		if err == io.EOF {
			break
		}
		ts.Require().NoError(err)

		var t schemes.Transaction
		ts.Require().NoError(schemes.BindFromCsv(&t, row, headers))

		tranS = append(tranS, t)
	}

	err = ts.repository.Seed(ts.ctx, tranS)
	ts.Require().NoError(err, "repository seed:")
}

func (ts *ServiceTestSuite) CleanupDB() {
	err := ts.repository.Cleanup(ts.ctx)
	ts.Require().NoError(err)
}
