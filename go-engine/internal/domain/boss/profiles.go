package boss

import "pi-devops/internal/domain/contracts"

// GetProfessorBoss retorna o perfil do professor baseado na disciplina/ID
func GetProfessorBoss(id string) contracts.Boss {
	profiles := map[string]contracts.Boss{
		"infra_boss": {
			ID:          "infra_boss",
			Name:        "Prof. de DevOps", // Substitua pelo nome real
			Class:       "DevOps",
			MaxHP:       20000,
			Weakness:    "devops",
			Description: "O mestre dos containers. Seus pods são indestrutíveis!",
		},
		"logic_boss": {
			ID:          "logic_boss",
			Name:        "Prof. de Back-end",
			Class:       "Back-end",
			MaxHP:       15000,
			Weakness:    "back-end",
			Description: "Especialista em algoritmos. Cuidado com o Stack Overflow!",
		},
		"security_boss": {
			ID:          "security_boss",
			Name:        "Prof. de Security",
			Class:       "Security",
			MaxHP:       18000,
			Weakness:    "security",
			Description: "O Guardião do Firewall. Suas vulnerabilidades serão expostas!",
		},
	}

	if boss, ok := profiles[id]; ok {
		return boss
	}
	
	// Default caso não encontre
	return profiles["infra_boss"]
}