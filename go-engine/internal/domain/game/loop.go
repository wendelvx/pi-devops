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

	// Mapa para gerenciar múltiplas salas simultaneamente (RF09)
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

		// 1. LÊ A MENSAGEM PRIMEIRO
		var action contracts.PlayerAction
		json.Unmarshal([]byte(msg.Payload), &action)

		// 2. Inicialização Dinâmica da Sala
		if _, exists := rooms[roomID]; !exists {
			log.Printf("🏰 Nova instância detectada: %s. Inicializando Boss...", roomID)

			targetBoss := "infra_boss" // Boss Padrão
			
			if action.Type == "join" && action.Payload != "" {
				targetBoss = action.Payload
			}

			professor := boss.GetProfessorBoss(targetBoss)
			rooms[roomID] = NewStateManager(professor)

			// Inicia o Game Loop de 1 Segundo para ESTA sala
			go startRoomGameLoop(ctx, rdb, rooms[roomID], roomID)
		}

		sm := rooms[roomID]
		currentState := sm.GetState()

		// =======================================================
		// 1. TRATAMENTO DE RESET
		// =======================================================
		if action.Type == "reset" {
			wasDead := currentState.Status != "fighting"
			sm.Reset() 

			if wasDead {
				go startRoomGameLoop(ctx, rdb, sm, roomID)
			}

			broadcastState(ctx, rdb, sm.GetState(), roomID)
			continue
		}

		// =======================================================
		// 2. TRAVA DE FIM DE JOGO
		// =======================================================
		if currentState.Status != "fighting" {
			continue
		}

		// =======================================================
		// 3. TRAVA DE INCIDENTE GLOBAL
		// =======================================================
		if currentState.ActiveIncident != nil && action.Type == "attack" {
			if action.Class != currentState.ActiveIncident.TargetClass {
				continue 
			}
		}

		// --- TRATAMENTO DE ATAQUE ---
		if action.Type == "attack" {
			dmg := CalculateDamage(action.Class, currentState.CurrentBoss.Weakness)
			msgAction := fmt.Sprintf("%s atacou!", action.Nickname)

			// CORREÇÃO AQUI: Passando o Nickname e a Class para registrar no MVP
			sm.ApplyDamage(dmg, msgAction, action.Nickname, action.Class)

			go db.Exec(`INSERT INTO rankings (nickname, class, total_damage, battle_id) 
				VALUES (?, ?, ?, 1) ON CONFLICT (nickname) 
				DO UPDATE SET total_damage = rankings.total_damage + EXCLUDED.total_damage`,
				action.Nickname, action.Class, dmg)
		}

		// --- TRATAMENTO DE RESOLUÇÃO ---
		if action.Type == "resolve" {
			success, logMsg := ValidateResolution(currentState.ActiveIncident, action.Class, action.Payload)

			if success {
				points := currentState.ActiveIncident.Points
				// CORREÇÃO AQUI: Passando o Nickname e a Class
				sm.ApplyDamage(points, fmt.Sprintf("✨ %s: %s", action.Nickname, logMsg), action.Nickname, action.Class)
				sm.ClearIncident()

				go db.Exec(`UPDATE rankings SET incidents_solved = incidents_solved + 1 WHERE nickname = ?`, action.Nickname)
			}
		}

		broadcastState(ctx, rdb, sm.GetState(), roomID)
	}
}

// O Coração da Sala (Roda a cada 1 segundo)
func startRoomGameLoop(ctx context.Context, rdb *redis.Client, sm *StateManager, roomID string) {
	ticker := time.NewTicker(1 * time.Second) 
	chaosCounter := 0

	for range ticker.C {
		state := sm.GetState()
		
		// Se o jogo acabou, a rotina morre
		if state.Status != "fighting" {
			ticker.Stop()
			return
		}

		needsBroadcast := false

		// 1. Se tem um incidente ativo, faz a contagem regressiva
		if state.ActiveIncident != nil {
			expired := sm.TickIncidentTimer()
			needsBroadcast = true

			// =======================================================
			// CORREÇÃO AQUI: APLICAR DANO NA EQUIPE SE O TEMPO ZERAR
			// =======================================================
			if expired {
				// Causa um dano gigante (ex: 2500) à vida da equipe
				sm.DealDamageToTeam(2500, "⚠️ INCIDENTE NÃO RESOLVIDO! A Equipe sofreu dano crítico!")
				sm.ClearIncident()
			}
		} else {
			// 2. Se NÃO tem incidente, conta o tempo para o próximo caos
			chaosCounter++
			if chaosCounter >= 45 {
				newIncident := GenerateIncident()
				sm.SetIncident(&newIncident)
				chaosCounter = 0
				needsBroadcast = true
			}
		}

		// Envia o estado atualizado para o React Native
		if needsBroadcast {
			broadcastState(ctx, rdb, sm.GetState(), roomID)
		}
	}
}

func broadcastState(ctx context.Context, rdb *redis.Client, state contracts.GameState, roomID string) {
	payload, _ := json.Marshal(state)
	rdb.Publish(ctx, fmt.Sprintf("room:%s:boss_updates", roomID), payload)
}