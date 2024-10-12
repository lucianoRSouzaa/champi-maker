package application

import "github.com/google/uuid"

type ChampionshipCreatedMessage struct {
	ChampionshipID uuid.UUID   `json:"championship_id"`
	TeamIDs        []uuid.UUID `json:"team_ids"`
}
