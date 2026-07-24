package auth

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jyanza11/qeko-backend/internal/domain"
	"gorm.io/gorm"
)

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("conflict")
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateOrganizer(
	ctx context.Context,
	email, passwordHash, fullName, orgName, orgSlug string,
) (*domain.User, *domain.Organization, error) {
	user := &domain.User{
		Email:        email,
		PasswordHash: passwordHash,
		FullName:     fullName,
	}
	org := &domain.Organization{
		Name: orgName,
		Slug: orgSlug,
	}

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return mapDBError(err)
		}

		if err := tx.Create(org).Error; err != nil {
			return mapDBError(err)
		}

		member := &domain.OrganizationMember{
			OrganizationID: org.ID,
			UserID:         user.ID,
			Role:           "owner",
		}
		if err := tx.Create(member).Error; err != nil {
			return mapDBError(err)
		}

		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return user, org, nil
}

func (r *Repository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindUserByID(ctx context.Context, id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindPrimaryOrganization(ctx context.Context, userID uuid.UUID) (*domain.Organization, error) {
	var org domain.Organization
	err := r.db.WithContext(ctx).
		Joins("JOIN organization_members AS m ON m.organization_id = organizations.id").
		Where("m.user_id = ?", userID).
		Order("m.created_at ASC").
		First(&org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &org, nil
}

func (r *Repository) SlugExists(ctx context.Context, slug string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&domain.Organization{}).
		Where("slug = ?", slug).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func mapDBError(err error) error {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return ErrConflict
	}
	return err
}
