package json_serializer

import (
	"encoding/json"
	"github.com/google/uuid"
	"gitlab.com/ProtectIdentity/pugcha-backend/models"
	"time"
)

type CreateCategoriesResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type UpdateCategoryResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type TokenSerializer struct {
	AccessToken  string
	RefreshToken string
}

type GoogleResponse struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

type FacebookResponse struct {
}

type UserResponse struct {
	ID       uuid.UUID    `json:"id"`
	Status   string       `json:"-"`
	Password string       `json:"-"`
	Email    string       `json:"email"`
	Type     string       `json:"type"`
	Profile  *interface{} `json:"profile"`
}

type SmallBusinessResponse struct {
	BusinessName string          `json:"name"`
	BusinessId   uuid.UUID       `json:"id"`
	CategoryName *string         `json:"category_name"`
	UserId       uuid.UUID       `json:"-"`
	Location     *models.Location `json:"location"`
	Picture      []string        `json:"picture"`
	Description  *string         `json:"description"`
	Contact      []string        `json:"contact"`
	Distance     string          `json:"distance"`
	FullCount    int             `json:"total"`
}

type FullBusinessResponse struct {
	BusinessId      uuid.UUID              `json:"id"`
	UserId          uuid.UUID              `json:"-"`
	CategoryName    *string                `json:"category_name"`
	BusinessName    *string                `json:"business_name"`
	PanNumber       *string                `json:"pan_number"`
	VatNumber       *string                `json:"vat_number"`
	Email           *string                `json:"email"`
	Website         *string                `json:"website"`
	EstablishedDate *time.Time             `json:"established_date"`
	Address         *[]models.Address      `json:"address"`
	Opening         []*models.OpeningHours `json:"opening"`
	Location        *models.Location        `json:"location"`
	Picture         []string               `json:"picture"`
	Description     *string                `json:"description"`
	Contact         []string               `json:"contact"`
	Distance        string                 `json:"distance"`
	FullCount       int                    `json:"total"`
}

func (u *FullBusinessResponse) MarshalJSON() ([]byte, error) {
	type Alias FullBusinessResponse

	if u.EstablishedDate != nil {
		open := (*u.EstablishedDate).String()[0:4]
		return json.Marshal(&struct {
			EstablishedDate *string `json:"established_date"`
			*Alias
		}{
			EstablishedDate: &open,
			Alias:           (*Alias)(u),
		})
	}
	return json.Marshal(&struct {
		*Alias
	}{(*Alias)(u)})
}
