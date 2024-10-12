package service

import (
	"champi-maker/internal/domain/entity"
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func filterMatchesByPhase(matches []*entity.Match, phase int) []*entity.Match {
	var filtered []*entity.Match
	for _, match := range matches {
		if match.Phase == phase {
			filtered = append(filtered, match)
		}
	}
	return filtered
}

func TestGenerateCupMatches_Fixed_4Teams(t *testing.T) {
	ctx := context.Background()
	championshipID := uuid.New()
	championship := &entity.Championship{
		ID:              championshipID,
		Name:            "Copa Teste",
		Type:            entity.ChampionshipTypeCup,
		ProgressionType: entity.ProgressionFixed,
	}

	// Criar 4 times
	teamIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New()}

	matchService := &matchService{}

	matches, err := matchService.generateCupMatches(ctx, championship, teamIDs)
	assert.NoError(t, err)
	assert.NotNil(t, matches)

	// Verificar o número de fases e partidas
	// Com 4 times, devemos ter 2 fases:
	// - Fase 1: 2 partidas
	// - Fase 2: 1 partida (final)

	// Verificar o número total de partidas
	assert.Equal(t, 3, len(matches), "Deve haver 3 partidas no total")

	// Separar as partidas por fase
	phase1Matches := filterMatchesByPhase(matches, 1)
	phase2Matches := filterMatchesByPhase(matches, 2)

	assert.Equal(t, 2, len(phase1Matches), "Deve haver 2 partidas na Fase 1")
	assert.Equal(t, 1, len(phase2Matches), "Deve haver 1 partida na Fase 2")

	// Verificar as relações de parentesco
	finalMatch := phase2Matches[0]
	for _, match := range phase1Matches {
		assert.Equal(t, &finalMatch.ID, match.ParentMatchID, "A partida da Fase 1 deve ter ParentMatchID igual ao ID da partida final")
	}

	// Verificar que a partida final tem LeftChildMatchID e RightChildMatchID corretos
	assert.NotNil(t, finalMatch.LeftChildMatchID)
	assert.NotNil(t, finalMatch.RightChildMatchID)
	assert.Contains(t, []uuid.UUID{phase1Matches[0].ID, phase1Matches[1].ID}, *finalMatch.LeftChildMatchID)
	assert.Contains(t, []uuid.UUID{phase1Matches[0].ID, phase1Matches[1].ID}, *finalMatch.RightChildMatchID)
}

func TestGenerateCupMatches_Fixed_6Teams(t *testing.T) {
	ctx := context.Background()
	championshipID := uuid.New()
	championship := &entity.Championship{
		ID:              championshipID,
		Name:            "Copa Teste",
		Type:            entity.ChampionshipTypeCup,
		ProgressionType: entity.ProgressionFixed,
	}

	// Criar 6 times
	teamIDs := []uuid.UUID{uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New(), uuid.New()}

	matchService := &matchService{}

	matches, err := matchService.generateCupMatches(ctx, championship, teamIDs)
	assert.NoError(t, err)
	assert.NotNil(t, matches)

	// Com 6 times, devemos ter 3 fases:
	// - Fase 1: 4 partidas (algumas com byes)
	// - Fase 2: 2 partidas
	// - Fase 3: 1 partida (final)

	assert.Equal(t, 7, len(matches), "Deve haver 7 partidas no total")

	// Filtrar as partidas por fase
	phase1Matches := filterMatchesByPhase(matches, 1)
	phase2Matches := filterMatchesByPhase(matches, 2)
	phase3Matches := filterMatchesByPhase(matches, 3)

	assert.Equal(t, 4, len(phase1Matches), "Deve haver 4 partidas na Fase 1")
	assert.Equal(t, 2, len(phase2Matches), "Deve haver 2 partidas na Fase 2")
	assert.Equal(t, 1, len(phase3Matches), "Deve haver 1 partida na Fase 3")

	// Verificar relações de parentesco
	// Partidas da Fase 1 devem ter ParentMatchID apontando para as partidas da Fase 2
	// Partidas da Fase 2 devem ter ParentMatchID apontando para a partida da Fase 3
	finalMatch := phase3Matches[0]

	for _, match := range phase2Matches {
		assert.Equal(t, &finalMatch.ID, match.ParentMatchID, "A partida da Fase 2 deve ter ParentMatchID igual ao ID da partida final")
	}

	for _, match := range phase1Matches {
		parentMatchID := match.ParentMatchID
		assert.NotNil(t, parentMatchID, "Partida da Fase 1 deve ter ParentMatchID")
		// Verificar se o ParentMatchID existe nas partidas da Fase 2
		found := false
		for _, m := range phase2Matches {
			if m.ID == *parentMatchID {
				found = true
				break
			}
		}
		assert.True(t, found, "ParentMatchID da partida da Fase 1 deve existir nas partidas da Fase 2")
	}
}

func TestGenerateCupMatches_Fixed_8Teams(t *testing.T) {
	ctx := context.Background()
	championshipID := uuid.New()
	championship := &entity.Championship{
		ID:              championshipID,
		Name:            "Copa 8 Times",
		Type:            entity.ChampionshipTypeCup,
		ProgressionType: entity.ProgressionFixed,
	}

	// Criar 8 times
	teamIDs := []uuid.UUID{}
	for i := 0; i < 8; i++ {
		teamIDs = append(teamIDs, uuid.New())
	}

	matchService := &matchService{}

	matches, err := matchService.generateCupMatches(ctx, championship, teamIDs)
	assert.NoError(t, err)
	assert.NotNil(t, matches)

	// Com 8 times, devemos ter 3 fases:
	// - Fase 1: 4 partidas
	// - Fase 2: 2 partidas
	// - Fase 3: 1 partida (final)

	assert.Equal(t, 7, len(matches), "Deve haver 7 partidas no total")

	phase1Matches := filterMatchesByPhase(matches, 1)
	phase2Matches := filterMatchesByPhase(matches, 2)
	phase3Matches := filterMatchesByPhase(matches, 3)

	assert.Equal(t, 4, len(phase1Matches), "Deve haver 4 partidas na Fase 1")
	assert.Equal(t, 2, len(phase2Matches), "Deve haver 2 partidas na Fase 2")
	assert.Equal(t, 1, len(phase3Matches), "Deve haver 1 partida na Fase 3")

	// Verificar as relações de parentesco de forma semelhante aos testes anteriores
	finalMatch := phase3Matches[0]

	for _, match := range phase2Matches {
		assert.Equal(t, &finalMatch.ID, match.ParentMatchID, "A partida da Fase 2 deve ter ParentMatchID igual ao ID da partida final")
	}

	for _, match := range phase1Matches {
		parentMatchID := match.ParentMatchID
		assert.NotNil(t, parentMatchID, "Partida da Fase 1 deve ter ParentMatchID")
		found := false
		for _, m := range phase2Matches {
			if m.ID == *parentMatchID {
				found = true
				break
			}
		}
		assert.True(t, found, "ParentMatchID da partida da Fase 1 deve existir nas partidas da Fase 2")
	}

	// Verificar que a partida final tem LeftChildMatchID e RightChildMatchID corretos
	assert.NotNil(t, finalMatch.LeftChildMatchID)
	assert.NotNil(t, finalMatch.RightChildMatchID)
}
