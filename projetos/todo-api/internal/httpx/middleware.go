package httpx

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger é um middleware que registra método, caminho, status e duração de
// cada requisição (≈ morgan no Express).
//
// Middleware em Go é só uma função que recebe um http.Handler e devolve outro,
// que "embrulha" o primeiro: faz algo antes, chama o próximo, faz algo depois.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r) // ≈ next() no Express

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start),
		)
	})
}

// statusRecorder guarda o status code que o handler escreveu, para o log.
//
// Embedding: ao colocar http.ResponseWriter sem nome de campo, a struct
// "herda" todos os métodos dele. Só sobrescrevemos o WriteHeader.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}
