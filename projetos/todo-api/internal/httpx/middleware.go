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
	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		start := time.Now()
		recorder := &statusRecorder{ResponseWriter: response, status: http.StatusOK}

		next.ServeHTTP(recorder, request) // ≈ next() no Express

		slog.Info("request",
			"method", request.Method,
			"path", request.URL.Path,
			"status", recorder.status,
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

// WriteHeader anota o status antes de repassá-lo ao ResponseWriter original.
func (recorder *statusRecorder) WriteHeader(statusCode int) {
	recorder.status = statusCode
	recorder.ResponseWriter.WriteHeader(statusCode)
}
