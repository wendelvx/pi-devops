package game

import (
	"math/rand"
	"pi-devops/internal/domain/contracts"
)

// CalculateDamage define o montante de HP subtraído do Boss.
// Implementa o Objetivo 1.1: o dano é derivado da coordenação e fraquezas técnicas.
func CalculateDamage(playerClass string, bossWeakness string) int {
	baseDamage := 10

	// 1. Bônus de Classe Dominante (Fraqueza do Boss)
	// Se a classe do aluno for a fraqueza do professor, dano massivo (RF04)
	if playerClass == bossWeakness {
		return baseDamage * 5 // Multiplicador de 500%
	}

	// 2. Chance de "Critical Hit" Aleatório (Diferencial Premium)
	// Pequena chance (5%) de qualquer classe dar dano dobrado por "excelência técnica"
	if rand.Intn(100) < 5 {
		return baseDamage * 2
	}

	return baseDamage
}

// ValidateResolution verifica se a tentativa de resolução do incidente é válida (RF05).
// Retorna um booleano e uma mensagem descritiva para o log (RF11).
func ValidateResolution(activeIncident *contracts.Incident, playerClass string, payload string) (bool, string) {
	// Proteção contra chamadas sem incidentes ativos
	if activeIncident == nil {
		return false, "Não há incidentes críticos no sistema."
	}

	// 1. Validação de Responsabilidade (TargetClass)
	// Garante que um DevOps não resolva um problema de Security
	if activeIncident.TargetClass != playerClass {
		return false, "Ação negada: Você não tem permissão para esta camada da infra!"
	}

	// 2. Validação de Lógica (Solution)
	// Compara o que o aluno fez (Payload) com o que a Engine espera
	if activeIncident.Solution != payload {
		return false, "Solução inválida: O erro persiste no log."
	}

	return true, "Incidente neutralizado com sucesso!"
}