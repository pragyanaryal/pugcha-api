package models

import (
	"github.com/google/uuid"
)

type Promotion struct {
	Id            uuid.UUID
	BusinessId    uuid.UUID
	StartedDate   string
	FinishDate    string
	Duration      string
	Active        bool
	Model         PromotionModel
	PromotionType string
}

type PromotionModel struct {
	Price          float64
	InitialPayment float64
	PaymentTillNow float64
	Payment        Installments
}

type Installments struct {
	Price       float64
	PaidOn      string
	PaymentType string
}
