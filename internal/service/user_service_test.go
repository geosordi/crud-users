package service

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"crud-users/internal/models"
)

type stubUserRepository struct {
	users  map[string]models.User
	nextID int
}

func (s *stubUserRepository) Create(ctx context.Context, user *models.User) (*models.User, error) {
	if user.Email == "duplicate@example.com" {
		return nil, errors.New("email already exists")
	}
	if user.ID == "" {
		s.nextID++
		user.ID = fmt.Sprintf("generated-id-%d", s.nextID)
	}
	s.users[user.ID] = *user
	return user, nil
}

func (s *stubUserRepository) List(ctx context.Context) ([]models.User, error) {
	users := make([]models.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

func (s *stubUserRepository) GetByID(ctx context.Context, id string) (*models.User, error) {
	user, exists := s.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return &user, nil
}

func (s *stubUserRepository) Update(ctx context.Context, id string, user models.User) error {
	_, exists := s.users[id]
	if !exists {
		return errors.New("user not found")
	}
	s.users[id] = user
	return nil
}

func (s *stubUserRepository) Delete(ctx context.Context, id string) error {
	_, exists := s.users[id]
	if !exists {
		return errors.New("user not found")
	}
	delete(s.users, id)
	return nil
}

func TestCreateUser(t *testing.T) {
	repo := &stubUserRepository{users: make(map[string]models.User)}
	svc := NewUserService(repo)

	user, err := svc.Create(context.Background(), CreateUserRequest{Name: "Alice", Email: "alice@example.com"})
	if err != nil {
		t.Fatalf("expected user to be created, got error: %v", err)
	}
	if user.ID == "" {
		t.Fatal("expected generated user id")
	}
	if user.Name != "Alice" || user.Email != "alice@example.com" {
		t.Fatalf("unexpected user payload: %+v", user)
	}
}

func TestCreateUserRejectsInvalidEmail(t *testing.T) {
	repo := &stubUserRepository{users: make(map[string]models.User)}
	svc := NewUserService(repo)

	_, err := svc.Create(context.Background(), CreateUserRequest{Name: "Bob", Email: "invalid-email"})
	if err == nil {
		t.Fatal("expected invalid email error")
	}
}

func TestUpdateUser(t *testing.T) {
	repo := &stubUserRepository{users: make(map[string]models.User)}
	svc := NewUserService(repo)

	created, err := svc.Create(context.Background(), CreateUserRequest{Name: "Carol", Email: "carol@example.com"})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}

	updated, err := svc.Update(context.Background(), created.ID, UpdateUserRequest{Name: "Carol Updated", Email: "carol.updated@example.com"})
	if err != nil {
		t.Fatalf("expected update to succeed: %v", err)
	}

	if updated.Name != "Carol Updated" || updated.Email != "carol.updated@example.com" {
		t.Fatalf("unexpected updated user: %+v", updated)
	}
	if updated.UpdatedAt.Before(created.UpdatedAt) && !updated.UpdatedAt.Equal(created.UpdatedAt) {
		t.Fatalf("expected updated timestamp to move forward")
	}
	if updated.CreatedAt != created.CreatedAt {
		t.Fatalf("expected created timestamp to stay unchanged")
	}
}

func TestDeleteUser(t *testing.T) {
	repo := &stubUserRepository{users: make(map[string]models.User)}
	svc := NewUserService(repo)

	created, err := svc.Create(context.Background(), CreateUserRequest{Name: "Diana", Email: "diana@example.com"})
	if err != nil {
		t.Fatalf("expected create to succeed: %v", err)
	}

	if err := svc.Delete(context.Background(), created.ID); err != nil {
		t.Fatalf("expected delete to succeed: %v", err)
	}

	_, err = svc.GetByID(context.Background(), created.ID)
	if err == nil {
		t.Fatal("expected user to be gone after delete")
	}
}

func TestListUsers(t *testing.T) {
	repo := &stubUserRepository{users: make(map[string]models.User)}
	svc := NewUserService(repo)

	_, _ = svc.Create(context.Background(), CreateUserRequest{Name: "One", Email: "one@example.com"})
	_, _ = svc.Create(context.Background(), CreateUserRequest{Name: "Two", Email: "two@example.com"})

	users, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("expected list to succeed: %v", err)
	}
	if len(users) != 2 {
		t.Fatalf("expected 2 users, got %d", len(users))
	}
	if users[0].CreatedAt.IsZero() || users[1].CreatedAt.IsZero() {
		t.Fatalf("expected created timestamps to be set, got %+v", users)
	}
	if users[0].UpdatedAt.IsZero() || users[1].UpdatedAt.IsZero() {
		t.Fatalf("expected updated timestamps to be set, got %+v", users)
	}
	if time.Since(users[0].CreatedAt) < 0 {
		t.Fatalf("expected created timestamp to be in the past")
	}
}
