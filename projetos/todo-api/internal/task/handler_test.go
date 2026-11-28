package task

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Toda função de teste começa com "Test" e recebe *testing.T, que controla o
// teste (marcar falha, criar subtestes...). O costume em Go é chamá-lo de `t`;
// aqui ele se chama `test` para facilitar a leitura.

// newServer monta as rotas com um store novo, isolado para cada teste.
func newServer() http.Handler {
	router := http.NewServeMux()
	NewHandler(NewMemoryStore()).RegisterRoutes(router)
	return router
}

// send executa uma requisição direto no roteador, sem abrir porta de rede (≈ supertest).
func send(test *testing.T, server http.Handler, method, path, body string) *httptest.ResponseRecorder {
	test.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	response := httptest.NewRecorder()
	server.ServeHTTP(response, request)
	return response
}

// decode é uma função genérica (≈ function decode<T>() do TS): decode[Task](...) devolve uma Task.
func decode[T any](test *testing.T, response *httptest.ResponseRecorder) T {
	test.Helper()
	var value T
	if err := json.Unmarshal(response.Body.Bytes(), &value); err != nil {
		test.Fatalf("resposta não é JSON válido: %v (corpo: %s)", err, response.Body.String())
	}
	return value
}

func TestCreateAndGet(test *testing.T) {
	server := newServer()

	response := send(test, server, "POST", "/tasks", `{"title": "  Estudar Go  "}`)
	if response.Code != http.StatusCreated {
		test.Fatalf("POST /tasks: status %d, esperado 201 (corpo: %s)", response.Code, response.Body.String())
	}
	created := decode[Task](test, response)
	if created.ID != 1 || created.Title != "Estudar Go" || created.Done {
		test.Fatalf("tarefa criada inesperada: %+v", created)
	}

	response = send(test, server, "GET", "/tasks/1", "")
	if response.Code != http.StatusOK {
		test.Fatalf("GET /tasks/1: status %d, esperado 200", response.Code)
	}
	if found := decode[Task](test, response); found.Title != "Estudar Go" {
		test.Errorf("GET /tasks/1: título %q, esperado %q", found.Title, "Estudar Go")
	}
}

func TestListStartsAsEmptyArray(test *testing.T) {
	response := send(test, newServer(), "GET", "/tasks", "")

	if body := strings.TrimSpace(response.Body.String()); body != "[]" {
		test.Errorf("GET /tasks sem tarefas: corpo %s, esperado []", body)
	}
}

func TestPatchUpdatesOnlySentFields(test *testing.T) {
	server := newServer()
	send(test, server, "POST", "/tasks", `{"title": "Estudar Go"}`)

	response := send(test, server, "PATCH", "/tasks/1", `{"done": true}`)
	if response.Code != http.StatusOK {
		test.Fatalf("PATCH /tasks/1: status %d, esperado 200 (corpo: %s)", response.Code, response.Body.String())
	}
	updated := decode[Task](test, response)
	if !updated.Done || updated.Title != "Estudar Go" {
		test.Errorf("PATCH só com done deveria manter o título: %+v", updated)
	}
}

func TestDelete(test *testing.T) {
	server := newServer()
	send(test, server, "POST", "/tasks", `{"title": "Apagar"}`)

	if response := send(test, server, "DELETE", "/tasks/1", ""); response.Code != http.StatusNoContent {
		test.Fatalf("DELETE /tasks/1: status %d, esperado 204", response.Code)
	}
	if response := send(test, server, "GET", "/tasks/1", ""); response.Code != http.StatusNotFound {
		test.Errorf("GET depois do DELETE: status %d, esperado 404", response.Code)
	}
}

// Teste "table-driven": uma tabela de casos e um loop. É o estilo mais comum em Go.
func TestErrors(test *testing.T) {
	tooLong := `{"title": "` + strings.Repeat("a", 201) + `"}`

	testCases := []struct {
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

	for _, testCase := range testCases {
		test.Run(testCase.name, func(test *testing.T) {
			response := send(test, newServer(), testCase.method, testCase.path, testCase.body)
			if response.Code != testCase.wantStatus {
				test.Errorf("%s %s: status %d, esperado %d (corpo: %s)",
					testCase.method, testCase.path, response.Code, testCase.wantStatus, response.Body.String())
			}
		})
	}
}
