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

// Função para forçar a Go Engine a devolver o estado atual
function triggerGoEngineUpdate(room_id) {
    pub.publish(`room:${room_id}:attacks`, JSON.stringify({
        type: 'ping',
        class: 'system',
        nickname: 'gateway',
        timestamp: Date.now()
    }));
}

module.exports = (io, socket) => {
    
    // Join Game com Validação de Vagas e Salas (RF02 & RF09)
    socket.on('join_game', async (data) => {
        const { nickname, class: playerClass, room_id = 'sala_padrao', boss_id = '' } = data;

        if (!playerClass || !nickname || !room_id) {
            return socket.emit('game_error', { message: "Dados de login incompletos." });
        }

        try {
            // ========================================================
            // BLOQUEIO DE SALAS FANTASMAS (Somente entra onde o Admin abriu)
            // ========================================================
            if (playerClass === 'admin') {
                // Se é o Painel Web (Professor), nós REGISTRAMOS a sala no Redis
                await pub.sadd('active_dungeon_rooms', room_id);
            } else {
                // Se é o Mobile (Aluno), nós VERIFICAMOS se a sala existe no Redis
                const roomExists = await pub.sismember('active_dungeon_rooms', room_id);
                
                if (!roomExists) {
                    // Chuta o usuário antes mesmo de incomodar a Go Engine!
                    return socket.emit('game_error', { 
                        message: `A sala "${room_id}" não existe ou não foi aberta pelo Professor.` 
                    });
                }
            }
            // ========================================================

            // A chave do Redis isolada por sala!
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

            // Inscreve o socket do Socket.io em uma "sala" (room)
            socket.join(room_id);

            // Notifica a Engine em Go
            pub.publish(`room:${room_id}:attacks`, JSON.stringify({
                type: 'join',
                class: playerClass,
                nickname: nickname,
                payload: boss_id, 
                timestamp: Date.now()
            }));

            socket.emit('joined', { status: 'success', nickname, playerClass, room_id });
            console.log(`🎮 ${nickname} entrou como ${playerClass} na sala [${room_id}]`);

            // Chama a função para atualizar a tela do Dashboard instantaneamente
            triggerGoEngineUpdate(room_id);

        } catch (err) {
            console.error("Erro ao processar join_game:", err);
            socket.emit('game_error', { message: "Erro interno ao validar vaga." });
        }
    });

    // Comando de Ataque (RF04)
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

    // Resolução de Incidentes (RF05)
    socket.on('resolve_incident', (data) => {
        const { playerClass, nickname, room_id } = socket.data;
        if (!playerClass || !room_id) return;

        // O 'data.payload' é a solução do desafio que o Mobile envia
        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'resolve',
            class: playerClass,
            nickname: nickname,
            payload: data ? data.payload : "", 
            timestamp: Date.now()
        }));
    });

    // Reset da Sala pelo Admin
    socket.on('admin_reset_room', (data) => {
        const { room_id } = data;
        const { playerClass } = socket.data;

        // Validação de segurança
        if (playerClass !== 'admin') return;

        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'reset',
            class: 'admin',
            nickname: 'GAME_MASTER',
            timestamp: Date.now()
        }));

        io.to(room_id).emit('room_reset', { message: "O professor resetou a sala!" });
        console.log(`🔄 O Mestre resetou a sala [${room_id}].`);
    });

    // Tratamento de Desconexão (RF08)
    socket.on('disconnect', async () => {
        const { playerClass, nickname, room_id } = socket.data;
        
        if (playerClass && nickname && room_id) {
            // Libera a vaga na sala correta do Redis
            await pub.srem(`room:${room_id}:class_members:${playerClass}`, nickname);
            console.log(`❌ ${nickname} saiu da sala [${room_id}].`);
            
            // Chama a função para atualizar a tela do Dashboard com a pessoa a menos
            triggerGoEngineUpdate(room_id);
        }
        
        lastAttack.delete(socket.id);
    });
};