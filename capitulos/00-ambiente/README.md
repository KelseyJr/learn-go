# Capítulo 00 — Ambiente

> **Objetivo:** entender o que é Go, como o toolchain funciona e como este
> repositório roda tudo dentro do Docker.

## O que é Go?

Go é uma linguagem **compilada**, **estaticamente tipada** e com **garbage collector**,
criada no Google em 2009. Ela foi pensada para ser simples (a especificação inteira
é pequena), rápida para compilar e boa em concorrência.

É a linguagem por trás do **Docker**, **Kubernetes**, **Terraform**, **Prometheus**
e de boa parte da infraestrutura de cloud. O próprio Docker que você está usando
para rodar este repositório foi escrito em Go.

## Compilado vs interpretado

| JavaScript (Node)                                        | Go                                                            |
| -------------------------------------------------------- | ------------------------------------------------------------- |
| O código é lido pelo V8 e compilado *em tempo de execução* (JIT) | O código é compilado *antes* para um binário nativo           |
| Erros de tipo aparecem quando o código roda (ou no TS)   | Erros de tipo impedem a compilação                            |
| O servidor precisa ter o Node instalado                  | O binário roda sozinho: não precisa instalar nada no servidor |
| Deploy = código + `node_modules`                         | Deploy = **um único arquivo** executável                      |

## O toolchain: um binário que substitui vários

No ecossistema JS você junta várias ferramentas. No Go, quase tudo vem no comando `go`:

| O que você faz no JS            | Em Go                     |
| ------------------------------- | ------------------------- |
| `npm init`                      | `go mod init <módulo>`    |
| `node index.js`                 | `go run .`                |
| `tsc` / bundler (build)         | `go build`                |
| `npm install pacote`            | `go get pacote`           |
| `npm prune` / limpar deps       | `go mod tidy`             |
| Prettier                        | `gofmt` (sem configuração, um estilo só para todo mundo) |
| ESLint                          | `go vet`                  |
| Jest / Vitest                   | `go test`                 |

> 💡 O `gofmt` não tem opções de configuração de propósito. Não existe briga de
> "tabs vs espaços" ou "ponto e vírgula ou não" em time Go: todo código Go do
> mundo tem a mesma cara.

## Módulos: o `package.json` do Go

O arquivo [`go.mod`](../../go.mod) na raiz define o **módulo**:

```go
module github.com/KelseyJr/learn-go   // nome do módulo (≈ "name" no package.json)

go 1.27                                // versão mínima do Go
```

- Quando você adiciona dependências, elas aparecem no `go.mod` e os hashes ficam
  no `go.sum` (≈ `package-lock.json`).
- **Não existe `node_modules`**: as dependências ficam num cache global da máquina
  (no nosso caso, num volume do Docker) e são compartilhadas entre projetos.
- O nome do módulo costuma ser a URL do repositório. É assim que o `go get`
  sabe de onde baixar um pacote.

## Como o Docker deste repositório funciona

- [`Dockerfile`](../../Dockerfile): imagem com o Go 1.27 instalado.
- [`compose.yml`](../../compose.yml): monta a pasta do projeto dentro do container,
  então você edita os arquivos normalmente no seu editor.
- [`Makefile`](../../Makefile): atalhos para não precisar decorar comandos do Docker.

Por baixo dos panos, `make run CAP=00` executa:

```bash
docker compose run --rm go go run ./capitulos/00-ambiente
```

## Anatomia do `main.go`

```go
package main          // todo executável Go pertence ao pacote "main"

import (              // imports da biblioteca padrão (≈ require/import do Node)
	"fmt"
	"runtime"
)

func main() {         // ponto de entrada do programa (no Node, é o arquivo todo)
	fmt.Println("...")
}
```

Detalhes que já valem notar:

- Não tem `;` no fim das linhas (o compilador insere sozinho).
- A `{` precisa ficar **na mesma linha** do `func`/`if`/`for`. Na linha de baixo é erro de compilação.
- **Import não usado é erro de compilação**, não um warning. Tente apagar as quatro
  linhas que usam `runtime.` e rode de novo: o compilador recusa o código porque
  o import `"runtime"` ficou sobrando.

## Rodando

```bash
make run CAP=00   # roda o programa
make ex  CAP=00   # roda os testes do exercício (vai falhar até você resolver)
```

A primeira execução demora um pouco porque o Docker baixa a imagem do Go.

## Atividades

### 🧪 1. Exercícios

Abra [`exercicios/exercicios.go`](exercicios/exercicios.go) e faça a função
`Saudacao` retornar `"Olá, Go!"`. Depois rode `make ex CAP=00` até aparecer `PASS`.

### 🔍 2. Explorar

Abra um terminal dentro do container com `make shell` e rode:

```bash
go version                      # versão do Go
go env GOOS GOARCH GOMODCACHE   # onde ficam as dependências (o "node_modules global")
go help                         # todos os subcomandos do toolchain

# Compile o capítulo para um binário de verdade e rode-o:
go build -o bin/info ./capitulos/00-ambiente
./bin/info
ls -lh bin/info                 # quanto ele pesa?

# Compile o MESMO código para Windows sem sair do Linux (cross-compile):
GOOS=windows GOARCH=amd64 go build -o bin/info.exe ./capitulos/00-ambiente
ls -lh bin/
```

Compare o tamanho do binário com o de uma pasta `node_modules` de um projeto seu.
(A pasta `bin/` está no `.gitignore`, então nada disso vai para o GitHub.)

### 💥 3. Quebre de propósito

Faça uma alteração por vez no `main.go`, rode `make run CAP=00`, leia o erro e desfaça:

1. Apague as quatro linhas que usam `runtime.` (o import fica sobrando).
2. Mova a `{` do `func main()` para a linha de baixo.
3. Troque `fmt.Println` por `fmt.println` (com `p` minúsculo).
4. Bagunce a indentação do arquivo e rode `make fmt`. O que aconteceu?

> O erro 3 é uma prévia do capítulo 11: em Go, só o que começa com letra
> **maiúscula** é exportado (público) de um pacote.

### 🏆 4. Desafio

Crie `capitulos/00-ambiente/desafio/main.go` (com `package main`) que imprima um
"cartão de apresentação" seu no terminal: nome, cidade, stack atual e por que
está aprendendo Go. Use pelo menos uma vez `fmt.Println` e uma vez `fmt.Printf`.

Rode com `docker compose run --rm go go run ./capitulos/00-ambiente/desafio`.

**Bônus:** consulte a documentação do pacote [`runtime`](https://pkg.go.dev/runtime)
e adicione ao cartão alguma informação que não aparece no `main.go` do capítulo.

### ❓ 5. Perguntas

1. O que precisa estar instalado num servidor para rodar um programa Go já compilado? E um programa Node?
2. Qual a diferença entre `go run` e `go build`?
3. Por que o `gofmt` não ter opções de configuração é considerado uma vantagem?
4. Onde ficam as dependências de um projeto Go, já que não existe `node_modules`?

## Resumo

- [ ] Sei a diferença entre Go (compilado) e Node (interpretado/JIT)
- [ ] Sei para que servem `go run`, `go build`, `go test`, `gofmt` e `go vet`
- [ ] Entendi que o `go.mod` faz o papel do `package.json`
- [ ] Rodei `make run CAP=00` e `make ex CAP=00` com sucesso
- [ ] Fiz o desafio e respondi as perguntas
