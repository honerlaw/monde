package service

import (
	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/client"
	"os"
)

type AccountService struct {
	client *client.API
}

func NewAccountService() (*AccountService) {
	return &AccountService{
		client: client.New(os.Getenv("STRIPE_SECRET_KEY"), nil),
	}
}

func (service *AccountService) Create(email string, country string) (*string, error) {
	account, err := service.client.Account.New(&stripe.AccountParams{
		Type: stripe.String(string(stripe.AccountTypeCustom)),
		Country: stripe.String("US"),
		Email: stripe.String(email),
		RequestedCapabilities: stripe.StringSlice([]string{
			"card_payments",
		}),
	})

	if err != nil {
		return nil, err
	}
	return &account.ID, nil
}
