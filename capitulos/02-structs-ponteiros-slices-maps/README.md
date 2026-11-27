# 02 — Structs, ponteiros, slices e maps

```bash
make run CAP=02
```

Código: [`fundamentos.go`](fundamentos.go) (exemplos) · [`tarefa.go`](tarefa.go) + [`repo.go`](repo.go) (início da API)

## Structs

O formato dos dados. Equivale a uma `interface` do TS + o objeto, só que com campos fixos.

```go
type Usuario struct {
	Nome  string
	Idade int
}

u := Usuario{Nome: "Ana", Idade: 28}   // ≈ { nome: "Ana", idade: 28 }
u.Idade = 29
var vazio Usuario                      // zero value: {Nome: "", Idade: 0}
```

**Struct é valor.** `copia := u` copia todos os campos, enquanto no JS copiaria só a referência.

## Ponteiros

Em Go, **tudo é passado como cópia**. Para alterar o original, passe o endereço dele:

```go
func aniversarioCopia(u Usuario) { u.Idade++ }   // altera a cópia
func aniversario(u *Usuario)     { u.Idade++ }   // altera o original

aniversario(&u)
```

| Símbolo     | Significa                                      |
| ----------- | ---------------------------------------------- |
| `&u`        | endereço de `u` (cria um ponteiro)             |
| `*Usuario`  | tipo "ponteiro para Usuario"                   |
| `p.Nome`    | acessa o campo (o Go resolve o ponteiro sozinho) |
| `nil`       | ponteiro vazio; acessar um campo dele causa *panic* |

**Use ponteiro quando** precisar alterar o valor ou quando ele for compartilhado
(handler, repositório, conexão com o banco). **Use valor** para dados pequenos.

## Métodos

Go não tem `class`. Métodos são funções com um **receiver**, que faz o papel do `this`:

```go
type Repo struct {
	tarefas map[int]Tarefa
}

func NovoRepo() *Repo {                        // ≈ constructor (convenção: NovoX / NewX)
	return &Repo{tarefas: make(map[int]Tarefa)}
}

func (r *Repo) Criar(titulo string) Tarefa {   // ≈ criar(titulo) { this.tarefas... }
	...
}

repo := NovoRepo()      // ≈ new Repo()
repo.Criar("Estudar")
```

- `func (r *Repo)` é um receiver ponteiro e altera o original. É o padrão para quase tudo.
- `func (r Repo)` é um receiver por valor e trabalha numa cópia. Alterações se perdem.
- Não misture os dois tipos de receiver no mesmo tipo.

## Slices (≈ arrays do JS)

```go
nums := []int{3, 1, 2}
nums = append(nums, 4)          // ≈ push, mas DEVOLVE o slice: sempre reatribua
nums[0]                         // índice
nums[1:3]                       // ≈ slice(1, 3)
len(nums)                       // ≈ length
make([]int, 0, 10)              // vazio, com capacidade para 10 já reservada
```

| JavaScript                    | Go                                     |
| ----------------------------- | -------------------------------------- |
| `arr.push(x)`                 | `arr = append(arr, x)`                 |
| `arr.includes(x)`             | `slices.Contains(arr, x)`              |
| `arr.indexOf(x)`              | `slices.Index(arr, x)`                 |
| `arr.sort((a, b) => a - b)`   | `slices.SortFunc(arr, func(a, b T) int { return a - b })` |
| `arr.filter(...)` / `.map(...)` | não existe: use `for` + `append`     |

```go
var pares []int
for _, n := range nums {
	if n%2 == 0 {
		pares = append(pares, n)
	}
}
```

⚠️ `nums[a:b]` **não copia**: o pedaço compartilha memória com o original, e alterar um altera o outro.

## Maps (≈ `Map` / objeto do JS)

```go
estoque := map[string]int{"banana": 10}   // literal
estoque := make(map[string]int)           // vazio
estoque["uva"] = 8                        // inserir/atualizar
delete(estoque, "uva")                    // remover
qtd := estoque["pera"]                    // não existe → zero value (0), não undefined

if qtd, ok := estoque["maçã"]; ok {       // "comma ok" ≈ map.has + map.get
	...
}
```

⚠️ **Três pegadinhas:**
1. A ordem do `range` num map é **aleatória**. Para ordenar, copie para um slice e ordene.
2. Um map declarado sem `make` é `nil`: ler funciona, mas escrever causa *panic*.
3. Num `map[int]Tarefa` não dá para alterar um campo direto: `m[id].Concluida = true` não compila.
   Pegue a tarefa, altere e salve de volta:
   ```go
   t := m[id]
   t.Concluida = true
   m[id] = t
   ```

## A API até aqui

- [`tarefa.go`](tarefa.go): a struct `Tarefa` com os métodos `Concluir` e `Resumo`.
- [`repo.go`](repo.go): o `Repo`, que guarda as tarefas num `map[int]Tarefa`, com os métodos:

| Método                    | Retorno            | Observação                                         |
| ------------------------- | ------------------ | -------------------------------------------------- |
| `Criar(titulo)`           | `Tarefa`           | gera o ID sequencial                               |
| `Buscar(id)`              | `(Tarefa, bool)`   | "comma ok"; vira `(Tarefa, error)` no capítulo 04  |
| `Listar()`                | `[]Tarefa`         | ordenado por ID (o map não tem ordem)              |
| `Atualizar(t)`            | `bool`             | `false` se o ID não existir                        |
| `Remover(id)`             | `bool`             | `false` se o ID não existir                        |

`Buscar` devolve uma **cópia** da tarefa. Por isso, para alterar: buscar → modificar → `Atualizar`.

**Próximo capítulo:** expor esse `Repo` por HTTP, com rotas `GET`/`POST`/`PUT`/`DELETE` respondendo JSON.
