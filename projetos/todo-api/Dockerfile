# Dockerfile de produção, em duas etapas (multi-stage).
#   docker build -t todo-api projetos/todo-api

# Etapa 1: compila o binário com a imagem completa do Go
FROM golang:1.27 AS build
WORKDIR /src

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 gera um binário estático, que roda sem nenhuma biblioteca do sistema
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/todo-api ./cmd/api

# Etapa 2: imagem final mínima, só com o binário (sem Go, sem shell, sem node_modules)
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/todo-api /todo-api
EXPOSE 8080
ENTRYPOINT ["/todo-api"]
