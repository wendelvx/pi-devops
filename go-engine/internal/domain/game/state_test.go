package game

import (
	"testing"

	"pi-devops/internal/domain/contracts"

	"github.com/stretchr/testify/assert"
)

func createTestBoss() contracts.Boss {
	return contracts.Boss{
		ID:    "test_boss",
		Name:  "Boss de Teste",
		MaxHP: 1000,
	}
}

// TESTES DE INICIALIZAÇÃO E ESTADO
func TestNewStateManager_InitialState(t *testing.T) {
	boss := createTestBoss()
	sm := NewStateManager(boss)

	state := sm.GetState()

	assert.Equal(t, boss.MaxHP, state.BossHP)
	assert.Equal(t, 10000, state.TeamHP)
	assert.Equal(t, "waiting", state.Status)
	assert.Equal(t, boss, state.CurrentBoss)
	assert.Contains(t, state.LastAction, "Aguardando")
	assert.NotNil(t, sm.damageBoard)
}

func TestStartBattle_ChangesStatus(t *testing.T) {
	sm := NewStateManager(createTestBoss())

	sm.StartBattle()
	state := sm.GetState()

	assert.Equal(t, "fighting", state.Status)
	assert.Contains(t, state.LastAction, "A BATALHA COMEÇOU")
}

// TESTES DE SISTEMA DE DANO E RANKING
func TestApplyDamage_IgnoredIfNotFighting(t *testing.T) {
	sm := NewStateManager(createTestBoss()) // Status inicial é "waiting"

	sm.ApplyDamage(100, "Ataque", "player1", "devops")
	state := sm.GetState()

	// O HP do boss não deve ter mudado porque a batalha não começou
	assert.Equal(t, createTestBoss().MaxHP, state.BossHP)
}

func TestApplyDamage_ReducesBossHP_AndTracksStats(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	sm.StartBattle()

	sm.ApplyDamage(100, "Player1 atacou!", "player1", "devops")
	state := sm.GetState()

	assert.Equal(t, 900, state.BossHP)
	assert.Equal(t, "Player1 atacou!", state.LastAction)

	assert.Equal(t, 100, sm.damageBoard["player1"].Damage)
}

func TestApplyDamage_CapsBossHPAtMax(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	sm.StartBattle()

	sm.ApplyDamage(-500, "Boss se curou!", "", "")
	state := sm.GetState()

	assert.Equal(t, createTestBoss().MaxHP, state.BossHP)
}

func TestApplyDamage_VictoryAndRanking(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	sm.StartBattle()

	sm.ApplyDamage(10, "P1", "player1", "devops")
	sm.ApplyDamage(40, "P2", "player2", "security")
	sm.ApplyDamage(20, "P3", "player3", "back-end")
	sm.ApplyDamage(30, "P4", "player4", "front-end")

	sm.ApplyDamage(900, "Ataque Fatal!", "player2", "security")

	state := sm.GetState()

	assert.Equal(t, 0, state.BossHP)
	assert.Equal(t, "victory", state.Status)
	assert.Equal(t, "O BOSS FOI DERROTADO!", state.LastAction)

	assert.Len(t, state.TopRank, 3, "O ranking deve limitar-se aos 3 melhores jogadores")

	assert.Equal(t, "player2", state.TopRank[0].Nickname)
	assert.Equal(t, 940, state.TopRank[0].Damage)

	assert.Equal(t, "player4", state.TopRank[1].Nickname)

	assert.Equal(t, "player3", state.TopRank[2].Nickname)
}

func TestDealDamageToTeam_DefeatCondition(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	sm.StartBattle()

	sm.DealDamageToTeam(1000, "Instabilidade menor")
	assert.Equal(t, 9000, sm.GetState().TeamHP)

	sm.DealDamageToTeam(10000, "Crash geral do servidor!")
	state := sm.GetState()

	assert.Equal(t, 0, state.TeamHP)
	assert.Equal(t, "defeat", state.Status)
	assert.Equal(t, "SISTEMA COMPROMETIDO. A EQUIPE FOI DERROTADA!", state.LastAction)
}

func TestDealDamageToTeam_IgnoredIfNotFighting(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	sm.DealDamageToTeam(500, "Dano em background")

	assert.Equal(t, 10000, sm.GetState().TeamHP)
}

// TESTES DE INCIDENTES
func TestIncidentManagement(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	incident := &contracts.Incident{
		ID:       "inc_01",
		Duration: 15,
	}

	// Testa SetIncident
	sm.SetIncident(incident)
	assert.Equal(t, incident, sm.GetState().ActiveIncident)
	assert.Equal(t, 15, sm.GetState().IncidentTimer)

	// Testa IncrementIncidentResolution
	sm.IncrementIncidentResolution()
	assert.Equal(t, 1, sm.GetState().ActiveIncident.CurrentResolutions)

	// Testa TickIncidentTimer
	expired := sm.TickIncidentTimer()
	assert.False(t, expired)
	assert.Equal(t, 14, sm.GetState().IncidentTimer)

	// Força a expiração para testar o retorno booleano
	sm.mu.Lock()
	sm.state.IncidentTimer = 1
	sm.mu.Unlock()

	expired = sm.TickIncidentTimer()
	assert.True(t, expired)

	// Testa ClearIncident
	sm.ClearIncident()
	assert.Nil(t, sm.GetState().ActiveIncident)
	assert.Equal(t, 0, sm.GetState().IncidentTimer)
}

func TestTickIncidentTimer_ReturnsFalseIfNoActiveIncident(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	expired := sm.TickIncidentTimer()
	assert.False(t, expired)
}

// TESTE DE REINICIALIZAÇÃO
func TestReset_RestoresInitialValues(t *testing.T) {
	sm := NewStateManager(createTestBoss())
	sm.StartBattle()

	sm.ApplyDamage(100, "Dano", "p1", "qa")
	sm.SetIncident(&contracts.Incident{ID: "inc"})

	sm.Reset()
	state := sm.GetState()

	assert.Equal(t, createTestBoss().MaxHP, state.BossHP)
	assert.Equal(t, 10000, state.TeamHP)
	assert.Equal(t, "waiting", state.Status)
	assert.Nil(t, state.ActiveIncident)
	assert.Equal(t, 0, state.IncidentTimer)
	assert.Nil(t, state.TopRank)
	assert.Empty(t, sm.damageBoard)
}
