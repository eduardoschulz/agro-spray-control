FROM node:lts-alpine

#instala um server http para servir o conteudo estatico
RUN npm install -g http-server

WORKDIR /app

COPY package*.json ./
RUN npm install

COPY . .

RUN npm run build

EXPOSE 8888

CMD ["http-server", "dist"]
