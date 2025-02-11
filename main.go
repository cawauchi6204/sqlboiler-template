package main

import (
	"log"
	"net/http"

	"todoapp/internal/handler"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// ルーティング定義
	h := handler.NewTodoHandler()
	e.GET("/todos", h.ListTodos)
	e.POST("/todos", h.CreateTodo)
	e.PUT("/todos/:id", h.UpdateTodo)
	e.DELETE("/todos/:id", h.DeleteTodo)

	// サーバ起動
	if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
		log.Fatal("サーバのシャットダウン中にエラーが発生しました: ", err)
	}
}
