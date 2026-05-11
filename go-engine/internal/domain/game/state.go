package game

import (
    "pi-devops/internal/domain/contracts"
    "sync"
)

// StateManager controla o estado global da partida de forma segura
type StateManager struct {
    mu    sync.RWMutex
    state contracts.GameState
}

// NewStateManager inicia o estado com um Boss específico
func NewStateManager(initialBoss contracts.Boss) *StateManager {
    return &StateManager{
        state: contracts.GameState{
            BossHP:         initialBoss.MaxHP,
            Status:         "fighting",
            CurrentBoss:    initialBoss,
            ActiveIncident: nil,
            IncidentTimer:  0, // Relógio começa zerado
            LastAction:     "A masmorra foi aberta!",
        },
    }
}

// GetState retorna uma cópia do estado para envio via Redis
func (sm *StateManager) GetState() contracts.GameState {
    sm.mu.RLock()
    defer sm.mu.RUnlock()
    return sm.state
}

// ApplyDamage subtrai o HP (ou adiciona, se for dano negativo/cura) e verifica vitória
func (sm *StateManager) ApplyDamage(damage int, action string) {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if sm.state.Status != "fighting" {
        return
    }

    sm.state.BossHP -= damage

    // Trava de "Overheal": Se curou e passou da vida máxima, trava no limite!
    if sm.state.BossHP > sm.state.CurrentBoss.MaxHP {
        sm.state.BossHP = sm.state.CurrentBoss.MaxHP
    }

    sm.state.LastAction = action

    if sm.state.BossHP <= 0 {
        sm.state.BossHP = 0
        sm.state.Status = "victory"
        sm.state.LastAction = "O BOSS FOI DERROTADO!"
    }
}

// SetIncident define o desafio atual e INICIA O RELÓGIO
func (sm *StateManager) SetIncident(incident *contracts.Incident) {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.state.ActiveIncident = incident
    if incident != nil {
        sm.state.IncidentTimer = incident.Duration // Seta os segundos que o time tem
    }
}

// ClearIncident remove o desafio após resolução e reseta o relógio
func (sm *StateManager) ClearIncident() {
    sm.mu.Lock()
    defer sm.mu.Unlock()
    sm.state.ActiveIncident = nil
    sm.state.IncidentTimer = 0
}

func (sm *StateManager) Reset() {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    sm.state.BossHP = sm.state.CurrentBoss.MaxHP
    sm.state.Status = "fighting"
    sm.state.ActiveIncident = nil
    sm.state.IncidentTimer = 0
    sm.state.LastAction = "A masmorra foi resetada pelo Professor!"
}

// TickIncidentTimer reduz 1 segundo do relógio de forma thread-safe
// Retorna true se o tempo acabar.
func (sm *StateManager) TickIncidentTimer() bool {
    sm.mu.Lock()
    defer sm.mu.Unlock()

    if sm.state.ActiveIncident != nil {
        sm.state.IncidentTimer--
        if sm.state.IncidentTimer <= 0 {
            return true // O tempo esgotou!
        }
    }
    return false
}