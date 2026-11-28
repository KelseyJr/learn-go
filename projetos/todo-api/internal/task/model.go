// Package task contém tudo sobre tarefas: o modelo, o armazenamento e as rotas HTTP.
//
// Organizar por funcionalidade (task, user, billing...) em vez de por camada
// (models/, controllers/, services/) é o mais comum em projetos Go.
package task

import (
	"errors"
	"time"
	"unicode/utf8"
)

// Task é o modelo da API.
//
// As struct tags `json:"..."` definem o nome de cada campo no JSON.
// Sem elas, o JSON sairia com os nomes do Go: "ID", "Title", "CreatedAt".
type Task struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Done      bool      `json:"done"`
	CreatedAt time.Time `json:"created_at"`
}

// Erros "sentinela": valores fixos que quem chama compara com errors.Is
// para decidir o que fazer (por exemplo, responder 404 em vez de 500).
var (
	ErrNotFound      = errors.New("tarefa não encontrada")
	ErrTitleRequired = errors.New("o título é obrigatório")
	ErrTitleTooLong  = errors.New("o título deve ter no máximo 200 caracteres")
)

const maxTitleLength = 200

// ValidateTitle aplica as regras de negócio do título.
func ValidateTitle(title string) error {
	if title == "" {
		return ErrTitleRequired
	}
	// Conta caracteres, não bytes: len("ção") é 5, mas são 3 caracteres
	if utf8.RuneCountInString(title) > maxTitleLength {
		return ErrTitleTooLong
	}
	return nil
}
