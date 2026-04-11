package boss

import (
	"time"

	models "github.com/wendelvx/pi-devops.git/internal/domain/contracts"
)

func GetBossProfile(id string) models.BossProfile {
	profiles := map[string]models.BossProfile{

		"arquiteto": {
			ID:           "arquiteto",
			Name:         "O Arquiteto (Backend)",
			BaseHP:       2500000,
			BaseDamage:   60,
			AttackSpeed:  4 * time.Second,
			IncidentBias: "database_lock",
		},

		"pixel_perfect": {
			ID:           "pixel_perfect",
			Name:         "O Pixel Perfect (Frontend)",
			BaseHP:       1200000,
			BaseDamage:   15,
			AttackSpeed:  1200 * time.Millisecond,
			IncidentBias: "error_500",
		},

		"tcc_titan": {
			ID:           "tcc_titan",
			Name:         "The TCC Titan",
			BaseHP:       2000000,
			BaseDamage:   35,
			AttackSpeed:  2500 * time.Millisecond,
			IncidentBias: "code_review",
		},

		"guardiao": {
			ID:           "guardiao",
			Name:         "O Guardião (Security/QA)",
			BaseHP:       1500000,
			BaseDamage:   40,
			AttackSpeed:  3 * time.Second,
			IncidentBias: "phishing",
		},
	}

	if profile, ok := profiles[id]; ok {
		return profile
	}
	return profiles["tcc_titan"]
}
