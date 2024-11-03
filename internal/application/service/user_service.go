package service

import (
	"champi-maker/internal/application/port"
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/domain/repository"
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	RegisterUser(ctx context.Context, user *entity.User, password string) error
	AuthenticateUser(ctx context.Context, email, password string) (string, error) // Retorna o token JWT
	GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) error
	DeleteUser(ctx context.Context, id uuid.UUID) error
}

type userService struct {
	userRepo      repository.UserRepository
	tokenProvider port.TokenProvider
}

func NewUserService(userRepo repository.UserRepository, tokenProvider port.TokenProvider) UserService {
	return &userService{
		userRepo:      userRepo,
		tokenProvider: tokenProvider,
	}
}

func (s *userService) RegisterUser(ctx context.Context, user *entity.User, password string) error {
	if user.Name == "" || user.Email == "" || password == "" {
		return errors.New("nome, email e senha são obrigatórios")
	}

	existingUser, err := s.userRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return err
	}
	if existingUser != nil {
		return errors.New("email já está em uso")
	}

	hashedPassword, err := hashPassword(password)
	if err != nil {
		return err
	}
	user.PasswordHash = *hashedPassword

	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Create(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) AuthenticateUser(ctx context.Context, email, password string) (string, error) {
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("credenciais inválidas")
	}

	// Verificar a senha
	if err := verifyPassword(user.PasswordHash, password); err != nil {
		return "", errors.New("credenciais inválidas")
	}

	// Gerar o token JWT
	token, err := s.tokenProvider.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *userService) GetUserByID(ctx context.Context, id uuid.UUID) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("usuário não encontrado")
	}
	return user, nil
}

func (s *userService) UpdateUser(ctx context.Context, user *entity.User) error {
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}

	return nil
}

func (s *userService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.GetByID(ctx, id)

	if err != nil {
		return err
	}

	if user == nil {
		return errors.New("usuário não encontrado")
	}

	if err := s.userRepo.Delete(ctx, id); err != nil {
		return err
	}
	return nil
}

func hashPassword(password string) (*string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	hashed := string(bytes)
	return &hashed, nil
}

func verifyPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
