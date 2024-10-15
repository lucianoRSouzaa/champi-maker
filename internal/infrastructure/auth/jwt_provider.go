package security

import (
	"champi-maker/internal/application/port"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type JWTService struct {
	secretKey []byte
	issuer    string
	expiry    time.Duration
}

func NewJWTService(secretKey string, issuer string, expiry time.Duration) port.TokenProvider {
	return &JWTService{
		secretKey: []byte(secretKey),
		issuer:    issuer,
		expiry:    expiry,
	}
}

func (j *JWTService) GenerateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"iss": j.issuer,
		"sub": userID.String(),
		"exp": time.Now().Add(j.expiry).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}
