package controller

import (
	"bytes"
	"encoding/json"
	"enigmacamp.com/golang-sample/model"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

var dummyCustomers = []model.Customer{
	{
		Id:      "C001",
		Nama:    "Dummy Name 1",
		Address: "Dummy Address 1",
	},
}

// controller buth usecase
type CustomerUseCaseMock struct {
	mock.Mock
}

// buat TestSuite
type CustomerControllerTestSuit struct {
	suite.Suite
	routerMock  *gin.Engine
	useCaseMock *CustomerUseCaseMock
}

func (suite *CustomerControllerTestSuit) SetupTest() {
	suite.routerMock = gin.Default()
	suite.useCaseMock = new(CustomerUseCaseMock)
}

func (c *CustomerUseCaseMock) RegisterCustomer(customer model.Customer) error {
	args := c.Called(customer)
	if args.Get(0) != nil {
		return args.Get(0).(error)
	}
	return nil
}

func (c *CustomerUseCaseMock) FindCustomerById(id string) (model.Customer, error) {
	args := c.Called(id)
	if args.Get(1) != nil {
		return model.Customer{}, args.Get(1).(error)
	}
	return args.Get(0).(model.Customer), nil
}

func (c *CustomerUseCaseMock) GetAllCustomer() ([]model.Customer, error) {
	args := c.Called()
	if args.Get(1) != nil {
		return nil, args.Get(1).(error)
	}
	return args.Get(0).([]model.Customer), nil
}

func (suite *CustomerControllerTestSuit) TestGetAllCustomerApi_Success() {
	customers := dummyCustomers
	suite.useCaseMock.On("GetAllCustomer").Return(customers, nil)
	NewCustomerController(suite.routerMock, suite.useCaseMock)

	//buat kondisi untuk HTTP status
	r := httptest.NewRecorder()
	//reques test yang sesuai
	request, err := http.NewRequest(http.MethodGet, "/customer", nil)
	suite.routerMock.ServeHTTP(r, request)
	var actualCustomers []model.Customer
	response := r.Body.String()
	json.Unmarshal([]byte(response), &actualCustomers)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	assert.Equal(suite.T(), 1, len(actualCustomers))
	assert.Equal(suite.T(), customers[0].Nama, actualCustomers[0].Nama)
	assert.Nil(suite.T(), err)
}

func (suite *CustomerControllerTestSuit) TestGetAllCustomerApi_Failed() {
	suite.useCaseMock.On("GetAllCustomer").Return(nil, errors.New("failed"))
	NewCustomerController(suite.routerMock, suite.useCaseMock)

	//buat kondisi untuk HTTP status
	r := httptest.NewRecorder()
	request, err := http.NewRequest(http.MethodGet, "/customer", nil)
	suite.routerMock.ServeHTTP(r, request)
	response := r.Body.String()
	var errorResponse struct{ Err string }
	json.Unmarshal([]byte(response), &errorResponse)
	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
	assert.Equal(suite.T(), "failed", errorResponse.Err)
	assert.Nil(suite.T(), err)
}

func (suite *CustomerControllerTestSuit) TestRegisterCustomerApi_Success() {
	dummyCustomer := dummyCustomers[0]
	suite.useCaseMock.On("RegisterCustomer", dummyCustomer).Return(nil)
	NewCustomerController(suite.routerMock, suite.useCaseMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(reqBody))
	suite.routerMock.ServeHTTP(r, request)
	response := r.Body.String()
	var actualCustomer model.Customer
	json.Unmarshal([]byte(response), &actualCustomer)
	assert.Equal(suite.T(), http.StatusOK, r.Code)
	assert.Equal(suite.T(), dummyCustomer.Nama, actualCustomer.Nama)
}

func (suite *CustomerControllerTestSuit) TestRegisterCustomerApi_FailedBinding() {
	r := httptest.NewRecorder()
	NewCustomerController(suite.routerMock, suite.useCaseMock)
	request, _ := http.NewRequest(http.MethodPost, "/customer", nil)
	suite.routerMock.ServeHTTP(r, request)
	var errprResponse struct{ Err string }
	response := r.Body.String()
	json.Unmarshal([]byte(response), &errprResponse)
	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
	assert.NotEmpty(suite.T(), errprResponse.Err)

}

func (suite *CustomerControllerTestSuit) TestRegisterCustomerApi_FailedUseCase() {
	dummyCustomer := dummyCustomers[0]
	suite.useCaseMock.On("RegisterCustomer", dummyCustomer).Return(errors.New("failed"))
	NewCustomerController(suite.routerMock, suite.useCaseMock)
	r := httptest.NewRecorder()
	reqBody, _ := json.Marshal(dummyCustomer)
	request, _ := http.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(reqBody))
	suite.routerMock.ServeHTTP(r, request)
	response := r.Body.String()
	var errorResponse struct{ Err string }
	json.Unmarshal([]byte(response), &errorResponse)
	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
	assert.Equal(suite.T(), "failed", errorResponse.Err)
}

func TestCustomerControllerTestSuite(t *testing.T) {
	suite.Run(t, new(CustomerControllerTestSuit))
}
