package main

import (
	"fmt"
	"strconv"
)

func exemploVariaveis() {
	var linguagem string = "Go" // tipo explícito
	versao := 1.27              // forma curta (≈ let): a mais usada, só dentro de funções
	a, b := 1, 2                // múltiplas de uma vez

	fmt.Println(linguagem, versao, a, b)

	// Zero values: toda variável nasce com um valor. Não existe undefined.
	var (
		i  int
		s  string
		ok bool
	)
	fmt.Printf("zero values → int: %d | string: %q | bool: %t\n", i, s, ok)

	// Divisão entre inteiros descarta a parte decimal (no JS, 7 / 2 === 3.5)
	fmt.Println("7 / 2 =", 7/2, "| 7.0 / 2 =", 7.0/2)
}

func exemploConversoes() {
	idade := 30
	altura := 1.75

	// int e float64 não se misturam: conversão é sempre explícita
	fmt.Printf("float64(idade) * altura = %.2f\n", float64(idade)*altura)

	// string ↔ número (≈ String(n) / Number(s))
	texto := strconv.Itoa(idade)
	numero, err := strconv.Atoi("42") // pode falhar, então devolve um erro
	fmt.Printf("Itoa: %q | Atoi: %d (err: %v)\n", texto, numero, err)

	_, err = strconv.Atoi("abc") // `_` descarta o valor
	fmt.Println(`Atoi("abc") →`, err)
}

// "Enum" em Go: tipo nomeado + constantes com iota (contador 0, 1, 2...)
type Status int

const (
	Pendente  Status = iota // 0
	Ativo                   // 1
	Cancelado               // 2
)

func exemploConstantes() {
	const limitePorPagina = 50
	var s Status = Ativo
	fmt.Println("limite:", limitePorPagina, "| status:", s)
}
