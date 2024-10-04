package entity

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestTeamValidation_Success(t *testing.T) {
	team := &Team{
		ID:        uuid.New(),
		Name:      "Team A",
		Logo:      "https://example.com/logo.png",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.NoError(t, err)
}

func TestTeamValidation_MissingID(t *testing.T) {
	team := &Team{
		// ID ausente
		Name:      "Team A",
		Logo:      "https://example.com/logo.png",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestTeamValidation_MissingName(t *testing.T) {
	team := &Team{
		ID: uuid.New(),
		// Name ausente
		Logo:      "https://example.com/logo.png",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestTeamValidation_MissingUserID(t *testing.T) {
	team := &Team{
		ID:   uuid.New(),
		Name: "Team A",
		Logo: "https://example.com/logo.png",
		// UserID ausente
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "UserID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestTeamValidation_MissingCreatedAt(t *testing.T) {
	team := &Team{
		ID:     uuid.New(),
		Name:   "Team A",
		Logo:   "https://example.com/logo.png",
		UserID: uuid.New(),
		// CreatedAt ausente
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "CreatedAt", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestTeamValidation_MissingUpdatedAt(t *testing.T) {
	team := &Team{
		ID:        uuid.New(),
		Name:      "Team A",
		Logo:      "https://example.com/logo.png",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		// UpdatedAt ausente
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "UpdatedAt", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestTeamValidation_InvalidLogoURL(t *testing.T) {
	team := &Team{
		ID:        uuid.New(),
		Name:      "Team A",
		Logo:      "not-a-valid-url",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Logo", validationErrors[0].Field())
	assert.Equal(t, "url", validationErrors[0].Tag())
}

func TestTeamValidation_ShortName(t *testing.T) {
	team := &Team{
		ID:        uuid.New(),
		Name:      "A", // Muito curto, mínimo é 2 caracteres
		Logo:      "https://example.com/logo.png",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "min", validationErrors[0].Tag())
}

func TestTeamValidation_LongName(t *testing.T) {
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "a"
	}

	team := &Team{
		ID:        uuid.New(),
		Name:      longName, // Muito longo, máximo é 100 caracteres
		Logo:      "https://example.com/logo.png",
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "max", validationErrors[0].Tag())
}

func TestTeamValidation_MultipleErrors(t *testing.T) {
	var invalidUUID uuid.UUID // UUID zero value

	team := &Team{
		// ID ausente
		Name:      "A",               // Muito curto
		Logo:      "not-a-valid-url", // URL inválida
		UserID:    invalidUUID,       // UUID inválido
		CreatedAt: time.Time{},       // Ausente (zero value)
		UpdatedAt: time.Time{},       // Ausente (zero value)
	}

	err := team.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Len(t, validationErrors, 6)

	fieldsWithErrors := map[string]string{}
	for _, fieldErr := range validationErrors {
		fieldsWithErrors[fieldErr.Field()] = fieldErr.Tag()
	}

	expectedErrors := map[string]string{
		"ID":        "required",
		"Name":      "min",
		"Logo":      "url",
		"UserID":    "required",
		"CreatedAt": "required",
		"UpdatedAt": "required",
	}

	assert.Equal(t, expectedErrors, fieldsWithErrors)
}

func TestTeamValidation_LogoOptional(t *testing.T) {
	team := &Team{
		ID:   uuid.New(),
		Name: "Team A",
		// Logo ausente (campo opcional)
		UserID:    uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := team.Validate()
	assert.NoError(t, err)
}
