const express = require('express');
const http = require('http');
const { Server } = require('socket.io');
const registerGameHandlers = require('./sockets/gameHandlers');
const initRedisSubscriber = require('./services/redisSubscriber');

const app = express();
const server = http.createServer(app);
const io = new Server(server, { cors: { origin: "*" } });

app.get('/api/health', (req, res) => {
  res.json({ status: "healthy", service: "node-gateway" });
});

initRedisSubscriber(io);

io.on('connection', (socket) => {
  console.log(`📡 Novo aluno conectado: ${socket.id}`);
  registerGameHandlers(io, socket);
});

const PORT = process.env.PORT || 3001;
server.listen(PORT, "0.0.0.0", () => {
  console.log(`Gateway online na porta ${PORT}`);
});