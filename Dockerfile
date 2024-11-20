# Use a imagem oficial do Go como base
FROM golang:1.23.2

# Setando o diretório de trabalho
WORKDIR /app

# Copiar o código fonte para o contêiner
COPY . .

# Baixar as dependências (caso existam)
RUN go mod tidy

# Compilar o aplicativo
#RUN go build -o streesstest .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/streesstest ./cmd/cli/main.go

# Definir o comando para rodar o programa
ENTRYPOINT ["./streesstest"]
