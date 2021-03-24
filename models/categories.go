package models

import (
	"github.com/google/uuid"
	"time"
)

type Categories struct {
	Id        uuid.UUID
	Name      string
	Picture   string
	CreatedBy uuid.UUID
	CreatedOn time.Time
	UpdatedOn time.Time
}
