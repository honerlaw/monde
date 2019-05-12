package service

import (
	"io/ioutil"
	"github.com/labstack/gommon/log"
	"encoding/json"
)

type CountryDataISO struct {
	Two    string `json:"two"`
	Three  string `json:"three"`
	Number string `json:"number"`
}

type CountryData struct {
	Name string         `json:"name" `
	ISO  CountryDataISO `json:"iso"`
	CurrencyCodes []string `json:"currency_codes"`
}

type CountryService struct {
	data []CountryData
}

func NewCountryService(dataPath string) (*CountryService) {
	bytes, err := ioutil.ReadFile(dataPath)
	if err != nil {
		log.Fatal(err)
	}

	data := make([]CountryData, 0)
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		log.Fatal(err)
	}

	return &CountryService{
		data: data,
	}
}

func (service *CountryService) GetByISOCode(code string) (*CountryData) {
	for _, val := range service.data {
		if val.ISO.Two == code || val.ISO.Three == code || val.ISO.Number == code {
			return &val
		}
	}
	return nil
}
