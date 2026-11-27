package main

import (
	"fmt"
	"strconv"
)

func exemploCondicionais() {
	// Sem parênteses, chaves obrigatórias, condição precisa ser bool (sem truthy/falsy)
	n := 7
	if n%2 == 0 {
		fmt.Println(n, "é par")
	} else {
		fmt.Println(n, "é ímpar")
	}

	// if com inicialização: o padrão mais comum em Go, principalmente com erros.
	// valor e err só existem dentro do if/else.
	if valor, err := strconv.Atoi("123"); err != nil {
		fmt.Println("erro:", err)
	} else {
		fmt.Println("convertido:", valor)
	}

	// switch: sem break, vários valores por case
	metodo := "PUT"
	switch metodo {
	case "GET":
		fmt.Println("leitura")
	case "POST", "PUT", "PATCH":
		fmt.Println(metodo, "→ escrita")
	default:
		fmt.Println("outro")
	}

	// switch sem expressão substitui if / else if
	status := 404
	switch {
	case status >= 500:
		fmt.Println(status, "→ erro do servidor")
	case status >= 400:
		fmt.Println(status, "→ erro do cliente")
	default:
		fmt.Println(status, "→ sucesso")
	}
}

func exemploLoops() {
	// for é o ÚNICO loop do Go

	for i := 0; i < 3; i++ { // clássico
		fmt.Print(i, " ")
	}
	fmt.Println("← clássico")

	n := 1
	for n < 100 { // "while"
		n *= 2
	}
	fmt.Println(n, "← while")

	for i := range 3 { // 0, 1, 2 (Go 1.22+)
		fmt.Print(i, " ")
	}
	fmt.Println("← range 3")

	nomes := []string{"Ana", "Bia", "Caio"} // slice (≈ array do JS), visto no capítulo 02
	for i, nome := range nomes {            // ≈ nomes.forEach((nome, i) => ...)
		fmt.Print(i, ":", nome, " ")
	}
	fmt.Println("← range slice")

	for { // infinito, sai com break
		n--
		if n < 125 {
			break
		}
	}
	fmt.Println(n, "← for + break")
}
