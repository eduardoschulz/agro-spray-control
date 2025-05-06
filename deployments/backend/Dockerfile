#peguei a ultima versao, nao sei se eh a melhor para esse caso
FROM golang:1.24.2-bookworm as base

WORKDIR /build
COPY go.mod go.sum ./ 

RUN go mod download
COPY . .
RUN go build -o go-backend

#nao sei se essa eh a melhor porta para isso
EXPOSE 8000
CMD ["/build/go-backend"]


