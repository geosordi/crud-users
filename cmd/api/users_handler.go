package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"crud-users/internal/service"
)

func registerUserRoutes(mux *http.ServeMux, userSvc *service.UserService) {
	withService := func(h func(w http.ResponseWriter, r *http.Request)) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if userSvc == nil {
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "service unavailable"})
				return
			}
			h(w, r)
		}
	}

	mux.HandleFunc("GET /users", withService(func(w http.ResponseWriter, r *http.Request) {
		listUsers(w, r, userSvc)
	}))

	mux.HandleFunc("POST /users", withService(func(w http.ResponseWriter, r *http.Request) {
		createUser(w, r, userSvc)
	}))

	mux.HandleFunc("GET /users/{id}", withService(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		getUser(w, r, userSvc, id)
	}))

	mux.HandleFunc("PUT /users/{id}", withService(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		updateUser(w, r, userSvc, id)
	}))

	mux.HandleFunc("DELETE /users/{id}", withService(func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		deleteUser(w, r, userSvc, id)
	}))
}

func createUser(w http.ResponseWriter, r *http.Request, svc *service.UserService) {
	var req service.CreateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user, err := svc.Create(r.Context(), req)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func listUsers(w http.ResponseWriter, r *http.Request, svc *service.UserService) {
	users, err := svc.List(r.Context())
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, users)
}

func getUser(w http.ResponseWriter, r *http.Request, svc *service.UserService, id string) {
	user, err := svc.GetByID(r.Context(), id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func updateUser(w http.ResponseWriter, r *http.Request, svc *service.UserService, id string) {
	var req service.UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	user, err := svc.Update(r.Context(), id, req)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func deleteUser(w http.ResponseWriter, r *http.Request, svc *service.UserService, id string) {
	if err := svc.Delete(r.Context(), id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusNoContent, nil)
}

func decodeJSON(r *http.Request, dst any) error {
	if r.Body == nil {
		return fmt.Errorf("request body is required")
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(dst)
}

func writeJSON(w http.ResponseWriter, statusCode int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if payload != nil {
		if err := json.NewEncoder(w).Encode(payload); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}
