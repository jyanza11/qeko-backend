package auth

import "github.com/google/uuid"

type RegisterRequest struct {
	FullName         string `json:"full_name" validate:"required,min=2,max=120"`
	Email            string `json:"email" validate:"required,email"`
	Password         string `json:"password" validate:"required,min=8,max=72"`
	OrganizationName string `json:"organization_name" validate:"required,min=2,max=120"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserResponse struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	FullName string    `json:"full_name"`
}

type OrganizationResponse struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Slug string    `json:"slug"`
}

type AuthResponse struct {
	AccessToken  string               `json:"access_token"`
	TokenType    string               `json:"token_type"`
	ExpiresIn    int64                `json:"expires_in"`
	User         UserResponse         `json:"user"`
	Organization OrganizationResponse `json:"organization"`
}

type MeResponse struct {
	User         UserResponse         `json:"user"`
	Organization OrganizationResponse `json:"organization"`
}
