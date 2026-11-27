# Aprendendo Go 🐹

Repositório de estudos de **Go (Golang)**, do básico ao avançado, escrito por quem
vem do **JavaScript**. Cada capítulo explica os conceitos fazendo pontes com o JS,
traz código executável e exercícios validados por testes automáticos.

👉 **[Veja o roadmap completo](ROADMAP.md)**

## Requisitos

- [Docker](https://www.docker.com/) (todo o código roda dentro de containers, nada precisa ser instalado na máquina)
- `make`
- *Opcional:* Go instalado localmente (`brew install go`) para o editor ter autocomplete
  e checagem de erros (extensão **Go** no VS Code ou GoLand).

## Como usar

```bash
make help            # lista todos os comandos
make list            # lista os capítulos

make run CAP=00      # roda o main.go do capítulo 00
make ex  CAP=00      # roda os testes dos exercícios do capítulo 00
make test CAP=00     # roda todos os testes do capítulo
make shell           # abre um terminal dentro do container Go

make fmt             # formata o código (gofmt, o "prettier" do Go)
make vet             # análise estática (go vet, tipo um eslint embutido)
```

O código é montado no container como volume, então você edita no seu editor e
roda com `make` sem precisar reconstruir a imagem.

## Estrutura

```
.
├── capitulos/
│   └── NN-nome/
│       ├── README.md          # teoria + comparação com JS
│       ├── main.go            # exemplo executável
│       └── exercicios/        # exercícios + testes
├── projetos/                  # projetos ao final de cada módulo
├── Dockerfile                 # imagem com o toolchain do Go
├── compose.yml                # serviços (Go, e mais tarde Postgres etc.)
├── Makefile                   # atalhos para rodar tudo
└── go.mod                     # o "package.json" do Go
```
