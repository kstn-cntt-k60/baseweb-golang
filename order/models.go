package order

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CustomerStore struct {
	Id         uuid.UUID `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	Customer   string    `json:"customer" db:"customer_name"`
	CustomerId uuid.UUID `json:"customerId" db:"customer_id"`
	Address    string    `json:"address" db:"address"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type ClientProduct struct {
	Id       int             `json:"id"`
	Quantity decimal.Decimal `json:"quantity"`
}

type ProductInfo struct {
	Id                int                 `json:"id" db:"id"`
	Name              string              `json:"name" db:"name"`
	CreatedBy         string              `json:"createdBy" db:"created_by"`
	Weight            decimal.NullDecimal `json:"weight" db:"weight"`
	WeightUomId       string              `json:"weightUomId" db:"weight_uom_id"`
	UnitUomId         string              `json:"unitUomId" db:"unit_uom_id"`
	CurrencyUomId     string              `json:"currencyUomId" db:"currency_uom_id"`
	Price             decimal.Decimal     `json:"price" db:"price"`
	EffectiveFrom     time.Time           `json:"effectiveFrom" db:"effective_from"`
	QuantityAvailable decimal.Decimal     `json:"quantityAvailable" db:"quantity_available"`
	CreatedAt         time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt         time.Time           `json:"updatedAt" db:"updated_at"`
}
