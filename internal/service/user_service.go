package service

import (
	"context"
	"errors"
	"fmt"
	"net/mail"
	"time"

	"crud-users/internal/models"
	"crud-users/internal/repository"
)

type CreateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Create(ctx context.Context, req CreateUserRequest) (*models.User, error) {
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	user := &models.User{
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}

func (s *UserService) List(ctx context.Context) ([]models.User, error) {
	return s.repo.List(ctx)
}

func (s *UserService) GetByID(ctx context.Context, id string) (*models.User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *UserService) Update(ctx context.Context, id string, req UpdateUserRequest) (*models.User, error) {
	if err := validateEmail(req.Email); err != nil {
		return nil, err
	}

	current, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	updated := *current
	updated.Name = req.Name
	updated.Email = req.Email
	updated.UpdatedAt = time.Now().UTC()

	if err := s.repo.Update(ctx, id, updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (s *UserService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func validateEmail(email string) error {
	if email == "" {
		return errors.New("email is required")
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("invalid email: %w", err)
	}
	return nil
}
