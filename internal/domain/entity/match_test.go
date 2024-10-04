package entity

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMatchValidation_Success(t *testing.T) {
	now := time.Now()
	matchDate := now.Add(24 * time.Hour)

	homeTeamID := uuid.New()
	awayTeamID := uuid.New()

	match := &Match{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		MatchDate:      &matchDate,
		Status:         MatchStatusScheduled,
		ScoreHome:      0,
		ScoreAway:      0,
		HasExtraTime:   false,
		HasPenalties:   false,
		Phase:          1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := match.Validate()
	assert.NoError(t, err)
}

func TestMatchValidation_MissingID(t *testing.T) {
	now := time.Now()

	match := &Match{
		// ID ausente
		ChampionshipID: uuid.New(),
		Phase:          1,
		Status:         MatchStatusScheduled,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestMatchValidation_MissingChampionshipID(t *testing.T) {
	now := time.Now()

	match := &Match{
		ID: uuid.New(),
		// ChampionshipID ausente
		Phase:     1,
		Status:    MatchStatusScheduled,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ChampionshipID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestMatchValidation_MissingPhase(t *testing.T) {
	now := time.Now()

	homeTeamID := uuid.New()
	awayTeamID := uuid.New()

	match := &Match{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		Status:         MatchStatusScheduled,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Phase", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestMatchValidation_MissingStatus(t *testing.T) {
	now := time.Now()

	homeTeamID := uuid.New()
	awayTeamID := uuid.New()

	match := &Match{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		Phase:          1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Status", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestMatchValidation_InvalidStatus(t *testing.T) {
	now := time.Now()

	homeTeamID := uuid.New()
	awayTeamID := uuid.New()

	match := &Match{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		Phase:          1,
		Status:         "invalid_status", // Inválido
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Status", validationErrors[0].Field())
	assert.Equal(t, "oneof", validationErrors[0].Tag())
}

func TestMatchValidation_NegativeScoreHome(t *testing.T) {
	now := time.Now()

	homeTeamID := uuid.New()
	awayTeamID := uuid.New()

	match := &Match{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		Phase:          1,
		Status:         MatchStatusScheduled,
		ScoreHome:      -1, // Inválido
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ScoreHome", validationErrors[0].Field())
	assert.Equal(t, "gte", validationErrors[0].Tag())
}

func TestMatchValidation_NegativeScoreAway(t *testing.T) {
	now := time.Now()

	homeTeamID := uuid.New()
	awayTeamID := uuid.New()

	match := &Match{
		ID:             uuid.New(),
		ChampionshipID: uuid.New(),
		HomeTeamID:     &homeTeamID,
		AwayTeamID:     &awayTeamID,
		Phase:          1,
		Status:         MatchStatusScheduled,
		ScoreAway:      -1, // Inválido
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ScoreAway", validationErrors[0].Field())
	assert.Equal(t, "gte", validationErrors[0].Tag())
}

func TestMatchValidation_MultipleErrors(t *testing.T) {
	match := &Match{
		// ID ausente
		ChampionshipID: uuid.Nil,    // UUID inválido (zero value)
		Phase:          0,           // Inválido
		Status:         "invalid",   // Inválido
		ScoreHome:      -1,          // Inválido
		ScoreAway:      -1,          // Inválido
		CreatedAt:      time.Time{}, // Ausente
		UpdatedAt:      time.Time{}, // Ausente
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Len(t, validationErrors, 8)

	fieldsWithErrors := map[string]string{}
	for _, fieldErr := range validationErrors {
		fieldsWithErrors[fieldErr.Field()] = fieldErr.Tag()
	}

	expectedErrors := map[string]string{
		"ID":             "required",
		"ChampionshipID": "required",
		"Phase":          "required",
		"Status":         "oneof",
		"ScoreHome":      "gte",
		"ScoreAway":      "gte",
		"CreatedAt":      "required",
		"UpdatedAt":      "required",
	}

	assert.Equal(t, expectedErrors, fieldsWithErrors)
}

func TestMatchValidation_NegativeScoreHomeExtraTime(t *testing.T) {
	now := time.Now()

	match := &Match{
		ID:                 uuid.New(),
		ChampionshipID:     uuid.New(),
		Phase:              1,
		Status:             MatchStatusScheduled,
		HasExtraTime:       true,
		ScoreHomeExtraTime: -1, // Inválido
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ScoreHomeExtraTime", validationErrors[0].Field())
	assert.Equal(t, "gte", validationErrors[0].Tag())
}

func TestMatchValidation_NegativeScoreAwayPenalties(t *testing.T) {
	now := time.Now()

	match := &Match{
		ID:                 uuid.New(),
		ChampionshipID:     uuid.New(),
		Phase:              1,
		Status:             MatchStatusScheduled,
		HasPenalties:       true,
		ScoreAwayPenalties: -1, // Inválido
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	err := match.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ScoreAwayPenalties", validationErrors[0].Field())
	assert.Equal(t, "gte", validationErrors[0].Tag())
}

func TestMatchValidation_WithParentAndChildMatches(t *testing.T) {
	now := time.Now()
	parentMatchID := uuid.New()
	leftChildMatchID := uuid.New()
	rightChildMatchID := uuid.New()

	match := &Match{
		ID:                uuid.New(),
		ChampionshipID:    uuid.New(),
		Phase:             2,
		Status:            MatchStatusScheduled,
		ParentMatchID:     &parentMatchID,
		LeftChildMatchID:  &leftChildMatchID,
		RightChildMatchID: &rightChildMatchID,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	err := match.Validate()
	assert.NoError(t, err)
}
