const { pub } = require('../config/redis');

module.exports = (io, socket) => {
  socket.on('attack', (data) => {
    if (!data.class || !data.nickname) return;

    const payload = JSON.stringify({
      type: 'attack',
      class: data.class,
      nickname: data.nickname,
      timestamp: Date.now()
    });

    pub.publish('player_attacks', payload);
  });

  socket.on('disconnect', () => {
    console.log(`❌ Usuário desconectado: ${socket.id}`);
  });
};