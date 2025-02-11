package repository

import (
	"errors"
	"sync"
	"todoapp/internal/model"
)

type TodoRepository interface {
	GetAll() ([]model.Todo, error)
	GetByID(id int) (model.Todo, error)
	Create(todo model.Todo) (model.Todo, error)
	Update(todo model.Todo) (model.Todo, error)
	Delete(id int) error
}

type todoRepository struct {
	mu     sync.Mutex
	todos  map[int]model.Todo
	nextID int
}

func NewTodoRepository() TodoRepository {
	return &todoRepository{
		todos:  make(map[int]model.Todo),
		nextID: 1,
	}
}

func (r *todoRepository) GetAll() ([]model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	list := make([]model.Todo, 0, len(r.todos))
	for _, todo := range r.todos {
		list = append(list, todo)
	}
	return list, nil
}

func (r *todoRepository) GetByID(id int) (model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	todo, ok := r.todos[id]
	if !ok {
		return model.Todo{}, errors.New("todo not found")
	}
	return todo, nil
}

func (r *todoRepository) Create(todo model.Todo) (model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	todo.ID = r.nextID
	r.todos[r.nextID] = todo
	r.nextID++
	return todo, nil
}

func (r *todoRepository) Update(todo model.Todo) (model.Todo, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.todos[todo.ID]; !exists {
		return model.Todo{}, errors.New("todo not found")
	}
	r.todos[todo.ID] = todo
	return todo, nil
}

func (r *todoRepository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.todos[id]; !exists {
		return errors.New("todo not found")
	}
	delete(r.todos, id)
	return nil
}
