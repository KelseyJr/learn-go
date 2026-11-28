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
func (s *MemoryStore) List(_ context.Context) ([]Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// make garante uma lista vazia ([]) em vez de nil, que viraria `null` no JSON
	tasks := make([]Task, 0, len(s.tasks))
	for _, t := range s.tasks {
		tasks = append(tasks, t)
	}
	// A ordem de um map é aleatória, então ordenamos por ID
	slices.SortFunc(tasks, func(a, b Task) int { return cmp.Compare(a.ID, b.ID) })
	return tasks, nil
}

// Get devolve a tarefa com o ID informado, ou ErrNotFound se ela não existir.
func (s *MemoryStore) Get(_ context.Context, id int64) (Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tasks[id]
	if !ok {
		return Task{}, ErrNotFound
	}
	return t, nil
}

// Create salva uma nova tarefa com o próximo ID disponível e a devolve.
func (s *MemoryStore) Create(_ context.Context, title string) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := Task{
		ID:        s.nextID,
		Title:     title,
		CreatedAt: time.Now().UTC(),
	}
	s.tasks[t.ID] = t
	s.nextID++
	return t, nil
}

// Update substitui a tarefa salva pela versão recebida (identificada pelo ID).
// Devolve ErrNotFound se a tarefa não existir.
func (s *MemoryStore) Update(_ context.Context, t Task) (Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[t.ID]; !ok {
		return Task{}, ErrNotFound
	}
	s.tasks[t.ID] = t
	return t, nil
}

// Delete remove a tarefa com o ID informado, ou devolve ErrNotFound se ela não existir.
func (s *MemoryStore) Delete(_ context.Context, id int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[id]; !ok {
		return ErrNotFound
	}
	delete(s.tasks, id)
	return nil
}
