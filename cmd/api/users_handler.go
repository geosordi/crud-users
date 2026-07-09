package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"crud-users/internal/service"
)

func registerUserRoutes(mux *http.ServeMux, userSvc *service.UserService) {
	mux.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		if userSvc == nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "service unavailable"})
			return
		}

		switch r.Method {
		case http.MethodGet:
			listUsers(w, r, userSvc)
		case http.MethodPost:
			createUser(w, r, userSvc)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})

	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		if userSvc == nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "service unavailable"})
			return
		}

		id := strings.TrimPrefix(r.URL.Path, "/users/")
		if id == "" {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}

		switch r.Method {
		case http.MethodGet:
			getUser(w, r, userSvc, id)
		case http.MethodPut:
			updateUser(w, r, userSvc, id)
		case http.MethodDelete:
			deleteUser(w, r, userSvc, id)
		default:
			writeJSON(w, http.StatusMethodNotAllowed, map[string]string{"error": "method not allowed"})
		}
	})
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
