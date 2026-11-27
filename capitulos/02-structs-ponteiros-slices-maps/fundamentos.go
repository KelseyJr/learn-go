package main

import (
	"fmt"
	"slices"
)

// ── Structs ─────────────────────────────────────────────

type Endereco struct {
	Cidade string
	UF     string
}

type Usuario struct {
	Nome     string
	Email    string
	Idade    int
	Endereco Endereco // struct dentro de struct
}

func exemploStructs() {
	u := Usuario{
		Nome:     "Ana",
		Email:    "ana@exemplo.com",
		Endereco: Endereco{Cidade: "Recife", UF: "PE"},
	}
	u.Idade = 28 // campos são lidos e alterados com ponto

	fmt.Printf("%+v\n", u)
	fmt.Println("cidade:", u.Endereco.Cidade)

	var vazio Usuario // zero value: todos os campos zerados
	fmt.Printf("zero value: %+v\n", vazio)

	// Struct é VALOR: atribuir copia tudo (no JS, copiaria só a referência)
	copia := u
	copia.Nome = "Bia"
	fmt.Println("original:", u.Nome, "| cópia:", copia.Nome)
}

// ── Ponteiros ───────────────────────────────────────────

func aniversarioCopia(u Usuario) { u.Idade++ } // recebe uma cópia
func aniversario(u *Usuario)     { u.Idade++ } // recebe o endereço do original

func exemploPonteiros() {
	u := Usuario{Nome: "Ana", Idade: 28}

	aniversarioCopia(u)
	fmt.Println("depois de aniversarioCopia(u):", u.Idade)

	aniversario(&u)
	fmt.Println("depois de aniversario(&u):    ", u.Idade)

	p := &u              // p é um *Usuario
	p.Nome = "Ana Maria" // altera o próprio u
	fmt.Println("u.Nome:", u.Nome)

	var nulo *Usuario // zero value de ponteiro é nil
	fmt.Println("ponteiro sem valor:", nulo, "(acessar nulo.Nome causaria panic)")
}

// ── Métodos ─────────────────────────────────────────────

type Contador struct {
	total int
}

func (c *Contador) Incrementar() { c.total++ }
func (c *Contador) Valor() int   { return c.total }

func exemploMetodos() {
	c := Contador{}
	c.Incrementar() // o Go faz (&c).Incrementar() sozinho
	c.Incrementar()
	fmt.Println("contador:", c.Valor())
}

// ── Slices ──────────────────────────────────────────────

func exemploSlices() {
	nums := []int{3, 1, 2}    // ≈ [3, 1, 2]
	nums = append(nums, 4, 5) // ≈ push, mas DEVOLVE o slice: sempre reatribua
	fmt.Println(nums, "len:", len(nums))

	fmt.Println("nums[0]:", nums[0], "| nums[1:3]:", nums[1:3]) // ≈ slice(1, 3)

	fmt.Println("contém 2?", slices.Contains(nums, 2)) // ≈ includes
	fmt.Println("índice do 4:", slices.Index(nums, 4)) // ≈ indexOf
	slices.Sort(nums)                                  // ≈ sort(), altera o próprio slice
	fmt.Println("ordenado:", nums)

	// Não existe .filter/.map/.reduce: usa-se for
	var pares []int // slice nil: funciona com append, len e range
	for _, n := range nums {
		if n%2 == 0 {
			pares = append(pares, n)
		}
	}
	fmt.Println("pares:", pares)

	// make(tipo, tamanho, capacidade): reserva espaço quando você sabe quantos itens terá
	dobros := make([]int, 0, len(nums))
	for _, n := range nums {
		dobros = append(dobros, n*2)
	}
	fmt.Println("dobros:", dobros)

	// Pegadinha: um pedaço (nums[a:b]) compartilha a memória do slice original
	pedaco := nums[0:2]
	pedaco[0] = 99
	fmt.Println("nums depois de alterar pedaco[0]:", nums)
}

// ── Maps ────────────────────────────────────────────────

func exemploMaps() {
	estoque := map[string]int{ // ≈ new Map([["banana", 10], ["maçã", 5]])
		"banana": 10,
		"maçã":   5,
	}
	estoque["uva"] = 8        // inserir ou atualizar
	delete(estoque, "banana") // remover
	fmt.Println(estoque, "len:", len(estoque))

	// Chave inexistente devolve o zero value (0), não undefined
	fmt.Println("pera:", estoque["pera"])

	// "comma ok": diferencia "não existe" de "existe com valor zero" (≈ map.has)
	if qtd, ok := estoque["maçã"]; ok {
		fmt.Println("maçã existe:", qtd)
	}

	// A ordem do range em um map é ALEATÓRIA (muda a cada execução)
	for fruta, qtd := range estoque {
		fmt.Print(fruta, "=", qtd, " ")
	}
	fmt.Println("← ordem aleatória")

	// Zero value 0 permite incrementar direto, sem checar se a chave existe
	contagem := make(map[string]int)
	for _, palavra := range []string{"go", "js", "go"} {
		contagem[palavra]++
	}
	fmt.Println("contagem:", contagem)
}
