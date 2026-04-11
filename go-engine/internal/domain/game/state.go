package game

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	models "github.com/wendelvx/pi-devops.git/internal/domain/contracts"
)

type GameState struct {
	Status         string                         `json:"status"`
	Boss           models.BossProfile             `json:"boss"`
	BossHP         int                            `json:"boss_hp"`
	TeamHP         float64                        `json:"team_hp"`
	MaxTeamHP      int                            `json:"max_team_hp"`
	Multiplier     float64                        `json:"multiplier"`
	ActiveIncident models.IncidentMeta            `json:"active_incident"`
	ClassCounts    map[string]int                 `json:"class_counts"`
	Players        map[string]*models.PlayerStats `json:"players"`
	ResetTrigger   int                            `json:"reset_trigger"`
	StartTime      time.Time                      `json:"start_time"`
	LastAttackAt   time.Time                      `json:"-"`
	Mu             sync.Mutex                     `json:"-"`
}

var State = GameState{
	Status: "WAITING",
	Boss: models.BossProfile{
		ID:           "arquiteto",
		Name:         "O Arquiteto (Padrão)",
		BaseHP:       1000000,
		BaseDamage:   20,
		AttackSpeed:  3 * time.Second,
		IncidentBias: "database_lock",
	},
	TeamHP:      1000,
	MaxTeamHP:   1000,
	Multiplier:  1.0,
	ClassCounts: make(map[string]int),
	Players:     make(map[string]*models.PlayerStats),
}

func BroadcastState(rdb *redis.Client) {

	payload, err := json.Marshal(&State)
	if err != nil {
		log.Printf("❌ Erro fatal ao gerar JSON do estado: %v\n", err)
		return
	}

	go func(data []byte) {
		//contexto com timeout
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		err := rdb.Publish(ctx, "boss_updates", data).Err()
		if err != nil {
			log.Printf("⚠️ Falha ao publicar no Redis: %v\n", err)
		}
	}(payload)
}
