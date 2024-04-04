package currency_service

import (
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/currency/currency_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
	"github.com/Bitstarz-eng/event-processing-challenge/util/math_utils"
)

type currencyService struct {
	repository currency_repository.CurrencyRepositoryType
}

type CurrencyServiceType interface {
	ConvertCurrency(event *casino.Event) (*casino.Event, error)
	ConvertAmountToFloat(amount int, currency string) float64
}

const conversionConstant float64 = 100
const btcConversionConstant float64 = 100000000

func (service *currencyService) ConvertCurrency(event *casino.Event) (*casino.Event, error) {
	if event.Type != casino.EventTypeBet && event.Type != casino.EventTypeDeposit {
		return event, nil
	}

	if event.Currency == casino.CurrencyEUR {
		event.AmountEUR = event.Amount
		return event, nil
	}

	exchangeRates, err := service.repository.GetExchangeRates()
	if err != nil {
		log.Println("Error getting exchange rates", err)
		return event, err
	}

	amountFloat := service.ConvertAmountToFloat(event.Amount, event.Currency)

	amountEurFloat := math_utils.ToFixed(amountFloat*exchangeRates[event.Currency], 2)
	event.AmountEUR = int(amountEurFloat * conversionConstant)

	return event, nil
}

func (*currencyService) ConvertAmountToFloat(amount int, currency string) float64 {
	cc := conversionConstant
	if currency == casino.CurrencyBTC {
		cc = btcConversionConstant
	}

	amountFloat := float64(amount) / cc
	return amountFloat
}

func NewCurrencyService(repository currency_repository.CurrencyRepositoryType) CurrencyServiceType {
	return &currencyService{
		repository,
	}
}
