package repository

import (
	"database/sql"
	"enigmacamp.com/golang-sample/model"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"log"
	"testing"
)

var dummyCustomers = []model.Customer{
	{
		Id:      "C001",
		Nama:    "Dummy Name 1",
		Address: "Dummy Address 1",
	},
	{
		Id:      "C002",
		Nama:    "Dummy Name 2",
		Address: "Dummy Address 2",
	},
}

type CustomerRepositoryTestSuite struct {
	suite.Suite
	mockDb  *sql.DB
	mockSql sqlmock.Sqlmock
}

func (suite *CustomerRepositoryTestSuite) SetupTest() {
	mockDb, mockSql, err := sqlmock.New()
	if err != nil {
		log.Fatalln("An Error When opening a stub database connection", err)
	}

	suite.mockDb = mockDb
	suite.mockSql = mockSql
}

func (suite *CustomerRepositoryTestSuite) TearDownTest() {
	suite.mockDb.Close()
}

func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Success() {
	//siapakan colum
	rows := sqlmock.NewRows([]string{"id", "name", "address"})
	for _, v := range dummyCustomers {
		rows.AddRow(v.Id, v.Nama, v.Address)
	}

	//buat sebuah query mocnya menggunakan regex
	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnRows(rows)

	//panggil repo aslinya
	repo := NewCustomerDbRepository(suite.mockDb)

	//panggil method yang mau di test
	actual, err := repo.RetrieveAll()

	//buat test assertion
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 2, len(actual))
	assert.Equal(suite.T(), "C001", actual[0].Id)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerRetrieveAll_Failed() {
	//siapakan colum
	rows := sqlmock.NewRows([]string{"ids", "names", "addressssss"})
	for _, v := range dummyCustomers {
		rows.AddRow(v.Id, v.Nama, v.Address)
	}

	//buat sebuah query mocnya menggunakan regex
	suite.mockSql.ExpectQuery("select \\* from customer").WillReturnError(errors.New("failed"))

	//panggil repo aslinya
	repo := NewCustomerDbRepository(suite.mockDb)

	//panggil method yang mau di test
	actual, err := repo.RetrieveAll()

	//buat test assertion
	assert.Nil(suite.T(), actual)
	assert.Error(suite.T(), err)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerCreate_Success() {
	dummyCustomer := dummyCustomers[0]
	suite.mockSql.ExpectExec("insert into customer values").WithArgs(dummyCustomers[0].Id, dummyCustomers[0].Nama, dummyCustomers[0].Address).WillReturnResult(sqlmock.NewResult(1, 1))
	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer)
	assert.Nil(suite.T(), err)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerCreate_Failed() {
	dummyCustomer := dummyCustomers[0]
	suite.mockSql.ExpectExec("insert into customer values").WillReturnError(errors.New("failed"))
	repo := NewCustomerDbRepository(suite.mockDb)
	err := repo.Create(dummyCustomer)

	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), errors.New("failed"), err)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerFindById_Success() {
	dummyCustomer := dummyCustomers[0]
	rows := sqlmock.NewRows([]string{"id", "name", "address"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)
	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnRows(rows)

	repo := NewCustomerDbRepository(suite.mockDb)
	actual, err := repo.FindById(dummyCustomer.Id)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), actual)
}

func (suite *CustomerRepositoryTestSuite) TestCustomerFindById_Failed() {
	dummyCustomer := dummyCustomers[0]
	rows := sqlmock.NewRows([]string{"ids", "names", "addresssss"})
	rows.AddRow(dummyCustomer.Id, dummyCustomer.Nama, dummyCustomer.Address)
	suite.mockSql.ExpectQuery("select \\* from customer where id").WillReturnError(errors.New("failed"))

	repo := NewCustomerDbRepository(suite.mockDb)
	actual, err := repo.FindById(dummyCustomer.Id)

	//buat test assertion
	//func() {
	//	defer func() {
	//		if r := recover(); r == nil {
	//			assert.Error(suite.T(), err)
	//		}
	//	}()
	//	repo.FindById(dummyCustomer.Id)
	//}()
	assert.NotEqual(suite.T(), dummyCustomer, actual)
	assert.Error(suite.T(), err)
}

func TestCustomerRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerRepositoryTestSuite))
}
