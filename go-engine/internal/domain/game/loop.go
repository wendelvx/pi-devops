package game

import (
	"log"
	"math/rand"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/wendelvx/pi-devops.git/internal/db"
	"github.com/wendelvx/pi-devops.git/internal/metrics"
)

func StartBossLoop(rdb *redis.Client) {
	for {
		time.Sleep(State.Boss.AttackSpeed)

		State.Mu.Lock()
		if State.Status == "BATTLE" {
			dano := State.Boss.BaseDamage

			if State.ActiveIncident.ID == "legacy_code_spill" {
				dano *= 2
			}

			State.TeamHP -= dano

			if State.TeamHP <= 0 {
				State.TeamHP = 0
				State.Status = "GAMEOVER"

				db.SaveBattleResult(
					State.Status,
					State.Boss.ID,
					State.StartTime,
					State.Players,
				)

				log.Printf("💀 DERROTA: O Boss %s venceu a turma.", State.Boss.Name)
			}

			metrics.TeamHPMetric.Set(float64(State.TeamHP))
			BroadcastState(rdb)
		}
		State.Mu.Unlock()
	}
}

func StartChaosGenerator(rdb *redis.Client) {
	for {
		time.Sleep(time.Duration(rand.Intn(12)+12) * time.Second)

		incidentTypes, _ := PossibleIncident()

		if len(incidentTypes) == 0 {
			continue
		}

		randomIndex := rand.Intn(len(incidentTypes))
		selectedIncident := incidentTypes[randomIndex]

		State.Mu.Lock()
		if State.Status == "BATTLE" && State.ActiveIncident.ID == "" {
			State.ActiveIncident = selectedIncident

			log.Printf("⚠️ NOVO INCIDENTE: %s! (Requer: %s)\n", State.ActiveIncident.ID, State.ActiveIncident.RequiredClass)
			BroadcastState(rdb)
		}
		State.Mu.Unlock()
	}
}
