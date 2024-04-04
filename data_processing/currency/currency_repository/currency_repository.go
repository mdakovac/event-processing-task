package currency_repository

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/cache_service"
	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/currency/currency_models"
	"github.com/Bitstarz-eng/event-processing-challenge/util/env_vars"
)

type currencyRepository struct {
	cache *cache_service.Cache
	mutex sync.Mutex
}

type CurrencyRepositoryType interface {
	GetExchangeRates() (currency_models.ExchangeRates, error)
}

func (repository *currencyRepository) GetExchangeRates() (currency_models.ExchangeRates, error) {
	var exchangeRates currency_models.ExchangeRates

	repository.mutex.Lock()
	defer repository.mutex.Unlock()
	
	exchangeRates = findFromCache(repository.cache)
	if exchangeRates != nil {
		return exchangeRates, nil
	}

	exchangeRates, err := findFromApi()
	if err != nil {
		return nil, err
	}

	repository.cache.Set("exchange_rates", exchangeRates, time.Minute)

	return exchangeRates, nil
}

func findFromApi() (currency_models.ExchangeRates, error) {
	response, err := http.Get(env_vars.EnvVariables.EXCHANGE_RATES_API_URL + "/exchange-rates")
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, errors.New("Request failed with status:" + response.Status)
	}

	var exchangeRates currency_models.ExchangeRates
	err = json.NewDecoder(response.Body).Decode(&exchangeRates)
	if err != nil {
		return nil, err
	}

	return exchangeRates, nil
}

func findFromCache(cache *cache_service.Cache) currency_models.ExchangeRates {
	cached, found := cache.Get("exchange_rates")

	if found {
		a := cached.(currency_models.ExchangeRates)
		return a
	}

	return nil
}

func NewCurrencyRepository(cache *cache_service.Cache) CurrencyRepositoryType {
	return &currencyRepository{
		cache: cache,
	}
}
