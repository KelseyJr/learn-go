# 🗺️ Roadmap — Go do básico ao avançado

Plano de estudos pensado para quem já programa em **JavaScript/TypeScript** e quer
usar Go profissionalmente (backend, APIs, CLIs, cloud e ferramentas de infraestrutura).

Cada capítulo traz uma ponte com o JS: o que é igual, o que é diferente e o que
o JS fazia "escondido" e o Go obriga você a enxergar.

---

## Como cada capítulo funciona

Cada capítulo é uma pasta em `capitulos/NN-nome/`:

```
capitulos/05-slices-e-maps/
├── README.md              # teoria + comparação com JS + armadilhas + ATIVIDADES
├── main.go                # código comentado e executável (make run CAP=05)
├── exercicios/
│   ├── exercicios.go      # funções com TODO para VOCÊ implementar
│   └── exercicios_test.go # testes prontos que validam sua resposta (make ex CAP=05)
└── desafio/               # criada por VOCÊ: solução do desafio do capítulo
```

### Atividades de fim de capítulo

Todo `README.md` de capítulo termina com uma seção **Atividades**, em níveis:

| Nível | Atividade              | Como é validada                                                       |
| ----- | ---------------------- | --------------------------------------------------------------------- |
| 1     | 🧪 **Exercícios**      | testes automáticos: `make ex CAP=NN` precisa ficar verde              |
| 2     | 🔍 **Explorar**        | rodar comandos e observar o comportamento (sem resposta "certa")      |
| 3     | 💥 **Quebre de propósito** | provocar erros e aprender a ler as mensagens do compilador        |
| 4     | 🏆 **Desafio**         | um programa pequeno escrito do zero em `desafio/`; eu reviso o código |
| 5     | ❓ **Perguntas**       | perguntas de fixação; responda com suas palavras antes de avançar     |

O capítulo só está concluído quando os exercícios passam e o desafio foi revisado.

### Fluxo de trabalho

1. **Ler** o `README.md` do capítulo.
2. **Rodar e mexer** no `main.go`: `make run CAP=05`.
3. **Fazer as atividades** (exercícios até `make ex CAP=05` ficar verde ✅, depois o desafio).
4. **Pedir revisão** do desafio e tirar dúvidas antes de avançar.
5. **Commit + push** fechando o capítulo:
   `git commit -m "feat(cap-05): slices e maps"` → `git push`.

### Projetos

Ao final de cada módulo há um **projeto**, em `projetos/NN-nome/`. Nos projetos
**você escreve o código** e eu faço o papel de revisor (code review como no trabalho).
É aqui que o conhecimento vira habilidade, e são eles que vão compor o portfólio.

---

## Módulo 1 — Fundamentos

A base da linguagem. Objetivo: ler e escrever Go simples com confiança.

| #   | Capítulo                     | Conteúdo                                                                                | Ponte com JS                                                                     |
| --- | ---------------------------- | --------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 00  | Ambiente                     | toolchain `go`, `go.mod`, Docker, `gofmt`, `go vet`                                     | `go` = node + npm + prettier + eslint + jest num único binário                  |
| 01  | Olá, Go                      | `package main`, `func main`, `fmt`, imports, compilação                                 | `node index.js` vs binário compilado                                            |
| 02  | Variáveis e tipos            | `var`, `:=`, `const`, `iota`, tipos básicos, *zero values*, conversão explícita         | `let`/`const`, mas sem `undefined` e sem coerção (`"1" + 1`)                    |
| 03  | Controle de fluxo            | `if` com inicialização, `for` (o único loop), `switch`, `defer`                          | um `for` faz o papel de `for`, `while` e `for...of`                             |
| 04  | Funções                      | múltiplos retornos, variádicas, funções como valor, closures                            | arrow functions e closures funcionam igual                                      |
| 05  | Arrays, slices e maps        | `len`/`cap`, `append`, fatiamento, `map`, `range`                                       | `Array` → slice, `Object`/`Map` → map                                           |
| 06  | Strings, bytes e runes       | UTF-8, imutabilidade, `strings`, `strconv`, `strings.Builder`                           | por que `len("ção")` não é 3                                                    |
| 07  | Structs e métodos            | structs, métodos, *receivers*, embedding (composição)                                   | `class` sem herança: composição no lugar de `extends`                           |
| 08  | Ponteiros                    | `&` e `*`, valor vs referência, quando usar ponteiro                                    | o que o JS faz escondido quando você passa um objeto                            |
| 09  | Interfaces                   | implementação implícita, interfaces pequenas, `any`, type switch                        | tipagem estrutural do TS, mas checada em tempo de compilação                    |
| 10  | Erros                        | `error` como valor, `errors.Is`/`As`, `%w`, `panic`/`recover`                           | adeus `try/catch`: olá `if err != nil`                                          |
| 11  | Pacotes e módulos            | visibilidade por maiúscula, `internal/`, `go get`, versionamento                        | `export` → letra maiúscula; `go.mod` ≈ `package.json`                           |
| 12  | Arquivos e JSON              | `os`, `encoding/json`, struct tags, leitura/escrita de arquivos                         | `JSON.parse`/`stringify` + `fs`                                                 |

🛠️ **Projeto 1 — CLI de tarefas**: um *todo* de linha de comando com persistência em JSON.

---

## Módulo 2 — Go idiomático e concorrência

