package currency_service

import (
	"math"

	"github.com/Bitstarz-eng/event-processing-challenge/data_enrichment/currency/currency_models"
	"github.com/Bitstarz-eng/event-processing-challenge/data_enrichment/currency/currency_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type CurrencyService struct {
	repository currency_repository.CurrencyRepositoryType
}

type CurrencyServiceType interface {
	GetExchangeRates() (currency_models.ExchangeRates, error)
}

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
			return nil, err
		}

		const conversionConstant float64 = 100
		const btcConversionConstant float64 = 100000000

		var amountFloat float64
		if event.Currency == "BTC" {
			amountFloat = float64(event.Amount) / btcConversionConstant
		} else {
			amountFloat = float64(event.Amount) / conversionConstant
		}

		amountEurFloat := toFixed(amountFloat*exchangeRates[event.Currency], 2)
		event.AmountEUR = int(amountEurFloat * conversionConstant)
	}

	return event, nil
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func NewCurrencyService(repository currency_repository.CurrencyRepositoryType) *CurrencyService {
	return &CurrencyService{
		repository,
	}
}
