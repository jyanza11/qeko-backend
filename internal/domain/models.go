package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Email        string    `gorm:"size:255;unique;not null"`
	PasswordHash string    `gorm:"column:password_hash;not null"`
	FullName     string    `gorm:"column:full_name;size:255;not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
}

func (User) TableName() string { return "users" }

type Organization struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Name      string    `gorm:"size:255;not null"`
	Slug      string    `gorm:"size:255;unique;not null"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
}

func (Organization) TableName() string { return "organizations" }

type OrganizationMember struct {
	OrganizationID uuid.UUID    `gorm:"type:uuid;primaryKey"`
	UserID         uuid.UUID    `gorm:"type:uuid;primaryKey"`
	Role           string       `gorm:"size:32;not null"`
	CreatedAt      time.Time    `gorm:"not null"`
	Organization   Organization `gorm:"constraint:OnDelete:CASCADE"`
	User           User         `gorm:"constraint:OnDelete:CASCADE"`
}

func (OrganizationMember) TableName() string { return "organization_members" }

type Event struct {
	ID             uuid.UUID    `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	OrganizationID uuid.UUID    `gorm:"type:uuid;not null;index"`
	Title          string       `gorm:"size:255;not null"`
	Description    string       `gorm:"not null;default:''"`
	Location       string       `gorm:"size:255;not null;default:''"`
	StartsAt       time.Time    `gorm:"not null;index"`
	EndsAt         time.Time    `gorm:"not null"`
	Status         string       `gorm:"size:32;not null;default:draft"`
	CreatedAt      time.Time    `gorm:"not null"`
	UpdatedAt      time.Time    `gorm:"not null"`
	Organization   Organization `gorm:"constraint:OnDelete:CASCADE"`
}

func (Event) TableName() string { return "events" }

type Ticket struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	EventID       uuid.UUID `gorm:"type:uuid;not null;index"`
	AttendeeName  string    `gorm:"column:attendee_name;size:255;not null"`
	AttendeeEmail string    `gorm:"column:attendee_email;size:255;not null;index"`
	QRToken       string    `gorm:"column:qr_token;size:255;unique;not null"`
	Status        string    `gorm:"size:32;not null;default:issued"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
	Event         Event     `gorm:"constraint:OnDelete:CASCADE"`
}

func (Ticket) TableName() string { return "tickets" }

type Checkin struct {
	ID          uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	TicketID    uuid.UUID  `gorm:"type:uuid;unique;not null"`
	EventID     uuid.UUID  `gorm:"type:uuid;not null;index"`
	CheckedInAt time.Time  `gorm:"column:checked_in_at;not null"`
	CheckedInBy *uuid.UUID `gorm:"column:checked_in_by;type:uuid"`
	Ticket      Ticket     `gorm:"constraint:OnDelete:CASCADE"`
	Event       Event      `gorm:"constraint:OnDelete:CASCADE"`
	Checker     *User      `gorm:"foreignKey:CheckedInBy;constraint:OnDelete:SET NULL"`
}

func (Checkin) TableName() string { return "checkins" }
