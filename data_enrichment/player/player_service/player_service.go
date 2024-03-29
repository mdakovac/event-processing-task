package player_service

import (
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/data_enrichment/player/player_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type PlayerService struct {
	repository player_repository.PlayerRepositoryType
}

type PlayerServiceType interface {
	AssignPlayerData(event *casino.Event) (*casino.Event, error)
}

func (service *PlayerService) AssignPlayerData(event *casino.Event) (*casino.Event, error) {
	player, err := service.repository.FindById(event.PlayerID)
	if err != nil {
		log.Println("Unable to find Player Data for event", event)
		return event, err
	}

	event.Player = player
	return event, nil
}

func NewPlayerService(repository player_repository.PlayerRepositoryType) *PlayerService {
	return &PlayerService{
		repository,
	}
}
