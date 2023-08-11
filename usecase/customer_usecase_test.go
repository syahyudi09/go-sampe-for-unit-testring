package usecase

import (
	"enigmacamp.com/golang-sample/model"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
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

type repoMock struct {
	mock.Mock
}

func (r *repoMock) Create(newCustomer model.Customer) error {
	args := r.Called(newCustomer)
	if args.Get(0) != nil {
		return args.Error(0)
	}
	return nil
}

func (r *repoMock) RetrieveAll() ([]model.Customer, error) {
	args := r.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Customer), nil
}

func (r *repoMock) FindById(id string) (model.Customer, error) {
	args := r.Called(id)
	if args.Get(1) != nil {
		return model.Customer{}, args.Error(1)
	}
	return args.Get(0).(model.Customer), nil
}

func (suite *CustomerUseCaseTestSuite) TestCustomerFIndById_Success() {
	dummyCustomer := dummyCustomers[0]
	suite.repoMock.On("FindById", dummyCustomer.Id).Return(dummyCustomer, nil)
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.FindCustomerById(dummyCustomer.Id)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), dummyCustomer.Id, customer.Id)
}

func (suite *CustomerUseCaseTestSuite) TestCustomerFIndById_Failed() {
	dummyCustomer := dummyCustomers[0]
	suite.repoMock.On("FindById", dummyCustomer.Id).Return(model.Customer{}, errors.New("failed"))
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	customer, err := customerUseCaseTest.FindCustomerById(dummyCustomer.Id)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "failed", err.Error())
	assert.Equal(suite.T(), "", customer.Id)
}

func (suite *CustomerUseCaseTestSuite) TestCustomerCreate_Success() {
	dummyCustomer := dummyCustomers[0]
	suite.repoMock.On("Create", dummyCustomer).Return(nil)
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	err := customerUseCaseTest.RegisterCustomer(dummyCustomer)
	assert.Nil(suite.T(), err)
}

func (suite *CustomerUseCaseTestSuite) TestCustomerCreate_Failed() {
	dummyCustomer := dummyCustomers[0]
	suite.repoMock.On("Create", dummyCustomer).Return(errors.New("failed"))
	customerUseCaseTest := NewCustomerUseCase(suite.repoMock)
	err := customerUseCaseTest.RegisterCustomer(dummyCustomer)
	assert.NotNil(suite.T(), err)
}

type CustomerUseCaseTestSuite struct {
	suite.Suite
	repoMock *repoMock
}

func (suite *CustomerUseCaseTestSuite) SetupTest() {
	suite.repoMock = new(repoMock)
}

func TestCustomerUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerUseCaseTestSuite))
}
