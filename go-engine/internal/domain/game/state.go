package game

import (
	"pi-devops/internal/domain/contracts"
	"sync"
)

type StateManager struct {
	mu          sync.RWMutex
	state       contracts.GameState
	// NOVO: Um mapa em memória para guardar o dano de cada aluno durante a partida
	damageBoard map[string]contracts.PlayerStat
}

func NewStateManager(initialBoss contracts.Boss) *StateManager {
	return &StateManager{
		damageBoard: make(map[string]contracts.PlayerStat),
		state: contracts.GameState{
			BossHP:      initialBoss.MaxHP,
			TeamHP:      10000, // Defina um HP alto para a equipe inteira
			MaxTeamHP:   10000,
			Status:      "fighting",
			CurrentBoss: initialBoss,
			LastAction:  "A masmorra foi aberta!",
		},
	}
}

func (sm *StateManager) GetState() contracts.GameState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.state
}

// Atualizado para receber quem atacou e calcular o MVP se o boss morrer
func (sm *StateManager) ApplyDamage(damage int, action string, nickname string, class string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state.Status != "fighting" {
		return
	}

	// Subtrai do Boss
	sm.state.BossHP -= damage

	if sm.state.BossHP > sm.state.CurrentBoss.MaxHP {
		sm.state.BossHP = sm.state.CurrentBoss.MaxHP
	}

	sm.state.LastAction = action

	// NOVO: Registra o dano no placar interno
	if nickname != "" && nickname != "GAME_MASTER" && damage > 0 {
		stat := sm.damageBoard[nickname]
		stat.Nickname = nickname
		stat.Class = class
		stat.Damage += damage
		sm.damageBoard[nickname] = stat
	}

	// VITÓRIA DA EQUIPE
	if sm.state.BossHP <= 0 {
		sm.state.BossHP = 0
		sm.state.Status = "victory"
		sm.state.LastAction = "O BOSS FOI DERROTADO!"
		sm.calculateMVP()
	}
}

// NOVO: Função para dar dano na vida dos alunos (se eles não resolverem o incidente)
func (sm *StateManager) DealDamageToTeam(damage int, action string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state.Status != "fighting" {
		return
	}

	sm.state.TeamHP -= damage
	sm.state.LastAction = action

	// DERROTA DA EQUIPE
	if sm.state.TeamHP <= 0 {
		sm.state.TeamHP = 0
		sm.state.Status = "defeat"
		sm.state.LastAction = "SISTEMA COMPROMETIDO. A EQUIPE FOI DERROTADA!"
		sm.calculateMVP() // Elege quem tentou salvar mesmo na derrota
	}
}

// NOVO: Lógica interna para varrer o map e achar quem deu mais dano
func (sm *StateManager) calculateMVP() {
	var mvp *contracts.PlayerStat
	highestDamage := 0

	for _, stat := range sm.damageBoard {
		if stat.Damage > highestDamage {
			highestDamage = stat.Damage
			statCopy := stat
			mvp = &statCopy
		}
	}
	sm.state.MVP = mvp
}

func (sm *StateManager) SetIncident(incident *contracts.Incident) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.ActiveIncident = incident
	if incident != nil {
		sm.state.IncidentTimer = incident.Duration
	}
}

func (sm *StateManager) ClearIncident() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.ActiveIncident = nil
	sm.state.IncidentTimer = 0
}

// Atualizado para limpar o placar e restaurar a vida da equipe
func (sm *StateManager) Reset() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.state.BossHP = sm.state.CurrentBoss.MaxHP
	sm.state.TeamHP = sm.state.MaxTeamHP
	sm.state.Status = "fighting"
	sm.state.ActiveIncident = nil
	sm.state.IncidentTimer = 0
	sm.state.MVP = nil
	sm.state.LastAction = "A masmorra foi resetada pelo Professor!"
	
	// Limpa o placar antigo
	sm.damageBoard = make(map[string]contracts.PlayerStat)
}

func (sm *StateManager) TickIncidentTimer() bool {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state.ActiveIncident != nil {
		sm.state.IncidentTimer--
		if sm.state.IncidentTimer <= 0 {
			return true
		}
	}
	return false
}