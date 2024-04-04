package player_service

import (
	"log"

	"github.com/Bitstarz-eng/event-processing-challenge/data_processing/player/player_repository"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type playerService struct {
	repository player_repository.PlayerRepositoryType
}

type PlayerServiceType interface {
	AssignPlayerData(event *casino.Event) (*casino.Event, error)
}

func (service *playerService) AssignPlayerData(event *casino.Event) (*casino.Event, error) {
	player, err := service.repository.FindById(event.PlayerID)
	if err != nil {
		log.Printf("Unable to find Player Data for event id %v", event.ID)
		return event, err
	}

	event.Player = player
	return event, nil
}

func NewPlayerService(repository player_repository.PlayerRepositoryType) PlayerServiceType {
	return &playerService{
		repository,
	}
}
