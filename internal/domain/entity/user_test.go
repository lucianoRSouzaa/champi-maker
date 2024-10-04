package entity

import (
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func TestUserValidation_Success(t *testing.T) {

	u := &User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := validate.Struct(u)
	assert.NoError(t, err)
}

func TestUserValidation_MissingID(t *testing.T) {
	user := &User{
		// ID ausente
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "ID", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestUserValidation_MissingName(t *testing.T) {
	user := &User{
		ID: uuid.New(),
		// Name ausente
		Email:        "john.doe@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestUserValidation_MissingEmail(t *testing.T) {
	user := &User{
		ID:   uuid.New(),
		Name: "John Doe",
		// Email ausente
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Email", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestUserValidation_MissingPasswordHash(t *testing.T) {
	user := &User{
		ID:    uuid.New(),
		Name:  "John Doe",
		Email: "john.doe@example.com",
		// PasswordHash ausente
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "PasswordHash", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestUserValidation_MissingCreatedAt(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "hashedpassword123",
		// CreatedAt ausente
		UpdatedAt: time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "CreatedAt", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestUserValidation_MissingUpdatedAt(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "john.doe@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		// UpdatedAt ausente
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "UpdatedAt", validationErrors[0].Field())
	assert.Equal(t, "required", validationErrors[0].Tag())
}

func TestUserValidation_InvalidEmail(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Name:         "John Doe",
		Email:        "invalid-email",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Email", validationErrors[0].Field())
	assert.Equal(t, "email", validationErrors[0].Tag())
}

func TestUserValidation_ShortName(t *testing.T) {
	user := &User{
		ID:           uuid.New(),
		Name:         "J", // Muito curto, mínimo é 2 caracteres
		Email:        "john.doe@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "min", validationErrors[0].Tag())
}

func TestUserValidation_LongName(t *testing.T) {
	longName := ""
	for i := 0; i < 101; i++ {
		longName += "a"
	}

	user := &User{
		ID:           uuid.New(),
		Name:         longName, // Muito longo, máximo é 100 caracteres
		Email:        "john.doe@example.com",
		PasswordHash: "hashedpassword123",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Equal(t, "Name", validationErrors[0].Field())
	assert.Equal(t, "max", validationErrors[0].Tag())
}

func TestUserValidation_MultipleErrors(t *testing.T) {
	user := &User{
		// ID ausente
		Name:         "J",             // Muito curto
		Email:        "invalid-email", // Formato inválido
		PasswordHash: "",              // Ausente
		CreatedAt:    time.Time{},     // Ausente (zero value)
		UpdatedAt:    time.Time{},     // Ausente (zero value)
	}

	err := validate.Struct(user)
	assert.Error(t, err)
	validationErrors := err.(validator.ValidationErrors)
	assert.Len(t, validationErrors, 6)

	fieldsWithErrors := map[string]string{}
	for _, fieldErr := range validationErrors {
		fieldsWithErrors[fieldErr.Field()] = fieldErr.Tag()
	}

	expectedErrors := map[string]string{
		"ID":           "required",
		"Name":         "min",
		"Email":        "email",
		"PasswordHash": "required",
		"CreatedAt":    "required",
		"UpdatedAt":    "required",
	}

	assert.Equal(t, expectedErrors, fieldsWithErrors)
}
