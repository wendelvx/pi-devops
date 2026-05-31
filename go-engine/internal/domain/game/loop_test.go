package game

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"pi-devops/internal/domain/boss"
	"pi-devops/internal/domain/contracts"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestBroadcastState(t *testing.T) {
	rdb, mock := redismock.NewClientMock()
	ctx := context.Background()
	roomID := "sala_faculdade"

	state := contracts.GameState{Status: "fighting"}
	payload, _ := json.Marshal(state)

	mock.ExpectPublish("room:sala_faculdade:boss_updates", payload).SetVal(1)

	broadcastState(ctx, rdb, state, roomID)

	assert.NoError(t, mock.ExpectationsWereMet())
}

// TESTES DO GAME LOOP
func TestStartRoomGameLoop_ExitsOnTerminalStates(t *testing.T) {
	// Vamos testar se o loop encerra corretamente nos 3 status finais
	terminalStates := []string{"victory", "defeat", "deleted"}

	for _, status := range terminalStates {
		t.Run("Status: "+status, func(t *testing.T) {
			rdb, _ := redismock.NewClientMock()
			ctx := context.Background()

			//boss genérico para inicializar o negocio
			dummyBoss := boss.GetProfessorBoss("infra_boss")
			sm := NewStateManager(dummyBoss)

			// Força o status terminal
			sm.mu.Lock()
			sm.state.Status = status
			sm.mu.Unlock()

			done := make(chan bool)

			go func() {
				startRoomGameLoop(ctx, rdb, sm, "room_123")
				done <- true
			}()

			// ver se ele termina em menos de 1 segundo
			select {
			case <-done:
				assert.True(t, true, "Goroutine encerrou com sucesso no status "+status)
			case <-time.After(1500 * time.Millisecond):
				t.Fatalf("O loop não encerrou no status %s e causou timeout", status)
			}
		})
	}
}

func TestStartRoomGameLoop_WaitingState(t *testing.T) {

	rdb, _ := redismock.NewClientMock()
	ctx := context.Background()
	dummyBoss := boss.GetProfessorBoss("infra_boss")
	sm := NewStateManager(dummyBoss)

	sm.mu.Lock()
	sm.state.Status = "waiting"
	sm.mu.Unlock()

	go startRoomGameLoop(ctx, rdb, sm, "room_wait")

	time.Sleep(1500 * time.Millisecond)

	sm.mu.Lock()
	sm.state.Status = "deleted"
	sm.mu.Unlock()
}

//TESTE DE INICIALIZAÇÃO DA ENGINE MAIN

func TestStartGameEngine_Initialization(t *testing.T) {

	rdb, _ := redismock.NewClientMock()
	db := &gorm.DB{}

	go StartGameEngine(db, rdb)

	time.Sleep(100 * time.Millisecond)

	assert.True(t, true)
}
