package game

import (
	"pi-devops/internal/domain/contracts"
	"sync"
)

type StateManager struct {
	mu          sync.RWMutex
	state       contracts.GameState
	damageBoard map[string]contracts.PlayerStat
}

func NewStateManager(initialBoss contracts.Boss) *StateManager {
	return &StateManager{
		damageBoard: make(map[string]contracts.PlayerStat),
		state: contracts.GameState{
			BossHP:      initialBoss.MaxHP,
			TeamHP:      10000, 
			MaxTeamHP:   10000,
			Status:      "waiting", // <--- MUDANÇA: Começa esperando o Start
			CurrentBoss: initialBoss,
			LastAction:  "Aguardando o Professor iniciar a batalha...",
		},
	}
}

func (sm *StateManager) GetState() contracts.GameState {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return sm.state
}

// NOVO: Método para liberar a batalha
func (sm *StateManager) StartBattle() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.state.Status = "fighting"
	sm.state.LastAction = "🔥 A BATALHA COMEÇOU! ATAQUEM O SERVIDOR!"
}

func (sm *StateManager) ApplyDamage(damage int, action string, nickname string, class string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state.Status != "fighting" {
		return
	}

	sm.state.BossHP -= damage

	if sm.state.BossHP > sm.state.CurrentBoss.MaxHP {
		sm.state.BossHP = sm.state.CurrentBoss.MaxHP
	}

	sm.state.LastAction = action

	if nickname != "" && nickname != "GAME_MASTER" && damage > 0 {
		stat := sm.damageBoard[nickname]
		stat.Nickname = nickname
		stat.Class = class
		stat.Damage += damage
		sm.damageBoard[nickname] = stat
	}

	if sm.state.BossHP <= 0 {
		sm.state.BossHP = 0
		sm.state.Status = "victory"
		sm.state.LastAction = "O BOSS FOI DERROTADO!"
		sm.calculateMVP()
	}
}

func (sm *StateManager) DealDamageToTeam(damage int, action string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if sm.state.Status != "fighting" {
		return
	}

	sm.state.TeamHP -= damage
	sm.state.LastAction = action

	if sm.state.TeamHP <= 0 {
		sm.state.TeamHP = 0
		sm.state.Status = "defeat"
		sm.state.LastAction = "SISTEMA COMPROMETIDO. A EQUIPE FOI DERROTADA!"
		sm.calculateMVP() 
	}
}

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

func (sm *StateManager) IncrementIncidentResolution() {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	if sm.state.ActiveIncident != nil {
		sm.state.ActiveIncident.CurrentResolutions++
	}
}

func (sm *StateManager) Reset() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.state.BossHP = sm.state.CurrentBoss.MaxHP
	sm.state.TeamHP = sm.state.MaxTeamHP
	sm.state.Status = "waiting" // <--- MUDANÇA: Volta para waiting no reset
	sm.state.ActiveIncident = nil
	sm.state.IncidentTimer = 0
	sm.state.MVP = nil
	sm.state.LastAction = "A masmorra foi resetada! Aguardando o Professor iniciar..."
	
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