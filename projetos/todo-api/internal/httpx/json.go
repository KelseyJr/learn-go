// Package httpx reúne utilitários HTTP usados por todos os handlers.
package httpx

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// WriteJSON escreve `value` como JSON com o status informado (≈ res.status(s).json(value)).
//
// `any` aceita qualquer tipo (≈ unknown do TS).
func WriteJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	if err := json.NewEncoder(response).Encode(value); err != nil {
		slog.Error("falha ao escrever JSON", "err", err)
	}
}

// WriteError responde no formato {"error": "mensagem"}.
func WriteError(response http.ResponseWriter, status int, message string) {
	WriteJSON(response, status, map[string]string{"error": message})
}

// DecodeJSON lê o corpo da requisição para dentro de `destination` (≈ JSON.parse(req.body)).
// `destination` precisa ser um ponteiro, para o Decode conseguir preencher a variável de quem chamou.
func DecodeJSON(response http.ResponseWriter, request *http.Request, destination any) error {
	request.Body = http.MaxBytesReader(response, request.Body, 1<<20) // limita o corpo a 1 MB

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields() // rejeita campos que não existem na struct

	if err := decoder.Decode(destination); err != nil {
		return fmt.Errorf("JSON inválido: %w", err)
	}
	return nil
}
