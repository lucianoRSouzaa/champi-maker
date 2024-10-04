package entity

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestChampionshipValidation_Success(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Brasileiro",
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.NoError(t, err)
}

func TestChampionshipValidation_MissingID(t *testing.T) {
	championship := &Championship{
		// ID ausente
		Name:             "Campeonato Brasileiro",
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestChampionshipValidation_MissingName(t *testing.T) {
	championship := &Championship{
		ID: uuid.New(),
		// Name ausente
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestChampionshipValidation_MissingType(t *testing.T) {
	championship := &Championship{
		ID:   uuid.New(),
		Name: "Campeonato Brasileiro",
		// Type ausente
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Type", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestChampionshipValidation_MissingTiebreakerMethod(t *testing.T) {
	championship := &Championship{
		ID:   uuid.New(),
		Name: "Campeonato Brasileiro",
		Type: ChampionshipTypeLeague,
		// TiebreakerMethod ausente
		ProgressionType: ProgressionFixed,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "TiebreakerMethod", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestChampionshipValidation_MissingProgressionType(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Brasileiro",
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		// ProgressionType ausente
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ProgressionType", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestChampionshipValidation_MissingCreatedAt(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Brasileiro",
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		// CreatedAt ausente
		UpdatedAt: time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "CreatedAt", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestChampionshipValidation_MissingUpdatedAt(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Brasileiro",
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		// UpdatedAt ausente
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "UpdatedAt", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestChampionshipValidation_InvalidType(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Brasileiro",
		Type:             "invalid_type", // Valor inválido
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Type", validationErrors[0].Field())
	assert.Equal(t, "oneof", validationErrors[0].Tag())
}

func TestChampionshipValidation_InvalidTiebreakerMethod(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Brasileiro",
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: "invalid_method", // Valor inválido
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "TiebreakerMethod", validationErrors[0].Field())
	assert.Equal(t, "oneof", validationErrors[0].Tag())
}

func TestChampionshipValidation_InvalidProgressionType(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "Campeonato Brasileiro",
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  "invalid_progression", // Valor inválido
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ProgressionType", validationErrors[0].Field())
	assert.Equal(t, "oneof", validationErrors[0].Tag())
}

func TestChampionshipValidation_ShortName(t *testing.T) {
	championship := &Championship{
		ID:               uuid.New(),
		Name:             "A", // Muito curto, mínimo é 2 caracteres
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "min", validationErrors[0].Tag())
}

func TestChampionshipValidation_LongName(t *testing.T) {
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "a"
	}

	championship := &Championship{
		ID:               uuid.New(),
		Name:             longName, // Muito longo, máximo é 100 caracteres
		Type:             ChampionshipTypeLeague,
		TiebreakerMethod: TiebreakerPenalties,
		ProgressionType:  ProgressionFixed,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "max", validationErrors[0].Tag())
}

func TestChampionshipValidation_MultipleErrors(t *testing.T) {
	championship := &Championship{
		// ID ausente
		Name:             "A",         // Muito curto
		Type:             "invalid",   // Valor inválido
		TiebreakerMethod: "invalid",   // Valor inválido
		ProgressionType:  "invalid",   // Valor inválido
		CreatedAt:        time.Time{}, // Ausente (zero value)
		UpdatedAt:        time.Time{}, // Ausente (zero value)
	}

	err := championship.Validate()
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Len(t, validationErrors, 7)

	fieldsWithErrors := map[string]string{}
	for _, fieldErr := range validationErrors {
		fieldsWithErrors[fieldErr.Field()] = fieldErr.Tag()
	}

	expectedErrors := map[string]string{
		"ID":               "required",
		"Name":             "min",
		"Type":             "oneof",
		"TiebreakerMethod": "oneof",
		"ProgressionType":  "oneof",
		"CreatedAt":        "required",
		"UpdatedAt":        "required",
	}

	// Pode haver pequenas diferenças na ordem ou número de erros
	for field, tag := range expectedErrors {
		assert.Equal(t, tag, fieldsWithErrors[field], "Erro no campo %s", field)
	}
}
