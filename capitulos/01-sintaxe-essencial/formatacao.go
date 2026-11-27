package main

import "fmt"

type usuario struct {
	Nome  string
	Idade int
}

func exemploFmt() {
	u := usuario{Nome: "Ana", Idade: 28}

	fmt.Printf("%%v  → %v\n", u)  // valor padrão (≈ console.log)
	fmt.Printf("%%+v → %+v\n", u) // com nomes dos campos: o melhor para debug
	fmt.Printf("%%T  → %T\n", u)  // tipo
	fmt.Printf("%%d → %d | %%.2f → %.2f | %%s → %s | %%q → %q\n", 42, 3.14159, "go", "go")

	// Sprintf devolve a string em vez de imprimir (≈ template literal)
	msg := fmt.Sprintf("%s tem %d anos", u.Nome, u.Idade)
	fmt.Println(msg)
}
