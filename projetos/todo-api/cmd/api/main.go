// Ponto de entrada da API. Em projetos Go, cada executável fica em cmd/<nome>/main.go,
// e o main só monta as peças: configuração, dependências e servidor.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"todo/internal/httpx"
	"todo/internal/task"
)

func main() {
	// Logger estruturado da biblioteca padrão (≈ pino/winston)
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	if err := run(); err != nil {
		slog.Error("servidor encerrado com erro", "err", err)
		os.Exit(1)
	}
}

// run concentra a inicialização e devolve error em vez de derrubar o processo
// em vários pontos. É um padrão comum no main de projetos Go.
func run() error {
	addr := ":" + envOr("PORT", "8080") // ≈ process.env.PORT ?? "8080"

	// Injeção de dependência manual, sem framework: cria as peças e passa adiante
	store := task.NewMemoryStore()
	handler := task.NewHandler(store)

	mux := http.NewServeMux() // o roteador da biblioteca padrão (≈ express.Router())
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})
	handler.RegisterRoutes(mux)

	srv := &http.Server{
		Addr:              addr,
		Handler:           httpx.Logger(mux), // todas as rotas passam pelo middleware
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	// ─── Graceful shutdown ───────────────────────────────────────────
	// Este bloco é praticamente igual em todo projeto Go. Não precisa
	// dominar goroutines e channels agora para entender a ideia.

	// ctx é cancelado quando o processo recebe Ctrl+C (SIGINT) ou `docker stop` (SIGTERM)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ListenAndServe bloqueia, por isso roda numa goroutine (em "segundo plano")
	serverErr := make(chan error, 1)
	go func() {
		slog.Info("servidor iniciado", "addr", addr)
		serverErr <- srv.ListenAndServe()
	}()

	// Espera o que acontecer primeiro: o servidor falhar ou chegar um sinal de parada
	select {
	case err := <-serverErr:
		if !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
	}

	// Para de aceitar conexões e espera as requisições em andamento terminarem (até 10s)
	slog.Info("encerrando servidor...")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return srv.Shutdown(shutdownCtx)
}

// envOr lê uma variável de ambiente, ou devolve o fallback se ela estiver vazia.
func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
