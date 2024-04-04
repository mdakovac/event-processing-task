package aggregation_service

import (
	"sync"
	"time"

	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/aggregation/aggregation_models"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type aggregationService struct {
	aggregation aggregation_models.Aggregation

	winsByPlayer     aggregation_models.CountByPlayerId
	betsByPlayer     aggregation_models.CountByPlayerId
	depositsByPlayer aggregation_models.CountByPlayerId
	eventsBySecond   aggregation_models.CountByTimestamp

	earliestEventTimestamp time.Time

	mutex sync.RWMutex
}

type AggregationServiceType interface {
	AddEventToAggregation(event *casino.Event)
	GetAggregation() *aggregation_models.Aggregation
}

func (service *aggregationService) AddEventToAggregation(event *casino.Event) {
	service.mutex.Lock()
	defer service.mutex.Unlock()

	service.aggregation.EventsTotal++
	incrementMapValue(service.eventsBySecond, event.CreatedAt.Truncate(time.Second).String())

	if service.earliestEventTimestamp.IsZero() || event.CreatedAt.Before(service.earliestEventTimestamp) {
		service.earliestEventTimestamp = event.CreatedAt
	}

	if event.Type == casino.EventTypeBet {
		incrementMapValue(service.betsByPlayer, event.PlayerID)
		key, value := findMaxValue(service.betsByPlayer)
		updateAggregationTopPlayer(&service.aggregation.TopPlayerBets, key, value)

		if event.HasWon {
			incrementMapValue(service.winsByPlayer, event.PlayerID)
			key, value := findMaxValue(service.winsByPlayer)
			updateAggregationTopPlayer(&service.aggregation.TopPlayerWins, key, value)
		}

	} else if event.Type == casino.EventTypeDeposit {
		incrementMapValue(service.depositsByPlayer, event.PlayerID)
		key, value := findMaxValue(service.depositsByPlayer)
		updateAggregationTopPlayer(&service.aggregation.TopPlayerDeposits, key, value)
	}
}

func (service *aggregationService) GetAggregation() *aggregation_models.Aggregation {
	// stil need RW Lock since getting the Aggregation also updates fields
	service.mutex.Lock()
	defer service.mutex.Unlock()

	updateAggregationEventsPerMinute(&service.earliestEventTimestamp, &service.aggregation)
	// TODO:  events_per_second_moving_average

	return &service.aggregation
}

func updateAggregationEventsPerMinute(earliestEventTimestamp *time.Time, aggregation *aggregation_models.Aggregation) {
	if earliestEventTimestamp.IsZero() {
		return
	}

	minutesElapsed := time.Since(*earliestEventTimestamp).Minutes()
	aggregation.EventsPerMinute = float64(aggregation.EventsTotal) / minutesElapsed
}

type IntOrString interface {
	string | int
}

func incrementMapValue[V IntOrString](m map[V]int, key V) {
	_, exists := m[key]
	if exists {
		m[key]++
	} else {
		m[key] = 0
		m[key]++
	}
}

func findMaxValue(m map[int]int) (int, int) {
	var maxKey int
	var maxValue int

	for key, value := range m {
		if value > maxValue {
			maxValue = value
			maxKey = key
		}
	}

	return maxKey, maxValue
}

func updateAggregationTopPlayer(s *aggregation_models.TopPlayer, id int, count int) {
	s.ID = id
	s.Count = count
}

/* func prettyPrint(v any) {
	forPrint1, _ := json.MarshalIndent(&v, "", "    ")
	log.Println("Aggregation event", string(forPrint1))
} */

func NewAggregationService() AggregationServiceType {
	return &aggregationService{
		winsByPlayer:     make(aggregation_models.CountByPlayerId),
		betsByPlayer:     make(aggregation_models.CountByPlayerId),
		depositsByPlayer: make(aggregation_models.CountByPlayerId),
		eventsBySecond:   make(aggregation_models.CountByTimestamp),
	}
}
