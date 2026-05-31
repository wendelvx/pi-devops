package game

import (
	"testing"

	"pi-devops/internal/domain/contracts"

	"github.com/stretchr/testify/assert"
)

func TestCalculateDamage_WhenPlayerClassMatchesBossWeakness(t *testing.T) {
	for i := 0; i < 100; i++ {
		damage := CalculateDamage("devops", "devops")
		assert.Equal(t, 50, damage, "O dano deveria ser multiplicado por 5 quando acerta a fraqueza")
	}
}

func TestCalculateDamage_WhenPlayerClassDoesNotMatchBossWeakness(t *testing.T) {
	seenNormal := false
	seenCritical := false

	for i := 0; i < 200; i++ {
		damage := CalculateDamage("devops", "security")

		// O dano só pode ser 10 ou 20
		assert.Contains(t, []int{10, 20}, damage, "O dano deve ser 10 (normal) ou 20 (crítico)")

		if damage == 10 {
			seenNormal = true
		}
		if damage == 20 {
			seenCritical = true
		}
	}

	assert.True(t, seenNormal, "Deveria ter ocorrido pelo menos um dano normal (10)")
	assert.True(t, seenCritical, "Deveria ter ocorrido pelo menos um dano crítico (20) em 200 tentativas")
}

// --- TESTES PARA ValidateResolution ---

func TestValidateResolution(t *testing.T) {
	// Criamos um incidente padrão para usar nos cenários abaixo
	validIncident := &contracts.Incident{
		TargetClass: "security",
		Solution:    "1-3-2",
	}

	// Estrutura de tabela para mapear cenários de teste
	tests := []struct {
		name           string
		activeIncident *contracts.Incident
		playerClass    string
		payload        string
		expectedValid  bool
		expectedMsg    string
	}{
		{
			name:           "Erro quando não há incidente ativo",
			activeIncident: nil,
			playerClass:    "security",
			payload:        "1-3-2",
			expectedValid:  false,
			expectedMsg:    "Não há incidentes críticos no sistema.",
		},
		{
			name:           "Erro quando a classe do jogador não tem permissão",
			activeIncident: validIncident,
			playerClass:    "devops", // Classe diferente de "security"
			payload:        "1-3-2",
			expectedValid:  false,
			expectedMsg:    "Ação negada: Você não tem permissão para esta camada da infra!",
		},
		{
			name:           "Erro quando a solução (payload) está errada",
			activeIncident: validIncident,
			playerClass:    "security",
			payload:        "ERRADO",
			expectedValid:  false,
			expectedMsg:    "Solução inválida: O erro persiste no log.",
		},
		{
			name:           "Sucesso quando todas as condições são atendidas",
			activeIncident: validIncident,
			playerClass:    "security",
			payload:        "1-3-2",
			expectedValid:  true,
			expectedMsg:    "Ação validada! Contribuindo para a resolução...",
		},
	}

	// Executa cada linha da tabela como um sub-teste isolado
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isValid, msg := ValidateResolution(tt.activeIncident, tt.playerClass, tt.payload)

			assert.Equal(t, tt.expectedValid, isValid)
			assert.Equal(t, tt.expectedMsg, msg)
		})
	}
}
