package main

import (
	"fmt"
	"net/http"
	"test_Auth/internal/auth"
	"test_Auth/internal/dbconn"
	"test_Auth/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type TokenMaker struct {
}

func main() {
	db := dbconn.Open()
	dbconn.CreateTable(db)
	tokenMaker := "12345678901"

	hdl := handler.NewHandler(tokenMaker)
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Post("/register", hdl.AuthUser)
	router.With(auth.GetAuthMiddlewareFunc(hdl.TokenMaker)).Post("/refresh", hdl.HdlRefresh)

	if err := http.ListenAndServe(":8082", router); err != nil {
		fmt.Println("failed to start server")
	}
	fmt.Println("server stopped")
}
