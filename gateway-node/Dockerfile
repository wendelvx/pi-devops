FROM node:18-alpine

WORKDIR /app

# Copiamos os arquivos de manifesto primeiro para aproveitar o cache das camadas
COPY package*.json ./

# Instala as dependências (incluindo o express que está faltando)
RUN npm install --omit=dev

# Agora copiamos o resto do código
COPY . .

EXPOSE 3001

# Boa prática: aguardar o Redis se possível, mas o depends_on no compose já ajuda
CMD ["node", "src/app.js"]