import { useEffect, useState } from 'react';
import { useParams } from 'react-router-dom';
import { io } from 'socket.io-client';

// O seu túnel Cloudflare
const SOCKET_URL = "https://sand-folder-reserve-circuits.trycloudflare.com"; 
const socket = io(SOCKET_URL, { autoConnect: false });

export default function TelaoScreen() {
  const { roomCode } = useParams();
  const [gameState, setGameState] = useState(null);

  useEffect(() => {
    socket.connect();
    
    const onConnect = () => {
      console.log("✅ Telão conectado ao Gateway!");
      socket.emit('join_game', { 
        type: 'join', 
        nickname: 'TELÃO', 
        class: 'spectator', 
        room_id: roomCode 
      });
    };

    const onBossUpdate = (data) => {
      setGameState(data);
    };

    socket.on('connect', onConnect);
    socket.on('boss_update', onBossUpdate);

    return () => {
      socket.off('connect', onConnect);
      socket.off('boss_update', onBossUpdate);
      socket.disconnect();
    };
  }, [roomCode]);

  if (!gameState) {
    return (
      <div className="min-h-screen bg-neutral-950 flex items-center justify-center text-red-500 font-mono text-xl animate-pulse">
        Conectando à Masmorra {roomCode}...
      </div>
    );
  }

  const currentBoss = gameState?.current_boss || {};
  const maxHp = currentBoss?.max_hp || 1000;
  const hpPercent = Math.max(0, ((gameState?.boss_hp || 0) / maxHp) * 100);
  const isCaos = gameState?.active_incident != null;

  // A MÁGICA DE DEFESA AQUI: translate="no" impede que o Google Chrome quebre o React
  return (
    <div translate="no" className={`min-h-screen flex flex-col font-sans transition-colors duration-700 ${isCaos ? 'bg-red-950' : 'bg-neutral-950'} text-white relative overflow-hidden`}>
      
      <header className="flex justify-between items-center p-6 border-b border-neutral-800/50 bg-black/20 z-10">
        <div className="bg-neutral-900 border border-neutral-700 px-6 py-2 rounded-full font-mono text-2xl tracking-widest shadow-lg">
          <span className="text-neutral-500">SALA:</span> {roomCode}
        </div>
        {isCaos && (
          <div className="bg-red-600 text-white px-8 py-2 rounded-full font-black animate-pulse shadow-[0_0_20px_rgba(220,38,38,0.8)] tracking-wider">
            ⚠️ INCIDENTE ATIVO
          </div>
        )}
      </header>

      <main className="flex-1 flex flex-col items-center justify-center p-8 z-10">
        <div className="flex flex-col items-center gap-8 w-full max-w-4xl">
          
          {/* A MÁGICA DA IMAGEM AQUI: Avatar novo e confiável */}
          <img 
            src={currentBoss?.avatar_url || 'https://ui-avatars.com/api/?name=B+O+S+S&background=111111&color=ff0033&size=300'} 
            alt="Boss Avatar" 
            className={`w-64 h-64 rounded-full border-4 object-cover transition-all duration-300 ${isCaos ? 'border-red-500 shadow-[0_0_50px_rgba(220,38,38,0.6)] animate-bounce' : 'border-neutral-800 shadow-2xl'}`}
          />
          
          <div className="text-center">
            <h2 className="text-6xl font-black tracking-tighter mb-2">{currentBoss?.name}</h2>
            <p className="text-neutral-500 uppercase tracking-widest text-lg font-bold">{currentBoss?.class}</p>
          </div>
          
          <div className="w-full space-y-4">
            <div className="w-full h-12 bg-neutral-900 rounded-full overflow-hidden border-2 border-neutral-800 relative shadow-inner">
              <div 
                className="h-full bg-gradient-to-r from-red-700 to-red-500 transition-all duration-300 ease-out"
                style={{ width: `${hpPercent}%` }}
              ></div>
            </div>
            <div className="flex justify-between items-center px-4">
              <span className="text-4xl font-mono text-red-500 font-black tracking-tighter">
                {gameState?.boss_hp} <span className="text-neutral-600 text-2xl">/ {maxHp} HP</span>
              </span>
              <span className="text-neutral-500 font-mono font-bold">{hpPercent.toFixed(1)}%</span>
            </div>
          </div>
        </div>
      </main>

      {isCaos && (
        <div className="absolute inset-0 bg-black/80 backdrop-blur-md flex flex-col items-center justify-center p-8 text-center z-50">
          <h3 className="text-7xl font-black text-red-500 mb-6 uppercase tracking-tighter drop-shadow-[0_5px_5px_rgba(220,38,38,0.5)]">
            {gameState?.active_incident?.title}
          </h3>
          <p className="text-3xl text-neutral-300 mb-12 max-w-4xl leading-relaxed">
            {gameState?.active_incident?.description}
          </p>
          
          <div className="bg-red-900/40 border-2 border-red-500 text-red-200 px-10 py-4 rounded-full text-3xl font-black mb-12 uppercase tracking-widest">
            AÇÃO REQUERIDA: <span className="text-white">{gameState?.active_incident?.target_class}</span>
          </div>
          
          <div className="text-[12rem] leading-none font-mono font-black text-red-500 animate-pulse drop-shadow-[0_0_30px_rgba(220,38,38,0.8)]">
            {gameState?.incident_timer}s
          </div>
        </div>
      )}

      <footer className="p-6 border-t border-neutral-800/50 bg-black/40 z-10">
        <p className="font-mono text-2xl text-green-400 text-center font-bold tracking-wide">
          {gameState?.last_action}
        </p>
      </footer>

    </div>
  );
}