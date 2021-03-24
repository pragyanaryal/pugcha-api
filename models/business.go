package models

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/jackc/pgtype"
	"time"
)

type Weekday int

const (
	Sunday Weekday = iota + 1
	Monday
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
)

type Business struct {
	ID         uuid.UUID
	Approved   bool
	UserId     uuid.UUID
	ApprovedBy uuid.UUID
	Blocked    bool
	CreatedOn  time.Time `json:"-"`
	UpdatedOn  time.Time `json:"-"`
}

type BusinessProfile struct {
	BusinessId      uuid.UUID      `form:"business_id,omitempty"`
	CategoryId      uuid.UUID      `form:"category_id,omitempty" validate:"required"`
	UserId          uuid.UUID      `form:"user_id,omitempty" validate:"required"`
	Name            string         `form:"name,omitempty" validate:"required"`
	PanNumber       *string        `form:"pan_number,omitempty"`
	Email           *string        `form:"email,omitempty" validate:"omitempty,email"`
	Website         *string        `form:"website,omitempty"`
	VatNumber       *string        `form:"vat_number,omitempty"`
	Picture         []string       `form:"image,omitempty" validate:"required"`
	Contact         []string       `form:"contact,omitempty" validate:"required,dive,numeric"`
	Location        Location       `form:"location,omitempty" validate:"omitempty"`
	EstablishedDate *time.Time     `form:"established_date,omitempty"`
	Owner           []Owner        `form:"owner,omitempty" validate:"omitempty,dive"`
	Address         []Address      `form:"address,omitempty"`
	OpeningHours    []OpeningHours `form:"opening_hours,omitempty" validate:"omitempty,unique=WeekDay"`
	Description     *string        `form:"description,omitempty"`
	CreatedOn       time.Time      `json:"-" form:"created_on,omitempty"`
	UpdatedOn       time.Time      `json:"-" form:"updated_on,omitempty"`
	Geom            *string        `json:"-"`
}

type Location struct {
	Latitude  *float64 `form:"latitude,omitempty" json:"latitude"`
	Longitude *float64 `form:"longitude,omitempty" json:"longitude"`
}

type OpeningHours struct {
	BusinessId  uuid.UUID `json:"-"`
	WeekDay     Weekday   `json:"week_day"`
	OpeningTime *string   `form:"opening_time,omitempty" json:"opening_time"`
	ClosingTime *string   `form:"closing_time,omitempty" json:"closing_time"`
	Opened      bool      `form:"opened,omitempty" json:"opened"`
}

type Owner struct {
	Id         uuid.UUID  `json:"id"`
	BusinessId uuid.UUID  `json:"business_id"`
	UserId     *uuid.UUID `form:"id,omitempty" json:"user_id"`
	Email      string     `form:"email,omitempty" validate:"omitempty,email" json:"email"`
	Contact    []string   `form:"contact,omitempty" validate:"dive,len=10,numeric" json:"contact"`
}

type Address struct {
	Id           uuid.UUID `json:"id"`
	BusinessId   uuid.UUID `json:"-"`
	Street       string    `form:"street,omitempty" json:"street"`
	Ward         int       `form:"ward,omitempty" json:"ward"`
	Municipality string    `form:"municipality,omitempty" json:"municipality"`
	District     string    `form:"district,omitempty" json:"district"`
	State        string    `form:"state,omitempty" json:"state"`
	Country      string    `form:"country,omitempty" json:"country"`
	Contact      []string  `form:"contact,omitempty" json:"contact"`
}

func (loc *Location) DecodeBinary(ci *pgtype.ConnInfo, src []byte) error {
	if src == nil {
		return errors.New("NULL values can't be decoded. Scan into a &*Location to handle NULLs")
	}
	if err := (pgtype.CompositeFields{&loc.Latitude, &loc.Longitude}).DecodeBinary(ci, src); err != nil {
		return err
	}
	return nil
}

func (loc Location) EncodeBinary(ci *pgtype.ConnInfo, buf []byte) (newBuf []byte, err error) {
	lat := pgtype.Float8{Float: *loc.Latitude, Status: pgtype.Present}
	lon := pgtype.Float8{Float: *loc.Longitude, Status: pgtype.Present}

	return (pgtype.CompositeFields{&lat, &lon}).EncodeBinary(ci, buf)
}

func (u *OpeningHours) MarshalJSON() ([]byte, error) {
	type Alias OpeningHours

	open := (*u.OpeningTime)[0:5]
	closes := (*u.ClosingTime)[0:5]

	return json.Marshal(&struct {
		OpeningTime *string `json:"opening_time"`
		ClosingTime *string `json:"closing_time"`
		*Alias
	}{
		OpeningTime: &open,
		ClosingTime: &closes,
		Alias:    (*Alias)(u),
	})
}