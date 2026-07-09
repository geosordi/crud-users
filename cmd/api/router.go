package main

import (
	"net/http"

	"crud-users/internal/service"
)

func newRouter(userSvc *service.UserService) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	registerSwaggerRoutes(mux)
	registerUserRoutes(mux, userSvc)

	return mux
}
