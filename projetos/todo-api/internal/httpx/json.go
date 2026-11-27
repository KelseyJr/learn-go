// Package httpx reúne utilitários HTTP usados por todos os handlers.
package httpx

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
)

// WriteJSON escreve `v` como JSON com o status informado (≈ res.status(s).json(v)).
//
// `any` aceita qualquer tipo (≈ unknown do TS).
func WriteJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		slog.Error("falha ao escrever JSON", "err", err)
	}
}

// WriteError responde no formato {"error": "mensagem"}.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, map[string]string{"error": msg})
}

// DecodeJSON lê o corpo da requisição para dentro de `dst` (≈ JSON.parse(req.body)).
// `dst` precisa ser um ponteiro, para o Decode conseguir preencher a variável de quem chamou.
func DecodeJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // limita o corpo a 1 MB

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // rejeita campos que não existem na struct

	if err := dec.Decode(dst); err != nil {
		return fmt.Errorf("JSON inválido: %w", err)
	}
	return nil
}
