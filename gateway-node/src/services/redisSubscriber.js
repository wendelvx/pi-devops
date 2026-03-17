const { sub } = require('../config/redis');

module.exports = (io) => {
  sub.subscribe('boss_updates');

  sub.on('message', (channel, message) => {
    if (channel === 'boss_updates') {
      try {
        const gameState = JSON.parse(message);
        io.emit('boss_update', gameState);
      } catch (err) {
        console.error("Erro ao processar boss_update:", err);
      }
    }
  });
};