package usecase

import (
	"errors"
	"todoapp/internal/model"
	"todoapp/internal/repository"
)

type TodoUsecase interface {
	ListTodos() ([]model.Todo, error)
	CreateTodo(title string) (model.Todo, error)
	UpdateTodo(id int, title, status string) (model.Todo, error)
	DeleteTodo(id int) error
}

type todoUsecase struct {
	repo repository.TodoRepository
}

func NewTodoUsecase() TodoUsecase {
	return &todoUsecase{
		repo: repository.NewTodoRepository(),
	}
}

func (u *todoUsecase) ListTodos() ([]model.Todo, error) {
	return u.repo.GetAll()
}

func (u *todoUsecase) CreateTodo(title string) (model.Todo, error) {
	todo := model.Todo{Title: title, Status: "pending"}
	return u.repo.Create(todo)
}

func (u *todoUsecase) UpdateTodo(id int, title, status string) (model.Todo, error) {
	todo, err := u.repo.GetByID(id)
	if err != nil {
		return model.Todo{}, err
	}
	todo.Title = title
	todo.Status = status
	return u.repo.Update(todo)
}

func (u *todoUsecase) DeleteTodo(id int) error {
	_, err := u.repo.GetByID(id)
	if err != nil {
		return errors.New("todo not found")
	}
	return u.repo.Delete(id)
}
