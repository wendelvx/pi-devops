package boss

import "sync"

type GameState struct {
	mu             sync.RWMutex
	BossHP         float64
	Multipliers    map[string]float64
	ActiveIncident string
}

func NewGameState(initialHP float64) *GameState {
	return &GameState{
		BossHP:         initialHP,
		Multipliers:    make(map[string]float64),
		ActiveIncident: "",
	}
}

func (g *GameState) ApplyDamage(baseDamage float64) {
	g.mu.Lock()
	defer g.mu.Unlock()

	mult, set := g.Multipliers["damage"]

	if !set {
		mult = 1.0
	}

	currentDamage := baseDamage * mult
	g.BossHP -= currentDamage

	if g.BossHP < 0 {
		g.BossHP = 0
	}
}
