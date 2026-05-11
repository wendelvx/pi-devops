const { pub } = require('../config/redis');
const { attackCounter } = require('../config/metrics');

// Configuração de limites por classe (RF02)
const CLASS_LIMITS = {
    'front-end': 15,
    'back-end': 15,
    'devops': 8,
    'qa': 10,
    'security': 8
};

const lastAttack = new Map();

module.exports = (io, socket) => {
    
    // 1. Join Game com Validação de Vagas e Salas (RF02 & RF09)
    socket.on('join_game', async (data) => {
        // EXTRAINDO O BOSS_ID AQUI!
        const { nickname, class: playerClass, room_id = 'sala_padrao', boss_id = '' } = data;

        if (!playerClass || !nickname) {
            return socket.emit('game_error', { message: "Dados de login incompletos." });
        }

        try {
            // A chave do Redis agora é isolada por sala!
            const roomClassKey = `room:${room_id}:class_members:${playerClass}`;
            
            // Consulta o Redis para saber quantos membros essa classe já possui NESTA SALA
            const currentCount = await pub.scard(roomClassKey);

            // Valida o limite de vagas (ignorando a classe 'admin' que não tem limite)
            if (CLASS_LIMITS[playerClass] !== undefined && currentCount >= CLASS_LIMITS[playerClass]) {
                return socket.emit('game_error', { 
                    message: `A classe ${playerClass} atingiu o limite na sala ${room_id}.` 
                });
            }

            // Adiciona o nickname ao set da classe na sala específica
            await pub.sadd(roomClassKey, nickname);

            // Salva os dados na sessão do socket
            socket.data.room_id = room_id;
            socket.data.playerClass = playerClass;
            socket.data.nickname = nickname;

            // Inscreve o socket do Socket.io em uma "sala" (room) para facilitar o broadcast depois
            socket.join(room_id);

            // Notifica a Engine, publicando no canal específico da sala
            pub.publish(`room:${room_id}:attacks`, JSON.stringify({
                type: 'join',
                class: playerClass,
                nickname: nickname,
                payload: boss_id, // MÁGICA AQUI: O payload avisa o Go qual Boss criar!
                timestamp: Date.now()
            }));

            socket.emit('joined', { status: 'success', nickname, playerClass, room_id });
            console.log(`🎮 ${nickname} entrou como ${playerClass} na sala [${room_id}]`);

        } catch (err) {
            console.error("Erro ao processar join_game:", err);
            socket.emit('game_error', { message: "Erro interno ao validar vaga." });
        }
    });

    // 2. Comando de Ataque (RF04)
    socket.on('attack', () => {
        const now = Date.now();
        const lastTime = lastAttack.get(socket.id) || 0;

        if (now - lastTime < 100) {
            return socket.emit('game_error', { message: "Calma! Respeite o Global Cooldown." });
        }
        
        lastAttack.set(socket.id, now);
        const { playerClass, nickname, room_id } = socket.data;

        if (!playerClass || !room_id) return;

        attackCounter.inc({ class: playerClass });

        // Publica o ataque no canal da sala correspondente
        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'attack',
            class: playerClass,
            nickname: nickname,
            timestamp: Date.now()
        }));
    });

    // 3. Resolução de Incidentes (RF05)
    socket.on('resolve_incident', (data) => {
        const { playerClass, nickname, room_id } = socket.data;
        if (!playerClass || !room_id) return;

        // O 'data.payload' é a solução do desafio (ex: "3-1-2" ou "15") que o Mobile envia
        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'resolve',
            class: playerClass,
            nickname: nickname,
            payload: data ? data.payload : "", 
            timestamp: Date.now()
        }));
    });

    socket.on('admin_reset_room', (data) => {
        const { room_id } = data;
        const { playerClass } = socket.data;

        // Validação de segurança: apenas quem entrou como admin pode resetar
        if (playerClass !== 'admin') return;

        // 1. Avisa a Go-Engine para zerar o estado interno (HP, Incidentes, Status)
        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'reset',
            class: 'admin',
            nickname: 'GAME_MASTER',
            timestamp: Date.now()
        }));

        // 2. Avisa os Apps Mobile conectados nesta sala para exibirem o alerta e saírem da tela de vitória
        io.to(room_id).emit('room_reset', { message: "O professor resetou a sala!" });
        console.log(`🔄 O Mestre resetou a sala [${room_id}].`);
    });

    // 4. Tratamento de Desconexão (RF08)
    socket.on('disconnect', async () => {
        const { playerClass, nickname, room_id } = socket.data;
        
        if (playerClass && nickname && room_id) {
            // Libera a vaga na sala correta do Redis
            await pub.srem(`room:${room_id}:class_members:${playerClass}`, nickname);
            console.log(`❌ ${nickname} saiu da sala [${room_id}].`);
        }
        
        lastAttack.delete(socket.id);
    });
};