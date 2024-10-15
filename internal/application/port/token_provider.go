package port

import "github.com/google/uuid"

type TokenProvider interface {
	GenerateToken(userID uuid.UUID) (string, error)
	// ValidateToken(token string) (uuid.UUID, error)
}
