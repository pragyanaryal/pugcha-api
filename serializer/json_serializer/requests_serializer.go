package json_serializer

import (
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"time"
)

type CreateCategoriesRequest struct {
	Name      string `validate:"required"`
	Picture   string `validate:"required"`
	CreatedBy uuid.UUID
}

type UpdateAbleBusiness struct {
	Approved        *bool                  `form:"approved,omitempty" validate:"omitempty"`
	UserId          *uuid.UUID             `form:"user_id,omitempty" validate:"omitempty"`
	ApprovedBy      *uuid.UUID             `form:"approved_by,omitempty" validate:"omitempty"`
	Blocked         *bool                  `form:"blocked,omitempty" validate:"omitempty"`
	CategoryId      *uuid.UUID             `form:"category_id,omitempty" validate:"omitempty"`
	Name            *string                `form:"name,omitempty" validate:"omitempty"`
	PanNumber       *string                `form:"pan_number,omitempty" validate:"omitempty"`
	Email           *string                `form:"email,omitempty" validate:"omitempty"`
	Website         *string                `form:"website,omitempty" validate:"omitempty"`
	VatNumber       *string                `form:"vat_number,omitempty" validate:"omitempty"`
	Contact         *[]string              `form:"contact,omitempty" validate:"omitempty,dive,len=10,numeric"`
	Location        *models.Location       `form:"location,omitempty" validate:"omitempty"`
	EstablishedDate *time.Time             `form:"established_date,omitempty" validate:"omitempty"`
	Owner           *[]models.Owner        `form:"owner,omitempty" validate:"omitempty,dive"`
	Address         *[]models.Address      `form:"address,omitempty" validate:"omitempty"`
	OpeningHours    *[]models.OpeningHours `form:"opening_hours,omitempty" validate:"omitempty,unique=WeekDay"`
	Description     *string                `form:"description,omitempty" validate:"omitempty"`
	Picture         *[]string              `form:"picture,omitempty" validate:"omitempty"`
}

type UpdateBusinessPicture struct {
	Picture *[]string `form:"picture,omitempty" validate:"omitempty"`
}

type UpdateAbleUser struct {
	Status      *string `form:"status,omitempty" validate:"omitempty,status"`
	Type        *string `form:"type,omitempty" validate:"omitempty,type"`
	FirstName   *string `form:"first_name,omitempty" validate:"omitempty,alphaunicode"`
	LastName    *string `form:"last_name,omitempty" validate:"omitempty,alphaunicode"`
	MiddleName  *string `form:"middle_name,omitempty" validate:"omitempty,alphaunicode"`
	Dob         *string `form:"dob,omitempty" validate:"omitempty,datetime=2006-01-02"`
	Gender      *string `form:"gender,omitempty" validate:"omitempty,gender"`
	Contact     *string `form:"contact,omitempty" validate:"omitempty,len=10,numeric"`
	ProfilePics *string `form:"profile_pics,omitempty" validate:"omitempty"`
}

type CreateNormalUser struct {
	SocialID       string `json:"id"`
	SocialPlatform string `json:"platform"`
	Email          string `json:"email"`
	VerifiedEmail  bool   `json:"verified_email"`
	FirstName      string `json:"given_name"`
	LastName       string `json:"family_name"`
	Name           string `json:"name"`
	Picture        string `json:"picture"`
}

type CreateAdminUserRequest struct {
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,alphaunicode"`
	LastName  string `json:"last_name,omitempty" validate:"alphaunicode"`
}

type PasswordSerializer struct {
	Password string `json:"password" validate:"required"`
}

type LoginSerializer struct {
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type Field struct {
	Name       string
	FilterOp   string
	FilterData []interface{}
}

type FilterRequest struct {
	Limit  uint16
	Offset uint16
	Fields map[string]*Field
	Sort   []string
}

type UpdateCategoriesRequest struct {
	Name    string
	Picture string
}
