package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
    "os"
    "time"

    "github.com/redis/go-redis/v9"
    "pi-devops/internal/db"
    "pi-devops/internal/domain/boss"
    "pi-devops/internal/domain/game"
    metrics "pi-devops/internal/metrics"
)

func main() {
    startedAt := time.Now()
    
    // 1. Inicializa Conexão com Postgres
    database := db.ConectDb()
    log.Println("✅ Banco de Dados Postgres conectado com sucesso!")

    // 2. Inicializa Cliente Redis (Pub/Sub)
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "redis:6379"
    }
    rdb := redis.NewClient(&redis.Options{
        Addr: redisAddr,
    })
    log.Printf("🔌 Redis conectado em: %s", redisAddr)

    // 3. Dispara a Game Engine (O Mestre da Masmorra)
    // Toda a lógica de Salas, Ataques e o novo RESET rodam em background a partir daqui!
    go game.StartGameEngine(database, rdb)

    // 4. Configuração do Servidor HTTP (API REST Auxiliar)
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    mux := http.NewServeMux()
    
    // --- Helper para Headers CORS (Essencial para o Front Web) ---
    enableCors := func(w http.ResponseWriter) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    }

    // --- Endpoint: Healthcheck ---
    mux.HandleFunc("/ping", func(writer http.ResponseWriter, request *http.Request) {
        enableCors(writer)
        writer.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(writer, `{"service":"go-engine","status":"pong"}`)
    })

    // --- Endpoint: Métricas (Monitoramento via Prometheus) ---
    mux.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
        writer.Header().Set("Content-Type", "text/plain; version=0.0.4")
        fmt.Fprint(writer, metrics.Render(startedAt))
    })

    // --- Endpoint: Info do Boss Atual ---
    mux.HandleFunc("/boss", func(writer http.ResponseWriter, request *http.Request) {
        enableCors(writer)
        if request.Method == "OPTIONS" { return } // Preflight do navegador

        currentBoss := boss.GetProfessorBoss("infra_boss") 
        
        writer.Header().Set("Content-Type", "application/json")
        json.NewEncoder(writer).Encode(currentBoss)
    })

    // --- Endpoint: Ranking Final (RF06) ---
    mux.HandleFunc("/ranking", func(writer http.ResponseWriter, request *http.Request) {
        enableCors(writer)
        if request.Method == "OPTIONS" { return }

        type RankResult struct {
            Nickname        string `json:"nickname"`
            Class           string `json:"class"`
            TotalDamage     int    `json:"total_damage"`
            IncidentsSolved int    `json:"incidents_solved"`
        }

        var rankings []RankResult
        database.Raw(`
            SELECT nickname, class, total_damage, incidents_solved 
            FROM rankings 
            ORDER BY total_damage DESC LIMIT 10
        `).Scan(&rankings)

        writer.Header().Set("Content-Type", "application/json")
        json.NewEncoder(writer).Encode(rankings)
    })

    // --- Painel do Professor ---
    mux.HandleFunc("/admin/rooms", func(writer http.ResponseWriter, request *http.Request) {
        enableCors(writer)
        if request.Method == "OPTIONS" { return }

        writer.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(writer, `{"status":"pending","message":"Endpoint de criação de sala em construção!"}`)
    })

    address := "0.0.0.0:" + port
    log.Printf("🚀 Dungeon Master API & Engine prontas em %s", address)
    
    if err := http.ListenAndServe(address, mux); err != nil {
        log.Fatal("❌ Erro no servidor HTTP:", err)
    }
}