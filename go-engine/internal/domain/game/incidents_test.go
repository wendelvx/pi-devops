package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateIncident_ReturnsValidIncident(t *testing.T) {
	incident := GenerateIncident()

	assert.NotEmpty(t, incident.ID, "O ID do incidente não deveria estar vazio")
	assert.NotEmpty(t, incident.Title, "O Título não deveria estar vazio")
	assert.NotEmpty(t, incident.Description, "A Descrição não deveria estar vazia")
	assert.NotEmpty(t, incident.TargetClass, "A Classe Alvo não deveria estar vazia")
	assert.NotEmpty(t, incident.Type, "O Tipo do incidente não deveria estar vazio")
	assert.NotEmpty(t, incident.Solution, "A Solução não deveria estar vazia")

	assert.Greater(t, incident.Points, 0, "Os pontos devem ser maiores que zero")
	assert.Greater(t, incident.Duration, 0, "A duração deve ser maior que zero")
	assert.Greater(t, incident.RequiredResolutions, 0, "As resoluções requeridas devem ser maiores que zero")

	validClasses := []string{"security", "devops", "back-end", "front-end", "qa"}
	assert.Contains(t, validClasses, incident.TargetClass, "A classe alvo deve ser uma das classes mapeadas no sistema")
}

func TestGenerateIncident_Randomness(t *testing.T) {
	firstIncident := GenerateIncident()
	isDifferent := false

	for i := 0; i < 20; i++ {
		newIncident := GenerateIncident()
		if newIncident.ID != firstIncident.ID {
			isDifferent = true
			break
		}
	}

	assert.True(t, isDifferent, "A função deveria gerar incidentes diferentes ao longo de múltiplas chamadas")
}
