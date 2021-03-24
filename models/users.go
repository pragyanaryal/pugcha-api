package models

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID        uuid.UUID `json:"id,omitempty"`
	Email     string    `json:"email,omitempty"`
	Status    string    `json:"status,omitempty"`
	Password  string    `json:"password,omitempty"`
	Type      string    `json:"type,omitempty"`
	CreatedOn time.Time `json:"created_on,omitempty"`
	UpdatedOn time.Time `json:"updated_on,omitempty"`
}

type UserProfile struct {
	UserId      uuid.UUID
	Email       string
	FirstName   string
	MiddleName  *string
	LastName    string
	Dob         *time.Time
	Gender      *string
	Contact     *string
	ProfilePics string
	CreatedOn   time.Time `json:"-"`
	UpdatedOn   time.Time `json:"-"`
}

type GoogleAccount struct {
	UserId     uuid.UUID
	Email      string
	GoogleID   string
	AccessKey  string
	RefreshKey string
}

type FacebookAccount struct {
	UserId     uuid.UUID
	Email      string
	FacebookID string
	AccessKey  string
	RefreshKey string
}
