# 01 — Sintaxe essencial

```bash
make run CAP=01
```

Código: [`variaveis.go`](variaveis.go) · [`controle.go`](controle.go) · [`funcoes.go`](funcoes.go) · [`formatacao.go`](formatacao.go)

Todos os arquivos de uma pasta formam **um único pacote** e se enxergam sem import.

## Variáveis

```go
var nome string = "Go"   // tipo explícito
x := 10                  // forma curta (≈ let), a mais usada; só dentro de funções
a, b := 1, 2             // múltiplas
const limite = 50        // constante
```

- **Zero values:** variável sem valor recebe o "zero" do tipo: `0`, `""`, `false` ou `nil`. Não existe `undefined`.
- **Variável ou import não usado é erro de compilação.**
- Tipos do dia a dia: `int`, `float64`, `string`, `bool`.

## Conversões

Nada é convertido automaticamente: `"1" + 1` e `int * float64` não compilam.

```go
float64(idade) * altura          // numérico → numérico: Tipo(valor)
strconv.Itoa(30)                 // int → string     (≈ String(30))
n, err := strconv.Atoi("42")     // string → int     (≈ Number("42")), pode falhar
```

## "Enum" com iota

```go
type Status int

const (
	Pendente  Status = iota // 0
	Ativo                   // 1
	Cancelado               // 2
)
```

`iota` é um contador (0, 1, 2...) dentro do bloco `const`.

## if e switch

```go
if n > 0 { ... }                              // sem parênteses, chaves obrigatórias

if v, err := strconv.Atoi(s); err != nil {    // if com inicialização: muito comum
	...
}
```

- Sem truthy/falsy: `if lista` não compila. Use `if len(lista) > 0`.
- Sem ternário (`a ? b : c`). Use `if`/`else`.

```go
switch metodo {
case "POST", "PUT":       // vários valores; sem break
	...
default:
	...
}

switch {                  // sem expressão = cadeia de if/else if
case status >= 500:
	...
case status >= 400:
	...
}
```

## for (o único loop)

| JavaScript                          | Go                          |
| ----------------------------------- | --------------------------- |
| `for (let i = 0; i < n; i++)`       | `for i := 0; i < n; i++`    |
| `while (cond)`                      | `for cond`                  |
| `while (true)`                      | `for`                       |
| `arr.forEach((x, i) => ...)`        | `for i, x := range arr`     |
| `for (const x of arr)`              | `for _, x := range arr`     |

`_` descarta um valor que você não vai usar.

## Funções

```go
func soma(a, b int) int {
	return a + b
}

func dividir(a, b float64) (float64, error) {   // múltiplos retornos
	if b == 0 {
		return 0, errors.New("divisão por zero")
	}
	return a / b, nil
}
```

**O padrão mais importante do Go.** Não existe `try/catch`: funções que podem falhar
devolvem `(resultado, error)` e quem chama verifica:

```go
resultado, err := dividir(10, 0)
if err != nil {
	return err   // ou trata o erro
}
// segue usando resultado
```

| JavaScript                    | Go                                            |
| ----------------------------- | --------------------------------------------- |
| `function f(...nums)`         | `func f(nums ...int)`                         |
| `f(...arr)`                   | `f(arr...)`                                   |
| `const f = (n) => n * 2`      | `f := func(n int) int { return n * 2 }`       |
| parâmetro com valor padrão    | não existe                                    |

Closures funcionam igual ao JS.

## defer

Executa uma chamada quando a função termina, por qualquer caminho (≈ `finally`):

```go
arquivo, err := os.Open("dados.txt")
if err != nil {
	return err
}
defer arquivo.Close()   // fecha no fim da função, aconteça o que acontecer
```

Vários `defer` executam em ordem inversa.

## fmt

| Verbo   | Uso                                    |
| ------- | -------------------------------------- |
| `%v`    | valor padrão (≈ `console.log`)         |
| `%+v`   | struct com nomes dos campos (debug)    |
| `%T`    | tipo                                   |
| `%d`    | inteiro                                |
| `%.2f`  | decimal com 2 casas                    |
| `%s` / `%q` | string / string entre aspas        |

`fmt.Printf` imprime; `fmt.Sprintf` devolve a string (≈ template literal).
