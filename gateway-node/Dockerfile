FROM node:18-alpine

WORKDIR /app

# Instala dependências primeiro para melhor cache em alterações de código
COPY package*.json ./
RUN npm install

# Copia o restante do código da aplicação
COPY . .

EXPOSE 3001

CMD ["node", "server.js"]
