// Capítulo 02 — Structs, ponteiros, slices e maps
//
// Rode com: make run CAP=02
//
//	fundamentos.go → exemplos de cada conceito
//	tarefa.go      → modelo Tarefa (início da API)
//	repo.go        → Repo: armazenamento em memória com CRUD
package main

import "fmt"

func main() {
	secao("Structs", exemploStructs)
	secao("Ponteiros", exemploPonteiros)
	secao("Métodos", exemploMetodos)
	secao("Slices", exemploSlices)
	secao("Maps", exemploMaps)
	secao("API: Repo de tarefas", exemploRepo)
}

func secao(titulo string, exemplo func()) {
	fmt.Printf("\n━━━ %s ━━━\n", titulo)
	exemplo()
}

func exemploRepo() {
	repo := NovoRepo()

	repo.Criar("Estudar Go")
	repo.Criar("Criar a API")
	repo.Criar("Fazer deploy")

	// Buscar devolve uma CÓPIA: altera a cópia e salva de volta
	if t, ok := repo.Buscar(2); ok {
		t.Concluir()
		repo.Atualizar(t)
	}

	repo.Remover(3)

	if _, ok := repo.Buscar(99); !ok {
		fmt.Println("tarefa 99 não encontrada")
	}

	for _, t := range repo.Listar() {
		fmt.Println(t.Resumo())
	}
}
