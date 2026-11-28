# TODO API

API REST de tarefas em Go, só com a biblioteca padrão e sem banco de dados
(os dados ficam em memória). A estrutura segue o que você encontra em projetos Go reais.

```bash
make todo        # sobe a API em http://localhost:8080 (Ctrl+C para parar)
make todo-test   # roda os testes
```

## Endpoints

| Método   | Rota            | Corpo                                  | Respostas          |
| -------- | --------------- | -------------------------------------- | ------------------ |
| `GET`    | `/health`       | —                                      | 200                |
| `GET`    | `/tasks`        | —                                      | 200                |
| `POST`   | `/tasks`        | `{"title": "Estudar Go"}`              | 201, 400           |
| `GET`    | `/tasks/{id}`   | —                                      | 200, 400, 404      |
| `PATCH`  | `/tasks/{id}`   | `{"title": "...", "done": true}` (ambos opcionais) | 200, 400, 404 |
| `DELETE` | `/tasks/{id}`   | —                                      | 204, 400, 404      |

Erros sempre no formato `{"error": "mensagem"}`.

```bash
curl -X POST localhost:8080/tasks -d '{"title": "Estudar Go"}'
curl localhost:8080/tasks
curl -X PATCH localhost:8080/tasks/1 -d '{"done": true}'
curl -X DELETE localhost:8080/tasks/1
```

## Estrutura

```
projetos/todo-api/
├── go.mod                     # module todo (≈ package.json)
├── cmd/api/main.go            # ponto de entrada: monta as peças e sobe o servidor
├── internal/
│   ├── task/                  # tudo sobre tarefas
│   │   ├── model.go           # modelo, erros e validação
│   │   ├── store.go           # armazenamento em memória
│   │   ├── handler.go         # rotas HTTP
│   │   └── handler_test.go    # testes das rotas
│   └── httpx/                 # utilitários HTTP compartilhados
│       ├── json.go            # ler e escrever JSON
│       └── middleware.go      # log de requisições
└── Dockerfile                 # imagem de produção (~16 MB)
```

| Go                              | Equivalente num projeto Express          |
| ------------------------------- | ---------------------------------------- |
| `cmd/api/main.go`               | `server.js`                              |
| `internal/task/handler.go`      | `routes/tasks.js` + controller           |
| `internal/task/store.go`        | repository / model                       |
| `internal/task/model.go`        | model + validação (zod/joi)              |
| `internal/httpx/middleware.go`  | `morgan`                                 |
| `internal/task/handler_test.go` | jest + supertest                         |

- **`go.mod`** declara o módulo `todo`, e é por isso que os imports são `todo/internal/task`.
  O projeto é autocontido: copiado para outro lugar, funciona sozinho.
- **`cmd/`** guarda os executáveis (um projeto pode ter `cmd/api`, `cmd/worker`...).
- **`internal/`** é uma regra do compilador: o que está ali só pode ser importado
  por este projeto, nunca por outro módulo. É o "privado" no nível de pacote.

> 💡 **Pasta = pacote, arquivo = gaveta.** `task.NewMemoryStore()` se refere à **pasta**
> `internal/task/`, não a um arquivo. Todos os arquivos dela começam com `package task`
> e formam um único pacote; o nome do arquivo não importa para o Go. Para saber o pacote
> de um arquivo, olhe a primeira linha dele.

## Caminho de uma requisição

```
PATCH /tasks/1   {"done": true}
   │
   ▼
http.Server ─────────────── cmd/api/main.go
   │
   ▼
httpx.Logger ────────────── middleware: marca o início, chama o próximo, loga no fim
   │
   ▼
ServeMux: "PATCH /tasks/{id}" → Handler.update ─── handler.go
   │   parseID → DecodeJSON → store.Get → aplica os campos enviados → store.Update
   ▼
MemoryStore.Update ──────── store.go: trava o mutex e salva no map
   │
   ▼
httpx.WriteJSON → 200 {"id": 1, "title": "...", "done": true, ...}
```

## Ordem de leitura sugerida

1. [`model.go`](internal/task/model.go): o modelo e as regras.
2. [`store.go`](internal/task/store.go): onde os dados ficam.
3. [`handler.go`](internal/task/handler.go): as rotas, o coração da API.
4. [`httpx/`](internal/httpx): JSON e middleware.
5. [`cmd/api/main.go`](cmd/api/main.go): como tudo é montado.
6. [`handler_test.go`](internal/task/handler_test.go): os testes, que também documentam o comportamento.

## O que aparece no código, e onde

| Conceito                                | Onde                                  | Em JS/TS seria                       |
| --------------------------------------- | ------------------------------------- | ------------------------------------ |
| Struct tags `json:"..."`                | `model.go`                            | os nomes das props no JSON           |
| Erros sentinela + `errors.Is`           | `model.go`, `respondError` em `handler.go` | classes de erro + `instanceof`  |
| Interface implícita (`Store`)           | `handler.go`                          | interface do TS, sem `implements`    |
| Receiver ponteiro `(h *Handler)`        | todos os métodos                      | `this`                               |
| Ponteiro como campo opcional (`*bool`)  | `updateRequest` em `handler.go`       | `done?: boolean`                     |
| `&req` para o Decode preencher          | `create` em `handler.go`              | —                                    |
| `r.Context()` repassado ao store        | todos os handlers                     | `AbortSignal`                        |
| Mutex                                   | `store.go`                            | não existe (event loop single-thread) |
| Middleware                              | `middleware.go`                       | `app.use(morgan(...))`               |
| Embedding (`statusRecorder`)            | `middleware.go`                       | um `extends` que sobrescreve um método |
| Goroutine + channel + `select`          | `main.go` (graceful shutdown)         | —                                    |
| Generics (`decode[T]`)                  | `handler_test.go`                     | `function decode<T>()`               |
| Teste table-driven                      | `TestErrors`                          | `test.each`                          |

## Decisões de projeto

- **Só biblioteca padrão.** Desde o Go 1.22 o `ServeMux` aceita método e parâmetros
  na rota (`"GET /tasks/{id}"`). Muitos projetos usam `chi`, `gin` ou `echo`, mas a
  estrutura é a mesma.
- **`[]` em vez de `null`.** Um slice `nil` vira `null` no JSON, por isso o `List`
  cria a lista com `make([]Task, 0, ...)`.
- **O store recebe `ctx` mesmo sem banco.** Trocar por Postgres é criar um
  `PostgresStore` com os mesmos métodos e mudar **uma linha** no `main.go`.
  O handler não muda, porque depende da interface `Store`, não do tipo concreto.
- **Mutex no store.** Cada requisição roda numa goroutine própria, então o map é
  acessado em paralelo. Os testes rodam com o race detector (`go test -race`) sem alertas.
- **Graceful shutdown.** Ao receber Ctrl+C ou `docker stop`, o servidor para de
  aceitar conexões e espera as requisições em andamento terminarem.
- **Nomes em inglês, comentários em português.** Nomes em inglês são o padrão do mercado.
- **Receivers com uma letra** (`s`, `h`), como manda a convenção do Go: o receiver
  aparece em quase toda linha do método e o contexto já diz quem ele é (≈ o `this`).
