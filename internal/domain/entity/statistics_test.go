package entity

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStatisticsValidation_Success(t *testing.T) {
	stats := &Statistics{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		TeamID:         uuid.New(),
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := stats.Validate()
	assert.NoError(t, err)
}

func TestStatisticsValidation_MissingID(t *testing.T) {
	stats := &Statistics{
		// ID ausente
		ChampionshipID: uuid.New(),
		TeamID:         uuid.New(),
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := stats.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestStatisticsValidation_MissingChampionshipID(t *testing.T) {
	stats := &Statistics{
		ID: uuid.New(),
		// ChampionshipID ausente
		TeamID:         uuid.New(),
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := stats.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ChampionshipID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestStatisticsValidation_MissingTeamID(t *testing.T) {
	stats := &Statistics{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		// TeamID ausente
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := stats.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "TeamID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestStatisticsValidation_NegativeMatchesPlayed(t *testing.T) {
	stats := &Statistics{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		TeamID:         uuid.New(),
		MatchesPlayed:  -1, // Valor inválido
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := stats.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "MatchesPlayed", validationErrors[0].Field())
	assert.Equal(t, "gte", validationErrors[0].Tag())
}

func TestStatisticsValidation_NegativeWins(t *testing.T) {
	stats := &Statistics{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		TeamID:         uuid.New(),
		MatchesPlayed:  10,
		Wins:           -1, // Valor inválido
		Draws:          2,
		Losses:         2,
		GoalsFor:       18,
		GoalsAgainst:   10,
		GoalDifference: 8,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := stats.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Wins", validationErrors[0].Field())
	assert.Equal(t, "gte", validationErrors[0].Tag())
}

func TestStatisticsValidation_NegativeGoalsFor(t *testing.T) {
	stats := &Statistics{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		TeamID:         uuid.New(),
		MatchesPlayed:  10,
		Wins:           6,
		Draws:          2,
		Losses:         2,
		GoalsFor:       -5, // Valor inválido
		GoalsAgainst:   10,
		GoalDifference: -15,
		Points:         20,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	err := stats.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "GoalsFor", validationErrors[0].Field())
	assert.Equal(t, "gte", validationErrors[0].Tag())
}

func TestStatisticsValidation_MultipleErrors(t *testing.T) {
	stats := &Statistics{
		// ID ausente
		ChampionshipID: uuid.Nil,    // UUID inválido
		TeamID:         uuid.Nil,    // UUID inválido
		MatchesPlayed:  -1,          // Valor inválido
		Wins:           -1,          // Valor inválido
		Draws:          -1,          // Valor inválido
		Losses:         -1,          // Valor inválido
		GoalsFor:       -1,          // Valor inválido
		GoalsAgainst:   -1,          // Valor inválido
		GoalDifference: -1,          // Valor inválido
		Points:         -1,          // Valor inválido
		CreatedAt:      time.Time{}, // Ausente
		UpdatedAt:      time.Time{}, // Ausente
	}

	err := stats.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)

	// Esperamos múltiplos erros
	assert.GreaterOrEqual(t, len(validationErrors), 10)
}
