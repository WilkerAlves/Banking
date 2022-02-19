package app

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/WilkerAlves/banking-lib/errs"
	"github.com/WilkerAlves/banking/dto"
	"github.com/WilkerAlves/banking/mocks/service"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
)

var router *mux.Router
var ch CustomerHandlers
var mockService *service.MockCustomerService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockService = service.NewMockCustomerService(ctrl)
	ch = CustomerHandlers{mockService}
	router = mux.NewRouter()
	router.HandleFunc("/customers", ch.getAllCustomers)
	return func() {
		router = nil
		defer ctrl.Finish()
	}
}

func Test_should_return_customers_with_status_code_200(t *testing.T) {
	// arrange
	teardown := setup(t)
	defer teardown()
	dummyCustomers := []dto.CustomerResponse{
		{Id: "1001", Name: "Ashish", City: "New Delhi", Zipcode: "110075", DateOfBirth: "2000-01-01", Status: "1"},
		{Id: "10012", Name: "Rob", City: "New Delhi", Zipcode: "110075", DateOfBirth: "2000-01-01", Status: "1"},
	}
	mockService.EXPECT().GetAllCustomer("").Return(dummyCustomers, nil)
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusOK {
		t.Error("Failed while testing the status code")
	}

}

func Test_should_return_status_code_500_with_error_message(t *testing.T) {
	// arrange
	teardown := setup(t)
	defer teardown()
	mockService.EXPECT().GetAllCustomer("").Return(
		nil,
		errs.NewUnexpectedError("some database error"),
	)
	request, _ := http.NewRequest(http.MethodGet, "/customers", nil)

	// Act
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	// Assert
	fmt.Println(recorder.Body.String())
	if recorder.Code != http.StatusInternalServerError {
		t.Error("Failed while testing the status code")
	}
}
