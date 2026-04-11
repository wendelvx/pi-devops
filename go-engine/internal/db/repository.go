package db

import (
	"log"
	"time"

	models "github.com/wendelvx/pi-devops.git/internal/domain/contracts"
	"gorm.io/gorm"
)

var db *gorm.DB

type Battle struct {
	ID       uint `gorm:"primaryKey"`
	BossID   string
	Result   string
	Duration int
	Rankings []Ranking `gorm:"foreignKey:BattleID"`
}

type Ranking struct {
	ID              uint `gorm:"primaryKey"`
	BattleID        uint
	Nickname        string
	Class           string
	TotalDamage     int
	IncidentsSolved int
}

// -------------------------------------

func SaveBattleResult(status string, bossID string, startTime time.Time, players map[string]*models.PlayerStats) {

	db = ConectDb()

	playerSnapshot := make([]models.PlayerStats, 0, len(players))
	for _, p := range players {
		playerSnapshot = append(playerSnapshot, *p)
	}

	go func(stats []models.PlayerStats) {
		duration := int(time.Since(startTime).Seconds())

		var rankings []Ranking
		for _, p := range stats {
			rankings = append(rankings, Ranking{
				Nickname:        p.Nickname,
				Class:           p.Class,
				TotalDamage:     p.TotalDamage,
				IncidentsSolved: p.IncidentsSolved,
			})
		}

		battle := Battle{
			BossID:   bossID,
			Result:   status,
			Duration: duration,
			Rankings: rankings,
		}

		result := db.Create(&battle)

		if result.Error != nil {
			log.Printf("Erro crítico ao salvar batalha: %v", result.Error)
			return
		}

		log.Printf(" Batalha #%d (%s) salva com %d jogadores!", battle.ID, status, len(rankings))
	}(playerSnapshot)
}
