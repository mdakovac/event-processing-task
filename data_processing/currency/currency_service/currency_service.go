package currency_service

import (
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/currency/currency_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/util/math_utils"
)

type CurrencyService struct {
	repository currency_repository.CurrencyRepositoryType
}

type CurrencyServiceType interface {
	ConvertCurrency(event *casino.Event) (*casino.Event, error)
	ConvertAmountToFloat(amount int, currency string) float64
}

const conversionConstant float64 = 100
const btcConversionConstant float64 = 100000000

func (service *CurrencyService) ConvertCurrency(event *casino.Event) (*casino.Event, error) {
	if event.Type != "bet" && event.Type != "deposit" {
		return event, nil
	}

	if event.Currency == "EUR" {
		event.AmountEUR = event.Amount
		return event, nil
	} else {
		exchangeRates, err := service.repository.GetExchangeRates()
		if err != nil {
			log.Println("Error getting exchange rates", err)
			return event, err
		}

		amountFloat := service.ConvertAmountToFloat(event.Amount, event.Currency)

		amountEurFloat := math_utils.ToFixed(amountFloat*exchangeRates[event.Currency], 2)
		event.AmountEUR = int(amountEurFloat * conversionConstant)
	}

	return event, nil
}

func (*CurrencyService) ConvertAmountToFloat(amount int, currency string) float64 {
	var amountFloat float64
	if currency == "BTC" {
		amountFloat = float64(amount) / btcConversionConstant
	} else {
		amountFloat = float64(amount) / conversionConstant
	}

	return amountFloat
}

func NewCurrencyService(repository currency_repository.CurrencyRepositoryType) CurrencyServiceType {
	return &CurrencyService{
		repository,
	}
}