O que faz o Go brilhar. Objetivo: escrever código testado e concorrente com segurança.

| #   | Capítulo                     | Conteúdo                                                                                | Ponte com JS                                                                     |
| --- | ---------------------------- | --------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 13  | Generics e iteradores        | type parameters, constraints, `iter.Seq`, range over func                               | `<T>` do TS; iteradores ≈ `function*` (generators)                              |
| 14  | Testes                       | `testing`, table-driven tests, subtests, benchmarks, cobertura, fuzzing                 | Jest/Vitest, só que embutido na linguagem                                       |
| 15  | Goroutines                   | `go`, `sync.WaitGroup`, scheduler do Go                                                  | paralelismo real vs event loop single-thread                                    |
| 16  | Channels e select            | channels com/sem buffer, `close`, `range`, `select`, timeouts                           | `Promise.race`, async iterators, filas                                          |
| 17  | Sincronização                | `Mutex`, `RWMutex`, `Once`, `atomic`, race detector                                      | problemas que o JS (single-thread) nunca te deu                                 |
| 18  | Context                      | cancelamento, deadlines, propagação entre chamadas                                      | `AbortController` / `AbortSignal`                                               |
| 19  | Padrões de concorrência      | worker pool, pipeline, fan-out/fan-in, `errgroup`, rate limiting                        | `Promise.all` com limite (tipo `p-limit`)                                       |
| 20  | I/O e CLIs                   | `io.Reader`/`Writer`, `bufio`, `flag`, variáveis de ambiente, sinais do SO              | streams do Node, `process.argv`, `process.env`                                  |

🛠️ **Projeto 2 — Verificador de links concorrente**: CLI que checa centenas de URLs
em paralelo com worker pool, timeouts via `context` e relatório final.

---

## Módulo 3 — Backend profissional

Go no dia a dia de trabalho. Objetivo: construir e colocar em produção uma API real.

| #   | Capítulo                     | Conteúdo                                                                                | Ponte com JS                                                                     |
| --- | ---------------------------- | --------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 21  | HTTP com net/http            | servidor, roteamento com método e path params, handlers, respostas JSON                 | Express, mas sem precisar de framework                                          |
| 22  | Middlewares e APIs REST      | logging, recovery, CORS, autenticação, validação                                        | `app.use()`                                                                     |
| 23  | Banco de dados               | `database/sql`, `pgx`, Postgres no Docker Compose, migrations, transações               | `pg` / Prisma / Knex                                                            |
| 24  | Arquitetura de projetos      | `cmd/`, `internal/`, camadas, injeção de dependência manual, configuração               | a organização do NestJS, mas sem mágica                                         |
| 25  | Logs e observabilidade       | `log/slog`, health checks, métricas, introdução ao OpenTelemetry                        | pino / winston                                                                  |
| 26  | Testes de integração         | `httptest`, testcontainers, mocks via interfaces                                        | supertest                                                                       |
| 27  | Produção                     | graceful shutdown, Dockerfile multi-stage, imagem distroless, cross-compile             | `node_modules` de 500 MB vs binário de 10 MB                                    |

🛠️ **Projeto 3 — API REST completa**: encurtador de URLs (ou gerenciador de tarefas)
com Postgres, autenticação, testes, Docker e CI no GitHub Actions.

---

## Módulo 4 — Avançado

Entender o que acontece por baixo dos panos. Objetivo: nível sênior.

| #   | Capítulo                     | Conteúdo                                                                                | Ponte com JS                                                                     |
| --- | ---------------------------- | --------------------------------------------------------------------------------------- | -------------------------------------------------------------------------------- |
| 28  | Reflection                   | pacote `reflect`, struct tags na prática, como o `encoding/json` funciona por dentro    | `Object.keys`, `Reflect`, decorators                                            |
| 29  | Performance                  | `pprof`, benchmarks, escape analysis, alocações, `sync.Pool`                            | Chrome DevTools / `--inspect`, mas muito mais fundo                             |
| 30  | Runtime por dentro           | scheduler (modelo GMP), garbage collector, memory model                                 | como o V8 e o libuv funcionam, só que no Go                                     |
| 31  | gRPC e Protobuf              | contratos com `.proto`, geração de código, streaming                                    | REST/tRPC entre microsserviços                                                  |
| 32  | Ferramentas e qualidade      | `golangci-lint`, `go generate`, `go:embed`, build tags, GitHub Actions                  | ESLint, scripts do `package.json`, bundlers                                     |

🏁 **Projeto final — você escolhe**: um sistema com API + worker assíncrono + fila
(ou uma ideia sua), com tudo que foi aprendido. É o projeto de vitrine do portfólio.

---

## Progresso

- [ ] Módulo 1 — Fundamentos (00–12) + Projeto 1
- [ ] Módulo 2 — Go idiomático e concorrência (13–20) + Projeto 2
- [ ] Módulo 3 — Backend profissional (21–27) + Projeto 3
- [ ] Módulo 4 — Avançado (28–32) + Projeto final

## Referências oficiais

- [A Tour of Go](https://go.dev/tour/): exercícios interativos no navegador
- [Effective Go](https://go.dev/doc/effective_go): como escrever Go idiomático
- [Go by Example](https://gobyexample.com/): exemplos curtos por tema
- [Documentação da biblioteca padrão](https://pkg.go.dev/std)
