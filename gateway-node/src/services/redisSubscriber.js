const { pub, sub } = require('../config/redis');

module.exports = (io) => {
    // 1. Array fixo com as classes do jogo para contagem
    const CLASSES = ['front-end', 'back-end', 'devops', 'qa', 'security'];

    // Em vez de 'subscribe', usamos 'psubscribe' para ouvir um padrão de canais.
    // O asterisco (*) funciona como um wildcard para qualquer ID de sala.
    sub.psubscribe('room:*:boss_updates', (err, count) => {
        if (err) {
            console.error("❌ Erro ao assinar padrão room:*:boss_updates no Redis:", err);
            return;
        }
        console.log(`📡 Gateway: Ouvindo atualizações de Boss em ${count} padrão(ões) de sala.`);
    });

    // Evento 'pmessage' é acionado quando uma mensagem chega em um canal assinado por padrão
    sub.on('pmessage', async (pattern, channel, message) => {
        try {
            // O canal chega no formato: "room:sala_padrao:boss_updates"
            // Vamos extrair o ID da sala separando a string pelo ":"
            const parts = channel.split(':');
            
            // Validação de segurança para garantir o formato correto
            if (parts.length >= 3 && parts[0] === 'room' && parts[2] === 'boss_updates') {
                const room_id = parts[1];
                let gameState = JSON.parse(message);

                // ========================================================
                // 2. NOVO: Conta a quantidade de players de cada classe no Redis
                // ========================================================
                const classCounts = {};
                
                // Promise.all executa as buscas no Redis de forma paralela e rápida
                await Promise.all(CLASSES.map(async (className) => {
                    const count = await pub.scard(`room:${room_id}:class_members:${className}`);
                    classCounts[className] = count;
                }));

                // Injeta as contagens dentro do objeto do estado do jogo
                gameState.class_counts = classCounts;
                // ========================================================

                // Envia a atualização APENAS para os sockets que deram 'join' nesta sala específica
                io.to(room_id).emit('boss_update', gameState);
            }
        } catch (err) {
            console.error(`Erro ao processar boss_update no canal ${channel}:`, err);
        }
    });
};