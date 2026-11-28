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
type Store interface {
	List(ctx context.Context) ([]Task, error)
	Get(ctx context.Context, id int64) (Task, error)
	Create(ctx context.Context, title string) (Task, error)
	Update(ctx context.Context, t Task) (Task, error)
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
func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /tasks", h.list)
	mux.HandleFunc("POST /tasks", h.create)
	mux.HandleFunc("GET /tasks/{id}", h.get)
	mux.HandleFunc("PATCH /tasks/{id}", h.update)
	mux.HandleFunc("DELETE /tasks/{id}", h.remove)
}

// list responde GET /tasks com todas as tarefas.
//
// Todo handler HTTP tem esta assinatura: w ≈ res e r ≈ req do Express.
func (h *Handler) list(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.store.List(r.Context())
	if err != nil {
		respondError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, tasks)
}

// get responde GET /tasks/{id} com uma tarefa.
func (h *Handler) get(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	t, err := h.store.Get(r.Context(), id)
	if err != nil {
		respondError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, t)
}

// createRequest é o corpo esperado no POST (≈ um DTO).
type createRequest struct {
	Title string `json:"title"`
}

// create responde POST /tasks: valida o título, cria a tarefa e devolve 201.
func (h *Handler) create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil { // &req para o Decode preencher a variável
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	title := strings.TrimSpace(req.Title)
	if err := ValidateTitle(title); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	t, err := h.store.Create(r.Context(), title)
	if err != nil {
		respondError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusCreated, t)
}

// updateRequest usa ponteiros para tornar os campos opcionais (≈ title?: string no TS):
// nil significa "não enviado". Sem ponteiro, `done` ausente viraria false
// (o zero value) e não daria para diferenciar de "done": false.
type updateRequest struct {
	Title *string `json:"title"`
	Done  *bool   `json:"done"`
}

// update responde PATCH /tasks/{id}: altera só os campos enviados no corpo.
func (h *Handler) update(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req updateRequest
	if err := httpx.DecodeJSON(w, r, &req); err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	t, err := h.store.Get(r.Context(), id)
	if err != nil {
		respondError(w, err)
		return
	}

	// Aplica só os campos enviados. *req.Title lê o valor para onde o ponteiro aponta.
	if req.Title != nil {
		title := strings.TrimSpace(*req.Title)
		if err := ValidateTitle(title); err != nil {
			httpx.WriteError(w, http.StatusBadRequest, err.Error())
			return
		}
		t.Title = title
	}
	if req.Done != nil {
		t.Done = *req.Done
	}

	t, err = h.store.Update(r.Context(), t)
	if err != nil {
		respondError(w, err)
		return
	}
	httpx.WriteJSON(w, http.StatusOK, t)
}

// remove responde DELETE /tasks/{id} com 204 (sucesso, sem corpo).
func (h *Handler) remove(w http.ResponseWriter, r *http.Request) {
	id, err := parseID(r)
	if err != nil {
		httpx.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.store.Delete(r.Context(), id); err != nil {
		respondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// parseID lê o {id} da URL (≈ req.params.id) e converte para número.
func parseID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("id inválido")
	}
	return id, nil
}

// respondError traduz erros da camada de dados em status HTTP.
func respondError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrNotFound):
		httpx.WriteError(w, http.StatusNotFound, err.Error())
	default:
		// Erro inesperado: registra o detalhe no log e devolve uma mensagem genérica
		slog.Error("erro inesperado", "err", err)
		httpx.WriteError(w, http.StatusInternalServerError, "erro interno")
	}
}
