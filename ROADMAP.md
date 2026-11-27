# 🗺️ Roadmap — Go para APIs REST

Trilha direta para **entender código Go e construir uma API REST**, escrita para
quem já programa em JavaScript. Sem exercícios: cada capítulo é um README objetivo
mais código executável.

A partir do capítulo 02, **todos os capítulos evoluem a mesma API** (um gerenciador
de tarefas). Cada pasta guarda a API naquele estágio, então dá para rodar e comparar
qualquer etapa.

## Capítulos

| #   | Capítulo                           | O que você aprende                                                        | Estado da API                                      |
| --- | ---------------------------------- | ------------------------------------------------------------------------- | -------------------------------------------------- |
| 00  | Ambiente                           | toolchain `go`, `go.mod`, Docker                                          | —                                                  |
| 01  | Sintaxe essencial                  | variáveis, tipos, `if`/`for`/`switch`, funções, `(valor, error)`, `defer` | —                                                  |
| 02  | Structs, ponteiros, slices e maps  | modelar dados, `&` e `*`, métodos, coleções                               | modelo `Tarefa` e armazenamento em memória         |
| 03  | HTTP e JSON                        | `net/http`, rotas, handlers, struct tags, status codes                    | CRUD em memória respondendo JSON                   |
| 04  | Interfaces e erros                 | interfaces implícitas, `errors.Is`, `%w`, erros → status HTTP             | camada de repositório; respostas 400/404/500       |
| 05  | Banco de dados                     | Postgres no Docker Compose, `pgx`, `context`                              | dados no Postgres, sem mudar os handlers           |
| 06  | Pronta para produção               | middleware, config por env, `log/slog`, Dockerfile, testes com `httptest` | API completa, pronta para o portfólio              |

## Extras (opcionais, depois da trilha)

- Concorrência: goroutines, channels, `sync`, `context` a fundo
- Generics
- Testes avançados e mocks
- Estrutura de projetos grandes (`cmd/`, `internal/`, injeção de dependência)

## Fluxo

1. Ler o `README.md` do capítulo.
2. Rodar: `make run CAP=NN`.
3. Commit + push: `git commit -m "feat(cap-NN): ..."` → `git push`.

## Referências

- [Go by Example](https://gobyexample.com/): exemplos curtos por tema
- [Documentação da biblioteca padrão](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go): como escrever Go idiomático
