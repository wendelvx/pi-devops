FROM node:18-alpine

WORKDIR /app

# Instala dependências primeiro para melhor cache em alterações de código
COPY package*.json ./
ENV NODE_ENV=production
RUN npm install --production

# Copia o restante do código da aplicação
COPY . .

EXPOSE 3001

CMD ["node", "server.js"]
