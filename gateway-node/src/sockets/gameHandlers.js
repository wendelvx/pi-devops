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
            // BLOQUEIO DE SALAS FANTASMAS
            if (playerClass === 'admin') {
                await pub.sadd('active_dungeon_rooms', room_id);
            } else {
                const roomExists = await pub.sismember('active_dungeon_rooms', room_id);
                
                if (!roomExists) {
                    return socket.emit('game_error', { 
                        message: `A sala "${room_id}" não existe ou não foi aberta pelo Professor.` 
                    });
                }
            }

            const roomClassKey = `room:${room_id}:class_members:${playerClass}`;
            const currentCount = await pub.scard(roomClassKey);

            if (CLASS_LIMITS[playerClass] !== undefined && currentCount >= CLASS_LIMITS[playerClass]) {
                return socket.emit('game_error', { 
                    message: `A classe ${playerClass} atingiu o limite na sala ${room_id}.` 
                });
            }

            // 💥 CORREÇÃO 1: Remover o aluno do Redis da sala ANTERIOR (se ele estiver trocando)
            if (socket.data.room_id && socket.data.playerClass && socket.data.nickname) {
                await pub.srem(`room:${socket.data.room_id}:class_members:${socket.data.playerClass}`, socket.data.nickname);
                
                // Dispara atualização para a engine recalcular a sala antiga
                triggerGoEngineUpdate(socket.data.room_id);
            }

            await pub.sadd(roomClassKey, nickname);

            // 💥 CORREÇÃO CRÍTICA: Desconecta das salas antigas no socket.io para evitar "Room Bleeding"
            socket.rooms.forEach(room => {
                if (room !== socket.id) socket.leave(room);
            });

            socket.data.room_id = room_id;
            socket.data.playerClass = playerClass;
            socket.data.nickname = nickname;

            socket.join(room_id);

            console.log(`🎮 ${nickname} entrou como ${playerClass} na sala [${room_id}]`);

            pub.publish(`room:${room_id}:attacks`, JSON.stringify({
                type: 'join',
                class: playerClass,
                nickname: nickname,
                payload: boss_id, 
                timestamp: Date.now()
            }));

            socket.emit('joined', { status: 'success', nickname, playerClass, room_id });

        } catch (err) {
            console.error("Erro ao processar join_game:", err);
            socket.emit('game_error', { message: "Erro interno ao validar vaga." });
        }
    });

    // ==========================================
    // CONEXÃO DE ESPECTADOR DO ADMIN
    // ==========================================
    socket.on('join_admin_spectator', async (data) => {
        const { room_id } = data;
        if (!room_id) return;
        
        // 💥 CORREÇÃO CRÍTICA: Garante que o Admin não continue ouvindo a sala anterior
        socket.rooms.forEach(room => {
            if (room !== socket.id) socket.leave(room);
        });

        await pub.sadd('active_dungeon_rooms', room_id);
        
        socket.join(room_id);
        socket.data.playerClass = 'admin';
        socket.data.room_id = room_id;
        console.log(`👁️ O Mestre voltou a assistir a sala [${room_id}].`);
    });

    socket.on('admin_request_state', (data) => {
        const { room_id } = data;
        if (!room_id) return;
        triggerGoEngineUpdate(room_id);
    });

    // ==========================================
    // START DA PARTIDA (Trava do Mestre)
    // ==========================================
    socket.on('admin_start_battle', (data) => {
        const { room_id } = data;
        if (socket.data.playerClass !== 'admin') return;

        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'start_battle',
            class: 'admin',
            nickname: 'GAME_MASTER',
            timestamp: Date.now()
        }));
        
        console.log(`🔥 O Mestre liberou a batalha na sala [${room_id}]!`);
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

        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'attack',
            class: playerClass,
            nickname: nickname,
            timestamp: Date.now()
        }));
    });

    // Resolução de Incidentes
    socket.on('resolve_incident', (data) => {
        const { playerClass, nickname, room_id } = socket.data;
        if (!playerClass || !room_id) return;

        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'resolve',
            class: playerClass,
            nickname: nickname,
            payload: data ? data.payload : "", 
            timestamp: Date.now()
        }));
    });

    // Reset da Sala
    socket.on('admin_reset_room', (data) => {
        const { room_id } = data;
        
        if (socket.data.playerClass !== 'admin') return;

        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'reset',
            class: 'admin',
            nickname: 'GAME_MASTER',
            timestamp: Date.now()
        }));

        io.to(room_id).emit('room_reset', { message: "O professor resetou a sala!" });
        console.log(`🔄 O Mestre resetou a sala [${room_id}].`);
    });

    // Desconexão
    socket.on('disconnect', async () => {
        const { playerClass, nickname, room_id } = socket.data;
        
        if (playerClass && nickname && room_id && playerClass !== 'admin') {
            await pub.srem(`room:${room_id}:class_members:${playerClass}`, nickname);
            console.log(`❌ ${nickname} saiu da sala [${room_id}].`);
            triggerGoEngineUpdate(room_id);
        }
        
        lastAttack.delete(socket.id);
    });

    // Destruição da Sala
    socket.on('admin_delete_room', async (data) => {
        const { room_id } = data;
        if (!room_id) return;

        // 1. Remove a sala do Redis para que ninguém novo consiga logar
        await pub.srem('active_dungeon_rooms', room_id);
        
        // 💥 CORREÇÃO 2: Expurgo total dos fantasmas no Redis
        const CLASSES = ['front-end', 'back-end', 'devops', 'qa', 'security'];
        await Promise.all(CLASSES.map(c => 
            pub.del(`room:${room_id}:class_members:${c}`)
        ));

        // 2. Notifica o Go Engine para matar o Game Loop
        pub.publish(`room:${room_id}:attacks`, JSON.stringify({
            type: 'delete',
            class: 'admin',
            nickname: 'GAME_MASTER',
            timestamp: Date.now()
        }));

        // 3. Emite um aviso aos alunos conectados e os expulsa da sala do socket
        io.to(room_id).emit('room_deleted', { message: "A instância foi encerrada permanentemente pelo Professor." });
        io.in(room_id).socketsLeave(room_id);
        
        console.log(`🗑️ Sala [${room_id}] deletada e todos os alunos/fantasmas foram expulsos.`);
    });
};