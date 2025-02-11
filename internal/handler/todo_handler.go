package handler

import (
	"net/http"
	"strconv"

	"todoapp/internal/usecase"

	"github.com/labstack/echo/v4"
)

type TodoHandler struct {
	Usecase usecase.TodoUsecase
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{
		Usecase: usecase.NewTodoUsecase(),
	}
}

func (h *TodoHandler) ListTodos(c echo.Context) error {
	todos, err := h.Usecase.ListTodos()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, todos)
}

func (h *TodoHandler) CreateTodo(c echo.Context) error {
	var req struct {
		Title string `json:"title"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	todo, err := h.Usecase.CreateTodo(req.Title)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, todo)
}

func (h *TodoHandler) UpdateTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	var req struct {
		Title  string `json:"title"`
		Status string `json:"status"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	todo, err := h.Usecase.UpdateTodo(id, req.Title, req.Status)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, todo)
}

func (h *TodoHandler) DeleteTodo(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID"})
	}
	err = h.Usecase.DeleteTodo(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusNoContent, nil)
}
