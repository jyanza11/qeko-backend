package auth

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jyanza11/qeko-backend/internal/domain"
	appErrors "github.com/jyanza11/qeko-backend/internal/shared/errors"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo   *Repository
	tokens *TokenManager
}

func NewService(repo *Repository, tokens *TokenManager) *Service {
	return &Service{repo: repo, tokens: tokens}
}

func (s *Service) Register(ctx context.Context, req RegisterRequest) (*AuthResponse, *appErrors.AppError) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, appErrors.Internal("failed to hash password")
	}

	slug, appErr := s.uniqueSlug(ctx, req.OrganizationName)
	if appErr != nil {
		return nil, appErr
	}

	user, org, err := s.repo.CreateOrganizer(
		ctx,
		strings.ToLower(strings.TrimSpace(req.Email)),
		string(passwordHash),
		strings.TrimSpace(req.FullName),
		strings.TrimSpace(req.OrganizationName),
		slug,
	)
	if errors.Is(err, ErrConflict) {
		return nil, appErrors.Conflict("email or organization already exists")
	}
	if err != nil {
		return nil, appErrors.Internal("failed to register organizer")
	}

	return s.issueAuth(user, org)
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, *appErrors.AppError) {
	user, err := s.repo.FindUserByEmail(ctx, strings.ToLower(strings.TrimSpace(req.Email)))
	if errors.Is(err, ErrNotFound) {
		return nil, appErrors.Unauthorized("invalid email or password")
	}
	if err != nil {
		return nil, appErrors.Internal("failed to login")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, appErrors.Unauthorized("invalid email or password")
	}

	org, err := s.repo.FindPrimaryOrganization(ctx, user.ID)
	if errors.Is(err, ErrNotFound) {
		return nil, appErrors.Forbidden("user has no organization")
	}
	if err != nil {
		return nil, appErrors.Internal("failed to login")
	}

	return s.issueAuth(user, org)
}

func (s *Service) Me(ctx context.Context, userID uuid.UUID) (*MeResponse, *appErrors.AppError) {
	user, err := s.repo.FindUserByID(ctx, userID)
	if errors.Is(err, ErrNotFound) {
		return nil, appErrors.Unauthorized("invalid token")
	}
	if err != nil {
		return nil, appErrors.Internal("failed to load user")
	}

	org, err := s.repo.FindPrimaryOrganization(ctx, user.ID)
	if errors.Is(err, ErrNotFound) {
		return nil, appErrors.Forbidden("user has no organization")
	}
	if err != nil {
		return nil, appErrors.Internal("failed to load organization")
	}

	return &MeResponse{
		User:         toUserResponse(user),
		Organization: toOrganizationResponse(org),
	}, nil
}

func (s *Service) issueAuth(user *domain.User, org *domain.Organization) (*AuthResponse, *appErrors.AppError) {
	token, _, err := s.tokens.Generate(user.ID, org.ID, user.Email)
	if err != nil {
		return nil, appErrors.Internal("failed to generate token")
	}

	return &AuthResponse{
		AccessToken:  token,
		TokenType:    "Bearer",
		ExpiresIn:    s.tokens.ExpiresInSeconds(),
		User:         toUserResponse(user),
		Organization: toOrganizationResponse(org),
	}, nil
}

func (s *Service) uniqueSlug(ctx context.Context, name string) (string, *appErrors.AppError) {
	base := slugify(name)
	slug := base

	for i := 0; i < 5; i++ {
		exists, err := s.repo.SlugExists(ctx, slug)
		if err != nil {
			return "", appErrors.Internal("failed to create organization slug")
		}
		if !exists {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%s", base, uuid.New().String()[:8])
	}

	return "", appErrors.Internal("failed to create organization slug")
}

func toUserResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:       user.ID,
		Email:    user.Email,
		FullName: user.FullName,
	}
}

func toOrganizationResponse(org *domain.Organization) OrganizationResponse {
	return OrganizationResponse{
		ID:   org.ID,
		Name: org.Name,
		Slug: org.Slug,
	}
}
