package main

import (
	"slices"
	"time"
)

// Repo guarda as tarefas em memória. No capítulo 05 ele passa a usar o Postgres.
type Repo struct {
	tarefas   map[int]Tarefa // id → tarefa (≈ new Map<number, Tarefa>())
	proximoID int
}

// NovoRepo é o "construtor". Ele precisa inicializar o map:
// um map nil pode ser lido, mas escrever nele causa panic.
func NovoRepo() *Repo {
	return &Repo{
		tarefas:   make(map[int]Tarefa),
		proximoID: 1,
	}
}

func (r *Repo) Criar(titulo string) Tarefa {
	t := Tarefa{
		ID:       r.proximoID,
		Titulo:   titulo,
		CriadaEm: time.Now(),
	}
	r.tarefas[t.ID] = t
	r.proximoID++
	return t
}

// Buscar segue o idioma "comma ok": (valor, existe).
// No capítulo 04 isso vira (Tarefa, error).
func (r *Repo) Buscar(id int) (Tarefa, bool) {
	t, ok := r.tarefas[id]
	return t, ok
}

// Listar devolve as tarefas ordenadas por ID, já que a ordem de um map é aleatória.
func (r *Repo) Listar() []Tarefa {
	lista := make([]Tarefa, 0, len(r.tarefas))
	for _, t := range r.tarefas {
		lista = append(lista, t)
	}
	slices.SortFunc(lista, func(a, b Tarefa) int { return a.ID - b.ID }) // ≈ sort((a, b) => a.id - b.id)
	return lista
}

// Atualizar substitui a tarefa salva. Devolve false se o ID não existir.
func (r *Repo) Atualizar(t Tarefa) bool {
	if _, ok := r.tarefas[t.ID]; !ok {
		return false
	}
	r.tarefas[t.ID] = t
	return true
}

// Remover apaga a tarefa. Devolve false se o ID não existir.
func (r *Repo) Remover(id int) bool {
	if _, ok := r.tarefas[id]; !ok {
		return false
	}
	delete(r.tarefas, id)
	return true
}
