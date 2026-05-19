package contracts

// Boss representa a entidade de um professor no jogo (RF10)
type Boss struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Class       string `json:"class"`       // Disciplina do professor
    MaxHP       int    `json:"max_hp"`
    Description string `json:"description"` // Bio pedagógica do prof
    Weakness    string `json:"weakness"`    // Classe técnica que tem vantagem
    AvatarURL   string `json:"avatar_url"`  // Foto do Professor
}

// Incident representa um evento de caos específico para uma classe (RF05)
type Incident struct {
    ID                  string `json:"id"`
    Title               string `json:"title"`       // Ex: "Ataque DDoS", "Bug em Produção"
    Description         string `json:"description"` // O que o aluno deve fazer
    TargetClass         string `json:"target_class"` // Qual classe pode resolver
    Type                string `json:"type"`         // RAPID_CLICK, SEQUENCE, PUZZLE
    Solution            string `json:"solution"`     // O valor esperado para validar a cura
    Points              int    `json:"points"`       // Dano causado ao Boss ao resolver
    Duration            int    `json:"duration"`     // Quantos segundos o time tem para resolver
    RequiredResolutions int    `json:"required_resolutions"` // NOVO: Quantas vezes precisa ser resolvido
    CurrentResolutions  int    `json:"current_resolutions"`  // NOVO: Quantas vezes já foi resolvido
}

// Estrutura para calcular quem bateu mais no Boss
type PlayerStat struct {
    Nickname string `json:"nickname"`
    Class    string `json:"class"`
    Damage   int    `json:"damage"`
}

// PlayerAction estrutura as mensagens vindas do Redis (Node -> Go)
type PlayerAction struct {
    Type      string `json:"type"`      // join, attack, resolve
    Class     string `json:"class"`     // Classe do aluno
    Nickname  string `json:"nickname"`  // Nickname do aluno
    Payload   string `json:"payload"`   // A resposta do desafio
    Timestamp int64  `json:"timestamp"`
}

// GameState representa o estado atual da arena enviado ao Mobile (RF03)
type GameState struct {
    BossHP         int         `json:"boss_hp"`
    TeamHP         int         `json:"team_hp"`         // Vida da Equipe
    MaxTeamHP      int         `json:"max_team_hp"`
    Status         string      `json:"status"`          // fighting, victory, defeat
    LastAction     string      `json:"last_action"`     // Log rápido para o feed (RF11)
    ActiveIncident *Incident   `json:"active_incident"` // Incidente atual (pode ser null)
    IncidentTimer  int         `json:"incident_timer"`  // Cronômetro regressivo
    CurrentBoss    Boss        `json:"current_boss"`    // Dados do professor atual
    TopRank        []PlayerStat `json:"top_rank"`            // Guarda o destaque da rodada
}