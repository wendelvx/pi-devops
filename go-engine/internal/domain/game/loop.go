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

        // 1. LÊ A MENSAGEM PRIMEIRO (Antes de verificar se a sala existe)
        var action contracts.PlayerAction
        json.Unmarshal([]byte(msg.Payload), &action)

        // 2. Inicialização Dinâmica da Sala
        if _, exists := rooms[roomID]; !exists {
            log.Printf("🏰 Nova instância detectada: %s. Inicializando Boss...", roomID)

            targetBoss := "infra_boss" // Boss Padrão
            
            // Se a sala está sendo criada agora e o evento é 'join',
            // nós olhamos o Payload para ver se o Admin mandou um boss específico!
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
        // 1. TRATAMENTO DE RESET (Tem que ser o primeiro de todos!)
        // =======================================================
        if action.Type == "reset" {
            wasDead := currentState.Status != "fighting"
            sm.Reset() // Restaura vida e status

            // PULO DO GATO: Se o boss estava morto, o relógio dessa sala tinha parado.
            // Precisamos "reviver" a Goroutine do relógio para os incidentes voltarem!
            if wasDead {
                go startRoomGameLoop(ctx, rdb, sm, roomID)
            }

            broadcastState(ctx, rdb, sm.GetState(), roomID)
            continue
        }

        // =======================================================
        // 2. TRAVA DE FIM DE JOGO
        // =======================================================
        // Se chegou aqui e o jogo acabou, ignora ataques e resoluções.
        if currentState.Status != "fighting" {
            continue
        }

        // =======================================================
        // 3. TRAVA DE INCIDENTE GLOBAL (NOVO)
        // =======================================================
        // Se houver um incidente ativo, o sistema "trava" ataques normais.
        // Apenas resoluções ou ataques da classe que pode resolver são permitidos!
        if currentState.ActiveIncident != nil && action.Type == "attack" {
            // Se o aluno que tentou atacar NÃO é a classe responsável por resolver...
            if action.Class != currentState.ActiveIncident.TargetClass {
                continue // Ignora o ataque solenemente! O servidor está travado.
            }
        }

        // --- TRATAMENTO DE ATAQUE ---
        if action.Type == "attack" {
            dmg := CalculateDamage(action.Class, currentState.CurrentBoss.Weakness)
            msgAction := fmt.Sprintf("%s atacou!", action.Nickname)

            sm.ApplyDamage(dmg, msgAction)

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
                sm.ApplyDamage(points, fmt.Sprintf("✨ %s: %s", action.Nickname, logMsg))
                sm.ClearIncident()

                go db.Exec(`UPDATE rankings SET incidents_solved = incidents_solved + 1 WHERE nickname = ?`, action.Nickname)
            }
        }

        broadcastState(ctx, rdb, sm.GetState(), roomID)
    }
}

// O Coração da Sala (Roda a cada 1 segundo)
func startRoomGameLoop(ctx context.Context, rdb *redis.Client, sm *StateManager, roomID string) {
    ticker := time.NewTicker(1 * time.Second) // Relógio de 1 segundo!
    chaosCounter := 0

    for range ticker.C {
        state := sm.GetState()
        
        // Se o Boss morrer, o relógio para e a rotina morre (poupa RAM)
        if state.Status != "fighting" {
            ticker.Stop()
            return
        }

        needsBroadcast := false

        // 1. Se tem um incidente ativo, faz a contagem regressiva
        if state.ActiveIncident != nil {
            expired := sm.TickIncidentTimer()
            needsBroadcast = true

            if expired {
                sm.ApplyDamage(-150, "⚠️ Tempo Esgotado! A falha curou o Boss!")
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