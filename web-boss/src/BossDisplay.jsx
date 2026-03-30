import { useEffect, useState } from 'react'
import { io } from 'socket.io-client'
import './BossDisplay.css'

import threatImg from './assets/threat.png'
import damagedImg from './assets/damaged.png'
import defeatedImg from './assets/defeated.png'

const socket = io('http://localhost:3000', { autoConnect: true })

function BossDisplay() {
  const [hp, setHp] = useState(100)
  const [bossState, setBossState] = useState('threat')

  const simulateDamage = () => {
    const newHp = Math.max(0, hp - 5)
    setHp(newHp)
    if (newHp === 0) {
      setBossState('defeated')
      return
    }
    setBossState('damaged')
    setTimeout(() => {
      setBossState('threat')
    }, 350)
  }
  const simulateBossVictory = () => {
    setBossState("victory")
  }

  useEffect(() => {
    socket.on("bossUpdate", (data) => {
      const { hp, type } = data
      if (hp !== undefined) {
        setHp(hp)
      }
      if (type === "damage") {
        setBossState("damaged")
        setTimeout(() => {
          setBossState("threat")
        }, 350)
      }
      if (type === "defeated") {
        setBossState("defeated")
      }
      if (type === "bossVictory") {
        setBossState("victory")
      }
    })

    return () => {
      socket.off('connect')
      socket.off('disconnect')
      socket.off('bossUpdate')
    }
  }, [])

  let bossImage = threatImg

  if (bossState === 'damaged') bossImage = damagedImg
  if (bossState === 'defeated') bossImage = defeatedImg

  return (
    <div className="display-container">

      {/* Simulação de dano para gerar animação */}
      <button className="dev-test-button" onClick={simulateDamage}>
        Simular Dano
      </button>
      {/* Simulação de vitória do boss para gerar animação */}
      <button className="dev-test-button victory" onClick={simulateBossVictory}>
        Simular Vitória
      </button>

      {bossState === 'defeated' &&
        <div className="victory-banner">
          <div className="victory-title">
            <span className="line"></span>
            BOSS DEFEATED
            <span className="line"></span>
          </div>

          <div className="victory-subtitle">
            The Boss Has Fallen
          </div>
        </div>
      }

      {bossState === 'victory' &&
        <div className="boss-victory-banner">
          <div className="victory-title">
            <span className="line"></span>
            BOSS WINS
            <span className="line"></span>
          </div>

          <div className="victory-subtitle">
            The Raid Has Been Destroyed
          </div>
        </div>
      }

      <div className="boss-area">
        <div className={`boss-sprite ${bossState === 'damaged' ? 'shake' : ''}`}>
          <img src={bossImage} alt="Boss"/>
        </div>
      </div>

      <div className="ui-footer">
        <div className="hp-bar-container">
          <div
            className={`hp-bar-fill ${hp < 30 ? 'critical' : ''}`}
            style={{ width: `${hp}%` }}
          />
          <span className="hp-label">
            BOSS HP: {hp}%
          </span>
        </div>
      </div>
    </div>
  )
}

export default BossDisplay