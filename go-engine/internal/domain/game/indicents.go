package game

import (
    "fmt"
    "math/rand"
    "pi-devops/internal/domain/contracts"
)

// GenerateIncident sorteia um desafio temático baseado nas classes do TDE (RF05)
func GenerateIncident() contracts.Incident {
    incidents := []contracts.Incident{
        // --- SECURITY ---
        {
            ID:          "sec_brute_force",
            Title:       "Tentativa de Brute Force",
            Description: "Múltiplos logins falhos! Security, ative o WAF na sequência correta!",
            TargetClass: "security",
            Type:        "SEQUENCE",
            Solution:    "1-3-2",
            Points:      350,
            Duration:    15, // Adicionado
        },
        {
            ID:          "sec_sql_inj",
            Title:       "SQL Injection Detectado",
            Description: "Alguém tentou um 'OR 1=1'! Sanitize os inputs agora!",
            TargetClass: "security",
            Type:        "PUZZLE",
            Solution:    "PREPARED_STMT",
            Points:      450,
            Duration:    15, // Adicionado
        },

        // --- DEVOPS ---
        {
            ID:          "dev_high_cpu",
            Title:       "Pico de Processamento",
            Description: "CPU em 99%! DevOps, escale o cluster manualmente!",
            TargetClass: "devops",
            Type:        "RAPID_CLICK",
            Solution:    "15",
            Points:      300,
            Duration:    10, // Menos tempo, precisa ser rápido!
        },
        {
            ID:          "dev_pod_eviction",
            Title:       "Pod Eviction",
            Description: "Nós sem memória! Mova as réplicas para o novo Node!",
            TargetClass: "devops",
            Type:        "SEQUENCE",
            Solution:    "2-2-1",
            Points:      400,
            Duration:    15,
        },

        // --- BACK-END ---
        {
            ID:          "back_deadlock",
            Title:       "Database Deadlock",
            Description: "Transações travadas! Back-end, libere os Mutexes!",
            TargetClass: "back-end",
            Type:        "SEQUENCE",
            Solution:    "1-2-3",
            Points:      300,
            Duration:    15,
        },

        // --- FRONT-END ---
        {
            ID:          "front_layout_shift",
            Title:       "Cumulative Layout Shift",
            Description: "A UI está quebrando no Safari! Front-end, fixe o Z-Index!",
            TargetClass: "front-end",
            Type:        "RAPID_CLICK",
            Solution:    "12",
            Points:      250,
            Duration:    10,
        },

        // --- QA / TESTER ---
        {
            ID:          "qa_regression",
            Title:       "Bug em Produção",
            Description: "Regressão detectada na Pipeline! QA, rode o Smoke Test!",
            TargetClass: "qa",
            Type:        "PUZZLE",
            Solution:    "RUN_SMOKE",
            Points:      350,
            Duration:    20, // Puzzles precisam de mais tempo para ler
        },
    }

    selected := incidents[rand.Intn(len(incidents))]
    
    fmt.Printf("⚠️ INCIDENTE ATIVO: %s para a classe %s (Duração: %ds)\n", selected.Title, selected.TargetClass, selected.Duration)
    
    return selected
}