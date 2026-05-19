package game

import (
    "math/rand"
    "pi-devops/internal/domain/contracts"
)

// CalculateDamage: NENHUMA MUDANÇA NECESSÁRIA AQUI.
func CalculateDamage(playerClass string, bossWeakness string) int {
    baseDamage := 10

    if playerClass == bossWeakness {
        return baseDamage * 5 
    }

    if rand.Intn(100) < 5 {
        return baseDamage * 2
    }

    return baseDamage
}

// ValidateResolution: AJUSTE APENAS A MENSAGEM FINAL
func ValidateResolution(activeIncident *contracts.Incident, playerClass string, payload string) (bool, string) {
    if activeIncident == nil {
        return false, "Não há incidentes críticos no sistema."
    }

    if activeIncident.TargetClass != playerClass {
        return false, "Ação negada: Você não tem permissão para esta camada da infra!"
    }

    if activeIncident.Solution != payload {
        return false, "Solução inválida: O erro persiste no log."
    }

    // MUDANÇA AQUI: Mensagem reflete que é um esforço coletivo
    return true, "Ação validada! Contribuindo para a resolução..."
}