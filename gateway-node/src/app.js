const express = require('express');
const http = require('http');
const { Server } = require('socket.io');
const { register } = require('./config/metrics');
const { pub, sub } = require('./config/redis'); // Importado para o Graceful Shutdown
const registerGameHandlers = require('./sockets/gameHandlers');
const initRedisSubscriber = require('./services/redisSubscriber');

const app = express();
const server = http.createServer(app);
const io = new Server(server, { cors: { origin: "*" } });

app.get('/metrics', async (req, res) => {
    try {
        res.set('Content-Type', register.contentType);
        res.end(await register.metrics());
    } catch (ex) {
        res.status(500).send(ex);
    }
});

app.get('/api/health', (req, res) => {
    res.json({ status: "healthy", service: "node-gateway" });
});

initRedisSubscriber(io);

io.on('connection', (socket) => {
    // MUDANÇA AQUI: Reflete que é apenas uma conexão "crua", antes de validar a sala.
    console.log(`🔌 Nova conexão de socket estabelecida: ${socket.id}`);
    registerGameHandlers(io, socket);
});

const PORT = process.env.PORT || 3001;

server.listen(PORT, "0.0.0.0", () => {
    console.log(`🚀 Gateway Metrics & Sockets online na porta ${PORT}`);
});

const shutdown = (signal) => {
    console.log(`\nRecebido ${signal}. Encerrando Dungeon Master Gateway...`);
    
    server.close(() => {
        console.log('Servidor HTTP fechado.');
        
        pub.quit();
        sub.quit();
        console.log('Conexões Redis encerradas.');
        
        process.exit(0);
    });

    setTimeout(() => {
        console.error('Forçando encerramento por timeout.');
        process.exit(1);
    }, 5000);
};

process.on('SIGTERM', () => shutdown('SIGTERM'));
process.on('SIGINT', () => shutdown('SIGINT'));