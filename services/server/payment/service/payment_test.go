package service

import (
	"testing"
	"services/server/test"
	"os"
)

var service *PaymentService

func TestMain(m *testing.M) {
	test.Setup("../../", false)

	// make sure to teardown on panic
	defer func() {
		if r := recover(); r != nil {
			test.Teardown(false)
		}
	}()

	service = NewPaymentService()

	code := m.Run()

	test.Teardown(false)

	os.Exit(code)
}

func TestSaveAccount(t *testing.T) {
	id, err := service.SaveAccount(&AccountSaveRequest{
		IPAddress: "1.2.3.4",
		FirstName: "Billy",
		LastName: "Fitzgerald",
		Email: "billy.fitzgerald@gmail.com",
		Country: "US",
		Currency: "USD",
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