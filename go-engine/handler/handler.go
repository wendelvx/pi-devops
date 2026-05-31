package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"pi-devops/internal/domain/boss"
	metrics "pi-devops/internal/metrics"

	"gorm.io/gorm"
)

func SetupRouter(database *gorm.DB, startedAt time.Time) *http.ServeMux {
	mux := http.NewServeMux()

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

	// --- Endpoint: Métricas ---
	mux.HandleFunc("/metrics", func(writer http.ResponseWriter, request *http.Request) {
		writer.Header().Set("Content-Type", "text/plain; version=0.0.4")
		fmt.Fprint(writer, metrics.Render(startedAt))
	})

	// --- Endpoint: Info do Boss Atual ---
	mux.HandleFunc("/boss", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(writer)
		if request.Method == "OPTIONS" {
			return
		}

		currentBoss := boss.GetProfessorBoss("infra_boss")
		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(currentBoss)
	})

	// --- Endpoint: Ranking Final ---
	mux.HandleFunc("/ranking", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(writer)
		if request.Method == "OPTIONS" {
			return
		}

		type RankResult struct {
			Nickname        string `json:"nickname"`
			Class           string `json:"class"`
			TotalDamage     int    `json:"total_damage"`
			IncidentsSolved int    `json:"incidents_solved"`
		}

		var rankings []RankResult
		if database != nil {
			database.Raw(`
				SELECT nickname, class, total_damage, incidents_solved 
				FROM rankings 
				ORDER BY total_damage DESC LIMIT 10
			`).Scan(&rankings)
		}

		writer.Header().Set("Content-Type", "application/json")
		json.NewEncoder(writer).Encode(rankings)
	})

	// --- Painel do Professor ---
	mux.HandleFunc("/admin/rooms", func(writer http.ResponseWriter, request *http.Request) {
		enableCors(writer)
		if request.Method == "OPTIONS" {
			return
		}

		writer.Header().Set("Content-Type", "application/json")
		fmt.Fprintf(writer, `{"status":"pending","message":"Endpoint de criação de sala em construção!"}`)
	})

	return mux
}
