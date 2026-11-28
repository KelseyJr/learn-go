package task

import (
	"cmp"
	"context"
	"slices"
	"sync"
	"time"
)

// MemoryStore guarda as tarefas em memória: tudo some quando o servidor reinicia.
//
// Cada requisição HTTP roda na sua própria goroutine, então várias requisições
// podem mexer no map AO MESMO TEMPO. O mutex garante a segurança:
//   - RLock: vários leitores ao mesmo tempo
//   - Lock: um escritor por vez, sem ninguém lendo
//
// No Node isso não existe, porque o event loop executa um callback de cada vez.
//
// Os métodos recebem context.Context para ter a mesma assinatura de um store com
// banco de dados (que repassaria o ctx para a query). Aqui ele não é usado, por isso o `_`.
//
// Neste projeto receivers, parâmetros e variáveis têm nomes por extenso (store,
// task, response, request) para facilitar a leitura. O costume em Go é usar nomes
// curtos quando o contexto é óbvio: (s *MemoryStore), t, w, r.
type MemoryStore struct {
	mu     sync.RWMutex
	tasks  map[int64]Task
	nextID int64
}

// NewMemoryStore cria um MemoryStore vazio, pronto para uso.
// O map precisa do make: escrever num map nil causa panic.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tasks:  make(map[int64]Task),
		nextID: 1,
	}
}

// List devolve todas as tarefas, ordenadas por ID.
func (store *MemoryStore) List(_ context.Context) ([]Task, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	// make garante uma lista vazia ([]) em vez de nil, que viraria `null` no JSON
	tasks := make([]Task, 0, len(store.tasks))
	for _, task := range store.tasks {
		tasks = append(tasks, task)
	}
	// A ordem de um map é aleatória, então ordenamos por ID (≈ sort((a, b) => a.id - b.id))
	slices.SortFunc(tasks, func(first, second Task) int {
		return cmp.Compare(first.ID, second.ID)
	})
	return tasks, nil
}

// Get devolve a tarefa com o ID informado, ou ErrNotFound se ela não existir.
func (store *MemoryStore) Get(_ context.Context, id int64) (Task, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	task, exists := store.tasks[id]
	if !exists {
		return Task{}, ErrNotFound
	}
	return task, nil
}

// Create salva uma nova tarefa com o próximo ID disponível e a devolve.
func (store *MemoryStore) Create(_ context.Context, title string) (Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	task := Task{
		ID:        store.nextID,
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
	store.tasks[task.ID] = task
	store.nextID++
	return task, nil
}

// Update substitui a tarefa salva pela versão recebida (identificada pelo ID).
// Devolve ErrNotFound se a tarefa não existir.
func (store *MemoryStore) Update(_ context.Context, task Task) (Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.tasks[task.ID]; !exists {
		return Task{}, ErrNotFound
	}
	store.tasks[task.ID] = task
	return task, nil
}

// Delete remove a tarefa com o ID informado, ou devolve ErrNotFound se ela não existir.
func (store *MemoryStore) Delete(_ context.Context, id int64) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, exists := store.tasks[id]; !exists {
		return ErrNotFound
	}
	delete(store.tasks, id)
	return nil
}
