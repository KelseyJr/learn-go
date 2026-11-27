package main

import (
	"errors"
	"fmt"
)

// Parâmetros do mesmo tipo podem ser agrupados: (a, b int)
func soma(a, b int) int {
	return a + b
}

// O padrão do Go: devolver (resultado, error) em vez de lançar exceção.
func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divisão por zero")
	}
	return a / b, nil // nil = sem erro
}

// Variádica (≈ rest parameters). Dentro da função, nums é um []int.
func somaTodos(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

func exemploFuncoes() {
	fmt.Println("soma(2, 3) =", soma(2, 3))

	// Este bloco é o que você mais vai ver em código Go:
	resultado, err := dividir(10, 0)
	if err != nil {
		fmt.Println("erro:", err)
	} else {
		fmt.Println("resultado:", resultado)
	}

	numeros := []int{10, 20, 30}
	fmt.Println("somaTodos(numeros...) =", somaTodos(numeros...)) // ≈ spread

	// Função anônima numa variável (≈ arrow function)
	dobrar := func(n int) int { return n * 2 }
	fmt.Println("dobrar(21) =", dobrar(21))
}

// Closure: igual ao JS. É a base dos middlewares HTTP (capítulo 06).
func gerarID(prefixo string) func() string {
	proximo := 0
	return func() string {
		proximo++
		return fmt.Sprintf("%s-%03d", prefixo, proximo)
	}
}

func exemploClosures() {
	novoPedido := gerarID("PED")
	fmt.Println(novoPedido(), novoPedido(), novoPedido())
}

// defer: executa quando a função terminar (≈ finally).
// Uso real: defer arquivo.Close(), defer rows.Close(), defer mutex.Unlock()
func exemploDefer() {
	fmt.Println("abre conexão")
	defer fmt.Println("fecha conexão (defer)")

	fmt.Println("usa conexão")
}
