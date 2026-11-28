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
// Neste projeto os receivers têm nome completo (store, handler) para facilitar
// a leitura. O costume em Go é usar uma ou duas letras: (s *MemoryStore).
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
	for _, t := range store.tasks {
		tasks = append(tasks, t)
	}
	// A ordem de um map é aleatória, então ordenamos por ID
	slices.SortFunc(tasks, func(a, b Task) int { return cmp.Compare(a.ID, b.ID) })
	return tasks, nil
}

// Get devolve a tarefa com o ID informado, ou ErrNotFound se ela não existir.
func (store *MemoryStore) Get(_ context.Context, id int64) (Task, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()

	t, ok := store.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return t, nil
}

// Create salva uma nova tarefa com o próximo ID disponível e a devolve.
func (store *MemoryStore) Create(_ context.Context, title string) (Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	t := Task{
		ID:        store.nextID,
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
	store.tasks[t.ID] = t
	store.nextID++
	return t, nil
}

// Update substitui a tarefa salva pela versão recebida (identificada pelo ID).
// Devolve ErrNotFound se a tarefa não existir.
func (store *MemoryStore) Update(_ context.Context, t Task) (Task, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.tasks[t.ID]; !ok {
		return Task{}, ErrNotFound
	}
	store.tasks[t.ID] = t
	return t, nil
}

// Delete remove a tarefa com o ID informado, ou devolve ErrNotFound se ela não existir.
func (store *MemoryStore) Delete(_ context.Context, id int64) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	if _, ok := store.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(store.tasks, id)
	return nil
}
