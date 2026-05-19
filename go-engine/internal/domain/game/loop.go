package game

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"pi-devops/internal/domain/boss"
	"pi-devops/internal/domain/contracts"
)

func StartGameEngine(db *gorm.DB, rdb *redis.Client) {
	ctx := context.Background()
	rooms := make(map[string]*StateManager)

	pubsub := rdb.PSubscribe(ctx, "room:*:attacks")
	ch := pubsub.Channel()

	log.Println("⚔️ Dungeon Master Multi-Room Engine: Operacional e ouvindo padrões 'room:*:attacks'")

	for msg := range ch {
		parts := strings.Split(msg.Channel, ":")
		if len(parts) < 2 {
			continue
		}
		roomID := parts[1]

		var action contracts.PlayerAction
		json.Unmarshal([]byte(msg.Payload), &action)

		if _, exists := rooms[roomID]; !exists {
			log.Printf("🏰 Nova instância detectada: %s. Inicializando Boss...", roomID)

			targetBoss := "infra_boss"

			if action.Type == "join" && action.Payload != "" {
				targetBoss = action.Payload
			}

			professor := boss.GetProfessorBoss(targetBoss)
			rooms[roomID] = NewStateManager(professor)

			go startRoomGameLoop(ctx, rdb, rooms[roomID], roomID)
		}

		sm := rooms[roomID]
		currentState := sm.GetState()

		// =======================================================
		// NOVO: TRATAMENTO DE DESTRUIÇÃO TOTAL DA SALA
		// =======================================================
		if action.Type == "delete" {
			sm.mu.Lock()
			sm.state.Status = "deleted" // Muda o status para a Goroutine morrer
			sm.mu.Unlock()

			// Remove a sala do map na memória RAM para liberar os recursos do servidor
			delete(rooms, roomID)
			log.Printf("🗑️ Engine: Sala [%s] removida da memória com sucesso.", roomID)
			continue
		}

		// =======================================================
		// 1. TRATAMENTO DE RESET
		// =======================================================
		if action.Type == "reset" {
			wasDead := currentState.Status == "victory" || currentState.Status == "defeat"
			sm.Reset()

			if wasDead {
				go startRoomGameLoop(ctx, rdb, sm, roomID)
			}

			broadcastState(ctx, rdb, sm.GetState(), roomID)
			continue
		}

		// =======================================================
		// 2. TRATAMENTO DO START MANUAL
		// =======================================================
		if action.Type == "start_battle" && currentState.Status == "waiting" {
			sm.StartBattle()
			broadcastState(ctx, rdb, sm.GetState(), roomID)
			continue
		}

		// =======================================================
		// REFRESH DE PLAYERS ENQUANTO AGUARDA
		// Garante que a web receba a atualização de conexão
		// =======================================================
		if action.Type == "join" && currentState.Status == "waiting" {
			broadcastState(ctx, rdb, sm.GetState(), roomID)
		}

		// =======================================================
		// TRAVA DE FIM DE JOGO E ESPERA
		// =======================================================
		if currentState.Status != "fighting" {
			continue // Ignora ataques e resoluções se não estiver lutando
		}

		// =======================================================
		// 3. TRAVA DE INCIDENTE GLOBAL
		// =======================================================
		if currentState.ActiveIncident != nil && action.Type == "attack" {
			if action.Class != currentState.ActiveIncident.TargetClass {
				continue
			}
		}

		if action.Type == "attack" {
			dmg := CalculateDamage(action.Class, currentState.CurrentBoss.Weakness)
			msgAction := fmt.Sprintf("%s atacou!", action.Nickname)

			sm.ApplyDamage(dmg, msgAction, action.Nickname, action.Class)

			go db.Exec(`INSERT INTO rankings (nickname, class, total_damage, battle_id) 
                VALUES (?, ?, ?, 1) ON CONFLICT (nickname) 
                DO UPDATE SET total_damage = rankings.total_damage + EXCLUDED.total_damage`,
				action.Nickname, action.Class, dmg)
		}

		if action.Type == "resolve" {
			success, logMsg := ValidateResolution(currentState.ActiveIncident, action.Class, action.Payload)

			if success {
				points := currentState.ActiveIncident.Points
				sm.ApplyDamage(points, fmt.Sprintf("✨ %s: %s", action.Nickname, logMsg), action.Nickname, action.Class)

				sm.IncrementIncidentResolution()
				newState := sm.GetState()

				if newState.ActiveIncident != nil && newState.ActiveIncident.CurrentResolutions >= newState.ActiveIncident.RequiredResolutions {
					sm.ClearIncident()
					sm.ApplyDamage(0, "🛡️ A EQUIPE NEUTRALIZOU O INCIDENTE!", "", "")
				}

				go db.Exec(`UPDATE rankings SET incidents_solved = incidents_solved + 1 WHERE nickname = ?`, action.Nickname)
			}
		}

		broadcastState(ctx, rdb, sm.GetState(), roomID)
	}
}

func startRoomGameLoop(ctx context.Context, rdb *redis.Client, sm *StateManager, roomID string) {
	ticker := time.NewTicker(1 * time.Second)
	chaosCounter := 0

	for range ticker.C {
		state := sm.GetState()

		// Se a equipe ganhou, perdeu ou a sala foi DELETADA, desliga o loop
		if state.Status == "victory" || state.Status == "defeat" || state.Status == "deleted" {
			ticker.Stop()
			return
		}

		// =======================================================
		// TRAVA DE ESTADO "WAITING"
		// Pula o segundo inteiro sem causar dano ou incidentes
		// =======================================================
		if state.Status == "waiting" {
			continue
		}

		needsBroadcast := false

		sm.DealDamageToTeam(75, "O Servidor está sofrendo degradação passiva...")
		needsBroadcast = true

		if state.ActiveIncident != nil {
			expired := sm.TickIncidentTimer()

			if expired {
				sm.DealDamageToTeam(2500, "⚠️ INCIDENTE NÃO RESOLVIDO! A Equipe sofreu dano crítico!")
				sm.ClearIncident()
			}
		} else {
			chaosCounter++
			if chaosCounter >= 45 {
				newIncident := GenerateIncident()
				sm.SetIncident(&newIncident)
				chaosCounter = 0
				needsBroadcast = true
			}
		}

		if needsBroadcast {
			broadcastState(ctx, rdb, sm.GetState(), roomID)
		}
	}
}

func broadcastState(ctx context.Context, rdb *redis.Client, state contracts.GameState, roomID string) {
	payload, _ := json.Marshal(state)
	rdb.Publish(ctx, fmt.Sprintf("room:%s:boss_updates", roomID), payload)
}