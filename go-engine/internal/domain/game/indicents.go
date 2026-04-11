package game

import (
	"log"

	models "github.com/wendelvx/pi-devops.git/internal/domain/contracts"
)

var AllIncidents = map[string]models.IncidentMeta{
	"error_500":         {ID: "error_500", TargetQuota: 15, RequiredClass: "DevOps"},
	"code_review":       {ID: "code_review", TargetQuota: 20, RequiredClass: "All"},
	"phishing":          {ID: "phishing", TargetQuota: 20, RequiredClass: "Security"},
	"database_lock":     {ID: "database_lock", TargetQuota: 40, RequiredClass: "Security"},
	"legacy_code_spill": {ID: "legacy_code_spill", TargetQuota: 50, RequiredClass: "Back-end"},
}

// SpawnRandomIncident escolhe um incidente aleatório compatível com as classes ativas
func PossibleIncident() ([]models.IncidentMeta, error) {
	State.Mu.Lock()
	defer State.Mu.Unlock()

	if State.ActiveIncident.ID != "" {
		return nil, nil
	}

	activeClasses := make(map[string]bool)
	for className, count := range State.ClassCounts {
		if count > 0 {
			activeClasses[className] = true
		}
	}

	var possibleIncidents []models.IncidentMeta

	for _, incident := range AllIncidents {
		if incident.RequiredClass == "All" || activeClasses[incident.RequiredClass] {
			possibleIncidents = append(possibleIncidents, incident)
		}
	}

	if len(possibleIncidents) == 0 {
		log.Println("Nenhum incidente compatível com as classes atuais.")
		return nil, nil
	}

	return possibleIncidents, nil
}
