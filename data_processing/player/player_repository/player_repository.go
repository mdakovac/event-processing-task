package player_repository

import (
	"context"

	"github.com/Bitstarz-eng/event-processing-challenge/db"
	"github.com/Bitstarz-eng/event-processing-challenge/internal/casino"
)

type PlayerRepository struct {
	conn *db.DbConnection
}

type PlayerRepositoryType interface {
	FindById(id int) (casino.Player, error)
}

func (repository *PlayerRepository) FindById(id int) (casino.Player, error) {
	var player casino.Player
	err := repository.conn.QueryRow(context.Background(), "SELECT email, last_signed_in_at FROM players WHERE id = $1", id).Scan(&player.Email, &player.LastSignedInAt)

	if err != nil {
		return player, err
	}

	return player, nil
}

func NewPlayerRepository(conn *db.DbConnection) PlayerRepositoryType {
	return &PlayerRepository{
		conn,
	}
}
