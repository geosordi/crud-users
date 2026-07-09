package main

import "net/http"

const swaggerSpec = `{
  "swagger": "2.0",
  "info": {
    "title": "CRUD Users API",
    "version": "1.0.0",
    "description": "API simples para cadastro, listagem, atualização e remoção de usuários"
  },
  "host": "localhost:8080",
  "schemes": ["http"],
  "paths": {
    "/health": {
      "get": {
        "summary": "Verifica se a API está funcionando",
        "responses": { "200": { "description": "ok" } }
      }
    },
    "/users": {
      "get": {
        "summary": "Lista todos os usuários",
        "responses": { "200": { "description": "Lista de usuários" } }
      },
      "post": {
        "summary": "Cria um novo usuário",
        "parameters": [
          {
            "name": "body",
            "in": "body",
            "required": true,
            "schema": { "$ref": "#/definitions/UserInput" }
          }
        ],
        "responses": { "201": { "description": "Usuário criado" } }
      }
    },
    "/users/{id}": {
      "get": {
        "summary": "Busca um usuário por ID",
        "parameters": [
          { "name": "id", "in": "path", "required": true, "type": "string" }
        ],
        "responses": { "200": { "description": "Usuário encontrado" } }
      },
      "put": {
        "summary": "Atualiza um usuário",
        "parameters": [
          { "name": "id", "in": "path", "required": true, "type": "string" },
          { "name": "body", "in": "body", "required": true, "schema": { "$ref": "#/definitions/UserInput" } }
        ],
        "responses": { "200": { "description": "Usuário atualizado" } }
      },
      "delete": {
        "summary": "Remove um usuário",
        "parameters": [
          { "name": "id", "in": "path", "required": true, "type": "string" }
        ],
        "responses": { "204": { "description": "Usuário removido" } }
      }
    }
  },
  "definitions": {
    "UserInput": {
      "type": "object",
      "required": ["name", "email"],
      "properties": {
        "name": { "type": "string" },
        "email": { "type": "string" }
      }
    }
  }
}`

const swaggerUIHTML = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <title>CRUD Users API Docs</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = () => {
        SwaggerUIBundle({
          url: '/swagger/doc.json',
          dom_id: '#swagger-ui'
        });
      };
    </script>
  </body>
</html>`

func registerSwaggerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
	})

	mux.HandleFunc("/swagger/", func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/swagger/" || r.URL.Path == "/swagger/index.html":
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = w.Write([]byte(swaggerUIHTML))
		case r.URL.Path == "/swagger/doc.json":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(swaggerSpec))
		default:
			http.NotFound(w, r)
		}
	})
}
