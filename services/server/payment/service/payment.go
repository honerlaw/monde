package service

import (
	"github.com/stripe/stripe-go/client"
	"os"
	"github.com/stripe/stripe-go"
	"log"
)

type CardSaveRequest struct {
	CardID    *string
	AccountID string
}

type AccountSaveRequest struct {
	AccountID *string
	FirstName string
	LastName  string
	Email     string
	Country   string
	Currency  string
	SSNLast4  string
	DOB       struct {
		Day   int64
		Month int64
		Year  int64
	}
	Address struct {
		LineOne    string
		LineTwo    string
		City       string
		State      string
		Country    string
		PostalCode string
	}
}

type PaymentService struct {
	client           *client.API
}

func NewPaymentService() (*PaymentService) {
	return &PaymentService{
		client:           client.New(os.Getenv("STRIPE_SECRET_KEY"), nil),
	}
}

// all user's can have an account associated with them... We don't NEED certain information
// unless the user wants to allow donations / subscriptions
func (service *PaymentService) SaveAccount(req *AccountSaveRequest) (*string, error) {
	params := &stripe.AccountParams{
		Type:            stripe.String(string(stripe.AccountTypeCustom)),
		Country:         stripe.String(req.Country),
		Email:           stripe.String(req.Email),
		DefaultCurrency: stripe.String(req.Currency),
		RequestedCapabilities: stripe.StringSlice([]string{
			"platform_payments",
		}),
		BusinessType: stripe.String("individual"),
		Individual: &stripe.PersonParams{
			FirstName: stripe.String(req.FirstName),
			LastName:  stripe.String(req.LastName),
			Address: &stripe.AccountAddressParams{
				City:       stripe.String(req.Address.City),
				Country:    stripe.String(req.Address.Country),
				Line1:      stripe.String(req.Address.LineOne),
				Line2:      stripe.String(req.Address.LineTwo),
				PostalCode: stripe.String(req.Address.PostalCode),
				State:      stripe.String(req.Address.State),
			},
			DOB: &stripe.DOBParams{
				Day:   stripe.Int64(req.DOB.Day),
				Month: stripe.Int64(req.DOB.Month),
				Year:  stripe.Int64(req.DOB.Year),
			},
			SSNLast4: stripe.String(req.SSNLast4),
		},
	}

	var account *stripe.Account
	var err error
	if req.AccountID != nil {
		account, err = service.client.Account.Update(*req.AccountID, params)
	} else {
		account, err = service.client.Account.New(params)
	}

	if err != nil {
		return nil, err
	}
	return &account.ID, nil
}

func (service *PaymentService) GetAccount(accountID string) (*stripe.Account) {
	account, err := service.client.Account.GetByID(accountID, nil)
	if err != nil {
		log.Print(err)
		return nil
	}
	return account
}

func (service *PaymentService) SaveCard(req *CardSaveRequest) (error) {
	params := &stripe.CardParams{
		Account: stripe.String(req.AccountID),
	}

	var err error
	if req.CardID != nil {
		_, err = service.client.Cards.Update(*req.CardID, params)
	} else {
		_, err = service.client.Cards.New(params)
	}

	if err != nil {
		return err
	}
	return nil
}

func (service *PaymentService) ListCards(accountID string) ([]stripe.Card) {
	iter := service.client.Cards.List(&stripe.CardListParams{
		Account: stripe.String(accountID),
	})

	if iter.Err() != nil {
		log.Print(iter.Err())
		return nil
	}

	cards := make([]stripe.Card, 0)
	for iter.Next() {
		cards = append(cards, *iter.Card())
	}

	return cards
}
