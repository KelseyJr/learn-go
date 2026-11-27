package exercicios

import "testing"

// Em Go, testes ficam em arquivos *_test.go e cada função de teste
// começa com "Test" e recebe *testing.T. Veremos isso a fundo no capítulo 14.
func TestSaudacao(t *testing.T) {
	esperado := "Olá, Go!"

	if obtido := Saudacao(); obtido != esperado {
		t.Errorf("Saudacao() retornou %q, mas o esperado era %q", obtido, esperado)
	}
}
