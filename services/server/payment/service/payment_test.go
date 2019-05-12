package service

import (
	"testing"
	"services/server/test"
	"os"
	service2 "services/server/core/service"
)

var paymentService *PaymentService

func TestMain(m *testing.M) {
	test.Setup("../../", false)

	// make sure to teardown on panic
	defer func() {
		if r := recover(); r != nil {
			test.Teardown(false)
		}
	}()

	countryService := service2.NewCountryService("../../assets/data/country.json")
	paymentService = NewPaymentService(countryService)

	code := m.Run()

	test.Teardown(false)

	os.Exit(code)
}

// @todo write a test that basically creates an account, adds a bank account, fetches the account
// @todo and checks that the account doesn't have anything in the requirements

func TestSaveAccount(t *testing.T) {
	id, err := paymentService.SaveAccount(&AccountSaveRequest{
		IPAddress: "1.2.3.4",
		FirstName: "Billy",
		LastName: "Fitzgerald",
		Email: "billy.fitzgerald@gmail.com",
		Country: "US",
		SSN: "123451234",
		DOB: &PaymentDOB{
			Day: 1,
			Month: 1,
			Year: 1980,
		},
		Address: &PaymentAddress{
			LineOne: "1234 Line Lane",
			City: "Random",
			State: "FL",
			Country: "US",
			PostalCode: "12345",
		},
	})

	if id == nil || err != nil {
		t.Error("Failed to save account")
	}
}