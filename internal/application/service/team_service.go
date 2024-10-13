package service

import (
	"champi-maker/internal/domain/entity"
	"champi-maker/internal/domain/repository"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type TeamService interface {
	CreateTeam(ctx context.Context, team *entity.Team) error
	GetTeamByID(ctx context.Context, teamID uuid.UUID) (*entity.Team, error)
	UpdateTeam(ctx context.Context, team *entity.Team) error
	DeleteTeam(ctx context.Context, teamID uuid.UUID) error
	ListTeams(ctx context.Context) ([]*entity.Team, error)
	ListTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Team, error)
}

type teamService struct {
	teamRepo repository.TeamRepository
	userRepo repository.UserRepository
}

func NewTeamService(teamRepo repository.TeamRepository, userRepo repository.UserRepository) TeamService {
	return &teamService{
		teamRepo: teamRepo,
		userRepo: userRepo,
	}
}

func (s *teamService) CreateTeam(ctx context.Context, team *entity.Team) error {
	// Validate the team entity
	if err := team.Validate(); err != nil {
		return err
	}

	// Check if the user exists
	user, err := s.userRepo.GetByID(ctx, team.UserID)
	if err != nil {
		return err
	}
	if user == nil {
		return fmt.Errorf("user with ID %s not found", team.UserID)
	}

	// Check if a team with the same name already exists
	existingTeams, err := s.teamRepo.List(ctx)
	if err != nil {
		return err
	}
	for _, existingTeam := range existingTeams {
		if existingTeam.Name == team.Name {
			return fmt.Errorf("team with name %s already exists", team.Name)
		}
	}

	// Set IDs and timestamps
	team.ID = uuid.New()
	team.CreatedAt = time.Now()
	team.UpdatedAt = time.Now()

	// Save the team in the repository
	return s.teamRepo.Create(ctx, team)
}

func (s *teamService) GetTeamByID(ctx context.Context, teamID uuid.UUID) (*entity.Team, error) {
	team, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return nil, err
	}
	if team == nil {
		return nil, fmt.Errorf("team with ID %s not found", teamID)
	}
	return team, nil
}

func (s *teamService) UpdateTeam(ctx context.Context, team *entity.Team) error {
	// Validate the team entity
	if err := team.Validate(); err != nil {
		return err
	}

	// Check if the team exists
	existingTeam, err := s.teamRepo.GetByID(ctx, team.ID)
	if err != nil {
		return err
	}
	if existingTeam == nil {
		return fmt.Errorf("team with ID %s not found", team.ID)
	}

	// Update timestamps
	team.UpdatedAt = time.Now()

	// Update the team in the repository
	return s.teamRepo.Update(ctx, team)
}

func (s *teamService) DeleteTeam(ctx context.Context, teamID uuid.UUID) error {
	// Check if the team exists
	existingTeam, err := s.teamRepo.GetByID(ctx, teamID)
	if err != nil {
		return err
	}
	if existingTeam == nil {
		return fmt.Errorf("team with ID %s not found", teamID)
	}

	// Delete the team
	return s.teamRepo.Delete(ctx, teamID)
}

func (s *teamService) ListTeams(ctx context.Context) ([]*entity.Team, error) {
	return s.teamRepo.List(ctx)
}

func (s *teamService) ListTeamsByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Team, error) {
	// Check if the user exists
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, fmt.Errorf("user with ID %s not found", userID)
	}

	return s.teamRepo.GetByUserID(ctx, userID)
}
