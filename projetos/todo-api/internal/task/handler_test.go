package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// newServer monta as rotas com um store novo, isolado para cada teste.
func newServer() http.Handler {
	mux := http.NewServeMux()
	NewHandler(NewMemoryStore()).RegisterRoutes(mux)
	return mux
}

// do executa uma requisição direto no roteador, sem abrir porta de rede (≈ supertest).
func do(t *testing.T, srv http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

// decode é uma função genérica (≈ function decode<T>() do TS): decode[Task](...) devolve uma Task.
func decode[T any](t *testing.T, rec *httptest.ResponseRecorder) T {
	t.Helper()
	var v T
	if err := json.Unmarshal(rec.Body.Bytes(), &v); err != nil {
		t.Fatalf("resposta não é JSON válido: %v (corpo: %s)", err, rec.Body.String())
	}
	return v
}

func TestCreateAndGet(t *testing.T) {
	srv := newServer()

	rec := do(t, srv, "POST", "/tasks", `{"title": "  Estudar Go  "}`)
	if rec.Code != http.StatusCreated {
		t.Fatalf("POST /tasks: status %d, esperado 201 (corpo: %s)", rec.Code, rec.Body.String())
	}
	created := decode[Task](t, rec)
	if created.ID != 1 || created.Title != "Estudar Go" || created.Done {
		t.Fatalf("tarefa criada inesperada: %+v", created)
	}

	rec = do(t, srv, "GET", "/tasks/1", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("GET /tasks/1: status %d, esperado 200", rec.Code)
	}
	if got := decode[Task](t, rec); got.Title != "Estudar Go" {
		t.Errorf("GET /tasks/1: título %q, esperado %q", got.Title, "Estudar Go")
	}
}

func TestListStartsAsEmptyArray(t *testing.T) {
	rec := do(t, newServer(), "GET", "/tasks", "")

	if body := strings.TrimSpace(rec.Body.String()); body != "[]" {
		t.Errorf("GET /tasks sem tarefas: corpo %s, esperado []", body)
	}
}

func TestPatchUpdatesOnlySentFields(t *testing.T) {
	srv := newServer()
	do(t, srv, "POST", "/tasks", `{"title": "Estudar Go"}`)

	rec := do(t, srv, "PATCH", "/tasks/1", `{"done": true}`)
	if rec.Code != http.StatusOK {
		t.Fatalf("PATCH /tasks/1: status %d, esperado 200 (corpo: %s)", rec.Code, rec.Body.String())
	}
	got := decode[Task](t, rec)
	if !got.Done || got.Title != "Estudar Go" {
		t.Errorf("PATCH só com done deveria manter o título: %+v", got)
	}
}

func TestDelete(t *testing.T) {
	srv := newServer()
	do(t, srv, "POST", "/tasks", `{"title": "Apagar"}`)

	if rec := do(t, srv, "DELETE", "/tasks/1", ""); rec.Code != http.StatusNoContent {
		t.Fatalf("DELETE /tasks/1: status %d, esperado 204", rec.Code)
	}
	if rec := do(t, srv, "GET", "/tasks/1", ""); rec.Code != http.StatusNotFound {
		t.Errorf("GET depois do DELETE: status %d, esperado 404", rec.Code)
	}
}

// Teste "table-driven": uma tabela de casos e um loop. É o estilo mais comum em Go.
func TestErrors(t *testing.T) {
	tooLong := `{"title": "` + strings.Repeat("a", 201) + `"}`

	cases := []struct {
		name         string
		method, path string
		body         string
		wantStatus   int
	}{
		{"título vazio", "POST", "/tasks", `{"title": ""}`, http.StatusBadRequest},
		{"título só com espaços", "POST", "/tasks", `{"title": "   "}`, http.StatusBadRequest},
		{"título longo demais", "POST", "/tasks", tooLong, http.StatusBadRequest},
		{"JSON malformado", "POST", "/tasks", `{"title":`, http.StatusBadRequest},
		{"campo desconhecido", "POST", "/tasks", `{"titulo": "x"}`, http.StatusBadRequest},
		{"id não numérico", "GET", "/tasks/abc", "", http.StatusBadRequest},
		{"tarefa inexistente", "GET", "/tasks/99", "", http.StatusNotFound},
		{"PATCH em inexistente", "PATCH", "/tasks/99", `{"done": true}`, http.StatusNotFound},
		{"DELETE em inexistente", "DELETE", "/tasks/99", "", http.StatusNotFound},
		{"método não permitido", "PUT", "/tasks/1", "", http.StatusMethodNotAllowed},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			rec := do(t, newServer(), c.method, c.path, c.body)
			if rec.Code != c.wantStatus {
				t.Errorf("%s %s: status %d, esperado %d (corpo: %s)",
					c.method, c.path, rec.Code, c.wantStatus, rec.Body.String())
			}
		})
	}
}
