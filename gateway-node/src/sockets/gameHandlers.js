const { pub } = require('../config/redis');
const { attackCounter } = require('../config/metrics');

module.exports = (io, socket) => {
    socket.on('join_game', (data) => {
        socket.data.playerClass = data.class;
        socket.data.nickname = data.nickname;
        
        pub.publish('player_attacks', JSON.stringify({
            type: 'join',
            class: data.class,
            nickname: data.nickname,
            timestamp: Date.now()
        }));
    });

    socket.on('attack', () => {
        const { playerClass } = socket.data;

        if (playerClass) {
            // ✅ Incrementa no Prometheus com a label da classe
            attackCounter.inc({ class: playerClass });
        }

        pub.publish('player_attacks', JSON.stringify({
            type: 'attack',
            class: playerClass || 'unknown',
            nickname: socket.data.nickname || 'anon',
            timestamp: Date.now()
        }));
    });
};