package service

import (
	"github.com/stripe/stripe-go/client"
	"os"
	"github.com/stripe/stripe-go"
	"log"
	"time"
	"services/server/core/service"
	"github.com/pkg/errors"
)

type BankAccountSaveRequest struct {
	BankAccountID *string
	AccountID string
	Country string
	AccountHolderName string
	AccountNumber string
	AccountRoutingNumber string
}

type CardSaveRequest struct {
	CardID    *string
	AccountID string
}

type PaymentDOB struct {
	Day   int64
	Month int64
	Year  int64
}

type PaymentAddress struct {
	LineOne    string
	LineTwo    string
	City       string
	State      string
	Country    string
	PostalCode string
}
type AccountSaveRequest struct {
	AccountID *string
	IPAddress string
	FirstName string
	LastName  string
	Email     string
	Country   string
	SSN       string
	DOB       *PaymentDOB
	Address   *PaymentAddress
}

type PaymentService struct {
	client *client.API
	countryService *service.CountryService
}

func NewPaymentService(countryService *service.CountryService) (*PaymentService) {
	return &PaymentService{
		client: client.New(os.Getenv("STRIPE_SECRET_KEY"), nil),
		countryService: countryService,
	}
}

// all user's can have an account associated with them... We don't NEED certain information
// unless the user wants to allow donations / subscriptions
func (service *PaymentService) SaveAccount(req *AccountSaveRequest) (*string, error) {
	countryData := service.countryService.GetByISOCode(req.Country)
	if countryData == nil {
		return nil, errors.New("invalid country code")
	}

	params := &stripe.AccountParams{
		Type:            stripe.String(string(stripe.AccountTypeCustom)),
		Country:         stripe.String(req.Country),
		Email:           stripe.String(req.Email),
		DefaultCurrency: stripe.String(countryData.CurrencyCodes[0]),
		RequestedCapabilities: stripe.StringSlice([]string{
			"platform_payments",
		}),
		BusinessType: stripe.String("individual"),
		BusinessProfile: &stripe.AccountBusinessProfileParams{
			ProductDescription: stripe.String("The user agreed to terms in order to sell content on our platform."),
		},
		TOSAcceptance: &stripe.AccountTOSAcceptanceParams{
			Date: stripe.Int64(time.Now().Unix()),
			IP:   stripe.String(req.IPAddress),
		},
		Individual: &stripe.PersonParams{
			IDNumber:  stripe.String(req.SSN),
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
			SSNLast4: stripe.String(string(req.SSN[len(req.SSN)-4:])),
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

func (service *PaymentService) SaveBankAccount(req *BankAccountSaveRequest) (error) {
	countryData := service.countryService.GetByISOCode(req.Country)
	if countryData == nil {
		return errors.New("invalid country code")
	}

	params := &stripe.BankAccountParams{
		Account: stripe.String(req.AccountID),
		AccountHolderName: stripe.String(req.AccountHolderName),
		AccountHolderType: stripe.String(string(stripe.BankAccountAccountHolderTypeIndividual)),
		AccountNumber: stripe.String(req.AccountNumber),
		RoutingNumber: stripe.String(req.AccountRoutingNumber),
		Currency: stripe.String(countryData.CurrencyCodes[0]),
		Country: stripe.String(req.Country),
	}

	var err error
	if req.BankAccountID != nil {
		_, err = service.client.BankAccounts.Update(*req.BankAccountID, params)
	} else {
		_, err = service.client.BankAccounts.New(params)
	}

	if err != nil {
		log.Print(err)
		return err
	}

	return nil
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
