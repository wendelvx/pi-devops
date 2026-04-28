import { useState } from 'react';
import { useNavigate } from 'react-router-dom';

export default function SetupScreen() {
  const [roomCode, setRoomCode] = useState('');
  const [bossName, setBossName] = useState('Prof. X');
  const navigate = useNavigate();

  const handleStart = () => {
    if (!roomCode) return alert('Digite o código da sala!');
    navigate(`/arena/${roomCode.toUpperCase()}`);
  };

  return (
    <div className="min-h-screen bg-neutral-950 flex items-center justify-center p-4 text-white font-sans">
      <div className="bg-neutral-900 border border-red-500/30 shadow-[0_0_30px_rgba(220,38,38,0.15)] rounded-2xl p-8 w-full max-w-md flex flex-col gap-8">
        
        <div className="text-center space-y-2">
          <h1 className="text-3xl font-black text-red-500 tracking-wider">DUNGEON MASTER</h1>
          <p className="text-neutral-400 uppercase tracking-widest text-xs font-bold">Painel do Professor</p>
        </div>
        
        <div className="space-y-4">
          <div>
            <label className="block text-xs font-bold text-neutral-500 mb-2 uppercase tracking-wider">Código da Sala</label>
            <input 
              type="text" 
              placeholder="Ex: ADS2026" 
              value={roomCode}
              onChange={(e) => setRoomCode(e.target.value)}
              className="w-full bg-neutral-950 border border-neutral-800 focus:border-red-500 focus:ring-1 focus:ring-red-500 rounded-xl px-4 py-4 text-white outline-none transition-all uppercase placeholder:text-neutral-700 font-mono text-lg"
            />
          </div>
          
          <div>
            <label className="block text-xs font-bold text-neutral-500 mb-2 uppercase tracking-wider">Selecione o Boss</label>
            <select 
              value={bossName} 
              onChange={(e) => setBossName(e.target.value)}
              className="w-full bg-neutral-950 border border-neutral-800 focus:border-red-500 focus:ring-1 focus:ring-red-500 rounded-xl px-4 py-4 text-white outline-none transition-all appearance-none cursor-pointer"
            >
              <option value="Prof. X">Professor de Infraestrutura</option>
              <option value="Prof. Y">Professor de Backend</option>
            </select>
          </div>
        </div>

        <button 
          onClick={handleStart}
          className="w-full bg-red-600 hover:bg-red-500 text-white font-black py-4 rounded-xl transition-all shadow-[0_0_15px_rgba(220,38,38,0.4)] hover:shadow-[0_0_25px_rgba(220,38,38,0.6)] uppercase tracking-widest active:scale-95 mt-2"
        >
          Iniciar Telão
        </button>
      </div>
    </div>
  );
}