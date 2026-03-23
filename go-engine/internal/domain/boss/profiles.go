package boss

import "sync"

type GameState struct {
	mu              sync.RWMutex
	BossHP          int
	Multiplicadores map[string]float64
	ActiveIncident  string
}
