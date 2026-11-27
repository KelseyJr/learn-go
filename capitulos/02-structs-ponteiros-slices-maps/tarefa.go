package main

import (
	"fmt"
	"time"
)

// Tarefa é o modelo central da API.
type Tarefa struct {
	ID        int
	Titulo    string
	Concluida bool
	CriadaEm  time.Time
}

// Concluir altera a tarefa, então precisa de receiver ponteiro.
func (t *Tarefa) Concluir() {
	t.Concluida = true
}

// Resumo só lê, mas usa ponteiro também: num mesmo tipo, não se misturam
// receivers por valor e por ponteiro (convenção do Go).
func (t *Tarefa) Resumo() string {
	marca := " "
	if t.Concluida {
		marca = "x"
	}
	return fmt.Sprintf("[%s] #%d %s", marca, t.ID, t.Titulo)
}
