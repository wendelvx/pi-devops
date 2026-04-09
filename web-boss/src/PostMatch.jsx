import { useEffect, useState } from "react"
import "./PostMatch.css"

function PostMatch() {
  const [data, setData] = useState(null)
  const [error, setError] = useState(false)

  useEffect(() => {
    fetch("http://localhost:3000/match-result")
      .then(res => {
        if (!res.ok) throw new Error("Erro na API")
        return res.json()
      })
      .then(setData)
      .catch(() => setError(true))
  }, [])

  if (error) return <div className="error">Erro ao carregar dados</div>
  if (!data) return <div className="loading">Loading...</div>

  const sortedPlayers = [...data.players].sort((a, b) => b.score - a.score)

  return (
    <div className="postmatch-container">

      <div className="mvp-card">
        <h2>🏆 MVP DA RODADA</h2>
        <p className="mvp-name">{data.mvp.nickname}</p>
        <p className="mvp-class">{data.mvp.class}</p>
        <p className="mvp-score">Score: {data.mvp.score}</p>
      </div>

      <div className="ranking">
        {sortedPlayers.map((p, i) => (
          <div
            key={p.nickname}
            className={`player-row ${i === 0 ? "first-place" : ""}`}
          >
            <div className="player-header">
              #{i + 1} - {p.nickname} ({p.class})
            </div>

            <div className="player-stats">
              Damage: {p.damage} | Incidents: {p.incidentsResolved} | Score: {p.score}
            </div>
          </div>
        ))}
      </div>

    </div>
  )
}

export default PostMatch