package main

import (
	"log"
	"net/http"
	"os"
	"pi-devops/handler"
	"pi-devops/internal/db"
	"pi-devops/internal/domain/game"
	"time"

	"github.com/redis/go-redis/v9"
)

func main() {
	startedAt := time.Now()
	database := db.ConectDb()
	log.Println("✅ Banco de Dados Postgres conectado com sucesso!")

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "redis:6379"
	}
	rdb := redis.NewClient(&redis.Options{Addr: redisAddr})
	log.Printf("🔌 Redis conectado em: %s", redisAddr)

	go game.StartGameEngine(database, rdb)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Chama a função isolada de rotas
	mux := handler.SetupRouter(database, startedAt)

	address := "0.0.0.0:" + port
	log.Printf("🚀 Dungeon Master API & Engine prontas em %s", address)
	if err := http.ListenAndServe(address, mux); err != nil {
		log.Fatal("❌ Erro no servidor HTTP:", err)
	}
}
