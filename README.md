# CRUD de Usuários em Go

Este projeto é um CRUD simples de usuários implementado em Go com MongoDB, pensado como um exercício de aprendizado para consolidar conceitos básicos de API REST, persistência, Docker e testes.

## O que foi implementado

- API REST com endpoints para criar, listar, buscar, atualizar e deletar usuários
- Persistência em MongoDB
- Execução do banco via Docker Compose usando MongoDB 8
- Configuração via arquivo .env
- Documentação Swagger disponível na aplicação
- Estrutura organizada em camadas: handlers, serviços, repositório, modelos e configuração
- Testes para o fluxo de serviço e para o roteamento básico da API

## Tecnologias

- Go
- MongoDB
- Docker Compose
- Swagger UI

## Pré-requisitos

- Go 1.26+
- Docker
- Docker Compose

## Configuração do ambiente

Copie o arquivo de exemplo e ajuste os valores conforme o seu ambiente:

```bash
cp .env.example .env
```

As variáveis abaixo são obrigatórias:

- MONGO_URI: URI de conexão com o MongoDB
- MONGO_DB: nome do banco de dados
- MONGO_COLLECTION: nome da coleção

Exemplo padrão fornecido em .env.example:

```env
MONGO_URI=mongodb://localhost:27017
MONGO_DB=crud_users
MONGO_COLLECTION=users
```

## Como subir o banco de dados

```bash
docker compose up -d
```

## Como rodar a aplicação

```bash
go run ./cmd/api
```

A aplicação ficará disponível em:

- Health check: http://localhost:8080/health
- Swagger UI: http://localhost:8080/swagger
- Swagger JSON: http://localhost:8080/swagger/doc.json

## Endpoints da API

### Criar usuário

```bash
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice","email":"alice@example.com"}'
```

### Listar usuários

```bash
curl http://localhost:8080/users
```

### Buscar usuário por ID

```bash
curl http://localhost:8080/users/{id}
```

### Atualizar usuário

```bash
curl -X PUT http://localhost:8080/users/{id} \
  -H "Content-Type: application/json" \
  -d '{"name":"Alice Updated","email":"alice.updated@example.com"}'
```

### Deletar usuário

```bash
curl -X DELETE http://localhost:8080/users/{id}
```

## Testes

Para executar a suíte de testes:

```bash
go test ./...
```
