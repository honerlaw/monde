package service

import (
	"testing"
	"services/server/test"
	"os"
	"services/server/core/repository"
	"services/server/user/model"
	repository2 "services/server/payment/repository"
)

var user *model.User
var service *PaymentService

func TestMain(m *testing.M) {
	test.Setup("../../")

	// save the user so we can map everything to them
	user = &model.User{
		Hash: "test-hash",
	}
	repository.GetRepository().Save(user)

	// init the service to test
	service = NewPaymentService(repository2.NewStripeRepository(repository.GetRepository()))

	// make sure to teardown on panic
	defer func() {
		if r := recover(); r != nil {
			test.Teardown()
		}
	}()

	code := m.Run()

	test.Teardown()

	os.Exit(code)
}

func TestSaveAccount(t *testing.T) {
	service.SaveAccount(&AccountSaveRequest{

	})
}