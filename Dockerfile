# Imagem de desenvolvimento com o toolchain completo do Go.
# Pense nela como o equivalente à imagem "node": traz compilador, gofmt, go test etc.
#
# Usamos a imagem baseada em Debian (e não Alpine) porque o race detector
# (capítulo 12) precisa de cgo/glibc. Imagens mínimas para produção ficam
# para o capítulo 22.
FROM golang:1.27

WORKDIR /app

# Copia primeiro só os arquivos de dependências (como package.json + package-lock.json)
# para aproveitar o cache de camadas: se só o código mudar, o download não se repete.
COPY go.mod go.sum* ./
RUN go mod download

COPY . .

# Capítulo executado por padrão. Exemplo de uso sem o compose:
#   docker build -t learn-go --build-arg CAPITULO=01-ola-go .
#   docker run --rm learn-go
ARG CAPITULO=00-ambiente
ENV CAPITULO=${CAPITULO}

CMD ["sh", "-c", "go run ./capitulos/${CAPITULO}"]
