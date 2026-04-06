const { pub } = require('../config/redis');
const { attackCounter } = require('../config/metrics');

const lastAttack = new Map();

module.exports = (io, socket) => {
    // 1. Join Game: Registro inicial do aluno
    socket.on('join_game', (data = {}) => {
        socket.data.playerClass = data.class;
        socket.data.nickname = data.nickname;
        
        pub.publish('player_attacks', JSON.stringify({
            type: 'join',
            class: data.class,
            nickname: data.nickname,
            timestamp: Date.now()
        }));
    });

    socket.on('admin_start_game', (data = {}) => {
        pub.publish('player_attacks', JSON.stringify({
            ...data,
            type: 'start_command',
            boss: data.boss,
            timestamp: Date.now()
        }));
    });

    socket.on('attack', () => {
        const now = Date.now();
        const lastTime = lastAttack.get(socket.id) || 0;

        if (now - lastTime < 100) {
            return socket.emit('game_error', { 
                message: "Calma, mestre! Code review em andamento (Clique muito rápido)." 
            });
        }
        
        lastAttack.set(socket.id, now);

        const { playerClass, nickname } = socket.data;

        if (playerClass) {
            attackCounter.inc({ class: playerClass });
        }

        pub.publish('player_attacks', JSON.stringify({
            type: 'attack',
            class: playerClass || 'unknown',
            nickname: nickname || 'anon',
            timestamp: Date.now()
        }));
    });

    socket.on('disconnect', () => {
        lastAttack.delete(socket.id);
        console.log(`❌ Aluno desconectado: ${socket.id}`);
    });
};
