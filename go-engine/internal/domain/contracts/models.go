package models

import "time"

type BossProfile struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	BaseHP       float64       `json:"base_hp"`
	BaseDamage   float64       `json:"base_damage"`
	AttackSpeed  time.Duration `json:"attack_speed"`
	IncidentBias string        `json:"incident_bias"`
}

type PlayerStats struct {
	Nickname        string `json:"nickname"`
	Class           string `json:"class"`
	TotalDamage     int    `json:"total_damage"`
	IncidentsSolved int    `json:"incidents_solved"`
}

type IncidentMeta struct {
	ID            string `json:"id"`
	TargetQuota   int    `json:"target_quota"`
	CurrentClicks int    `json:"current_clicks"`
	RequiredClass string `json:"required_class"`
}
