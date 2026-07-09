package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerEndpointsAreExposed(t *testing.T) {
	router := newRouter(nil)

	t.Run("swagger index", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
		if !strings.Contains(rr.Body.String(), "swagger") {
			t.Fatalf("expected swagger UI content, got %q", rr.Body.String())
		}
	})

	t.Run("swagger spec", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		if rr.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
		}
		body := rr.Body.String()
		if !strings.Contains(body, "\"swagger\"") || !strings.Contains(body, "\"info\"") {
			t.Fatalf("expected swagger spec content, got %q", body)
		}
	})
}
