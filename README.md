# Aprendendo Go 🐹

Estudos de **Go (Golang)** focados em entender código Go e construir uma **API REST**,
escritos para quem vem do **JavaScript**. Cada capítulo é um README direto com código executável.

👉 **[Veja o roadmap](ROADMAP.md)**

## Requisitos

- [Docker](https://www.docker.com/) e `make` (todo o código roda em containers)
- *Opcional:* Go local (`brew install go`) para o editor ter autocomplete

## Como usar

```bash
make help            # lista os comandos
make list            # lista os capítulos
make run CAP=01      # roda um capítulo
make test CAP=06     # roda os testes de um capítulo
make shell           # terminal dentro do container Go
make fmt             # formata o código (gofmt)
make vet             # análise estática (go vet)
```

O código é montado no container como volume: edite no editor e rode com `make`, sem rebuild.

## Projetos

- **[TODO API](projetos/todo-api)**: API REST de tarefas com a estrutura de um projeto
  Go real (`cmd/`, `internal/`, testes, Dockerfile de produção). Rode com `make todo`.

## Estrutura

```
.
├── capitulos/
│   └── NN-nome/
│       ├── README.md      # conteúdo do capítulo
│       └── *.go           # código executável
├── projetos/
│   └── todo-api/          # API REST completa
├── Dockerfile             # imagem com o toolchain do Go
├── compose.yml            # serviços (Go e, mais tarde, Postgres)
├── Makefile               # atalhos
└── go.mod                 # o "package.json" do Go
```
