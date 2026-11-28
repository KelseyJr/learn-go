package task

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"todo/internal/httpx"
)

// Store é tudo o que o Handler precisa da camada de dados.
//
// Em Go não existe "implements": qualquer tipo que tenha esses métodos já
// satisfaz a interface. O MemoryStore satisfaz hoje; um PostgresStore com os
// mesmos métodos satisfaria amanhã, sem mudar uma linha deste arquivo.
//
// Dois nomes continuam curtos de propósito, porque o nome longo daria conflito:
//   - ctx: chamá-lo de `context` esconderia o pacote context dentro da função
//   - err: chamá-lo de `error` esconderia o tipo embutido error
type Store interface {
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, id int64) (Task, error)
	Create(ctx context.Context, title string) (Task, error)
	Update(ctx context.Context, task Task) (Task, error)
	Delete(ctx context.Context, id int64) error
}

// Handler é a camada HTTP das tarefas (≈ controller no Express).
type Handler struct {
	store Store
}

// NewHandler cria um Handler que usa o store informado para ler e gravar tarefas.
func NewHandler(store Store) *Handler {
	return &Handler{store: store}
}

// RegisterRoutes associa cada rota a um método (≈ router.get('/tasks', ...)).
// Os métodos dos handlers começam com letra minúscula, então são privados ao
// pacote: de fora, só RegisterRoutes é visível.
func (handler *Handler) RegisterRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /tasks", handler.list)
	router.HandleFunc("POST /tasks", handler.create)
	router.HandleFunc("GET /tasks/{id}", handler.get)
	router.HandleFunc("PATCH /tasks/{id}", handler.update)
	router.HandleFunc("DELETE /tasks/{id}", handler.remove)
}

// list responde GET /tasks com todas as tarefas.
//
// Todo handler HTTP tem esta assinatura: response ≈ res e request ≈ req do Express.
func (handler *Handler) list(response http.ResponseWriter, request *http.Request) {
	tasks, err := handler.store.List(request.Context())
	if err != nil {
		respondError(response, err)
		return
	}
	httpx.WriteJSON(response, http.StatusOK, tasks)
}

// get responde GET /tasks/{id} com uma tarefa.
func (handler *Handler) get(response http.ResponseWriter, request *http.Request) {
	id, err := parseID(request)
	if err != nil {
		httpx.WriteError(response, http.StatusBadRequest, err.Error())
		return
	}

	task, err := handler.store.Get(request.Context(), id)
	if err != nil {
		respondError(response, err)
		return
	}
	httpx.WriteJSON(response, http.StatusOK, task)
}

// createRequest é o corpo esperado no POST (≈ um DTO).
type createRequest struct {
	Title string `json:"title"`
}

// create responde POST /tasks: valida o título, cria a tarefa e devolve 201.
func (handler *Handler) create(response http.ResponseWriter, request *http.Request) {
	var body createRequest
	if err := httpx.DecodeJSON(response, request, &body); err != nil { // &body para o Decode preencher a variável
		httpx.WriteError(response, http.StatusBadRequest, err.Error())
		return
	}

	title := strings.TrimSpace(body.Title)
	if err := ValidateTitle(title); err != nil {
		httpx.WriteError(response, http.StatusBadRequest, err.Error())
		return
	}

	task, err := handler.store.Create(request.Context(), title)
	if err != nil {
		respondError(response, err)
		return
	}
	httpx.WriteJSON(response, http.StatusCreated, task)
}

// updateRequest usa ponteiros para tornar os campos opcionais (≈ title?: string no TS):
// nil significa "não enviado". Sem ponteiro, `done` ausente viraria false
// (o zero value) e não daria para diferenciar de "done": false.
type updateRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

// update responde PATCH /tasks/{id}: altera só os campos enviados no corpo.
func (handler *Handler) update(response http.ResponseWriter, request *http.Request) {
	id, err := parseID(request)
	if err != nil {
		httpx.WriteError(response, http.StatusBadRequest, err.Error())
		return
	}

	var body updateRequest
	if err := httpx.DecodeJSON(response, request, &body); err != nil {
		httpx.WriteError(response, http.StatusBadRequest, err.Error())
		return
	}

	task, err := handler.store.Get(request.Context(), id)
	if err != nil {
		respondError(response, err)
		return
	}

	// Aplica só os campos enviados. *body.Title lê o valor para onde o ponteiro aponta.
	if body.Title != nil {
		title := strings.TrimSpace(*body.Title)
		if err := ValidateTitle(title); err != nil {
			httpx.WriteError(response, http.StatusBadRequest, err.Error())
			return
		}
		task.Title = title
	}
	if body.Done != nil {
		task.Done = *body.Done
	}

	task, err = handler.store.Update(request.Context(), task)
	if err != nil {
		respondError(response, err)
		return
	}
	httpx.WriteJSON(response, http.StatusOK, task)
}

// remove responde DELETE /tasks/{id} com 204 (sucesso, sem corpo).
func (handler *Handler) remove(response http.ResponseWriter, request *http.Request) {
	id, err := parseID(request)
	if err != nil {
		httpx.WriteError(response, http.StatusBadRequest, err.Error())
		return
	}

	if err := handler.store.Delete(request.Context(), id); err != nil {
		respondError(response, err)
		return
	}
	response.WriteHeader(http.StatusNoContent)
}

// parseID lê o {id} da URL (≈ req.params.id) e converte para número.
func parseID(request *http.Request) (int64, error) {
	id, err := strconv.ParseInt(request.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("id inválido")
	}
	return id, nil
}

// respondError traduz erros da camada de dados em status HTTP.
func respondError(response http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(response, http.StatusNotFound, err.Error())
	default:
		// Erro inesperado: registra o detalhe no log e devolve uma mensagem genérica
		slog.Error("erro inesperado", "err", err)
		httpx.WriteError(response, http.StatusInternalServerError, "erro interno")
	}
}
