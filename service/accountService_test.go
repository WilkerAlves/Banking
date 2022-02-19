package service

import (
	"github.com/WilkerAlves/banking-lib/errs"
	realdomain "github.com/WilkerAlves/banking/domain"
	"github.com/WilkerAlves/banking/dto"

	"testing"

	"github.com/WilkerAlves/banking/mocks/domain"
	"github.com/golang/mock/gomock"
)

var mockRepo *domain.MockAccountRepository
var service AccountService

func setup(t *testing.T) func() {
	ctrl := gomock.NewController(t)
	mockRepo = domain.NewMockAccountRepository(ctrl)
	service = NewAccountService(mockRepo)
	return func() {
		service = nil
		defer ctrl.Finish()
	}
}

func Test_should_return_a_validation_error_response_when_the_request_is_not_validated(t *testing.T) {
	// Arrange
	request := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      0,
	}
	service := NewAccountService(nil)
	// Act
	_, appError := service.NewAccount(request)
	// Assert
	if appError == nil {
		t.Error("failed while testing the new account validation")
	}
}

func Test_should_return_an_error_from_the_server_side_if_the_new_account_cannot_be_created(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      6000,
	}

	account := realdomain.NewAccount(req.CustomerId, req.AccountType, req.Amount)

	mockRepo.EXPECT().Save(account).Return(nil, errs.NewUnexpectedError("Unexpected database error"))
	// Act
	_, appError := service.NewAccount(req)

	// Assert
	if appError == nil {
		t.Error("Test failed while validating error for new account")
	}

}

func Test_should_return_new_account_response_when_a_new_account_is_saved_successfully(t *testing.T) {
	// Arrange
	teardown := setup(t)
	defer teardown()

	req := dto.NewAccountRequest{
		CustomerId:  "100",
		AccountType: "saving",
		Amount:      6000,
	}
	account := realdomain.NewAccount(req.CustomerId, req.AccountType, req.Amount)
	accountWithId := account
	accountWithId.AccountId = "201"
	mockRepo.EXPECT().Save(account).Return(&accountWithId, nil)
	// Act
	newAccount, appError := service.NewAccount(req)

	// Assert
	if appError != nil {
		t.Error("Test failed while creating new account")
	}
	if newAccount.AccountId != accountWithId.AccountId {
		t.Error("Failed while matching new account id")
	}
}
