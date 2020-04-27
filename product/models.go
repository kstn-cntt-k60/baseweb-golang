package product

import (
	"baseweb/basic"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	Id          int64               `json:"id" db:"id"`
	Name        string              `json:"name" db:"name"`
	CreatedBy   uuid.UUID           `json:"createdBy" db:"created_by_user_login_id"`
	Description string              `json:"description" db:"description"`
	Weight      decimal.NullDecimal `json:"weight" db:"weight"`
	WeightUomId string              `json:"weightUomId" db:"weight_uom_id"`
	UnitUomId   string              `json:"unitUomId" db:"unit_uom_id"`
	CreatedAt   time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time           `json:"updatedAt" db:"updated_at"`
}

type ClientProduct struct {
	Id          int64               `json:"id" db:"id"`
	Name        string              `json:"name" db:"name"`
	CreatedBy   string              `json:"createdBy" db:"created_by"`
	Description string              `json:"description" db:"description"`
	Weight      decimal.NullDecimal `json:"weight" db:"weight"`
	WeightUomId string              `json:"weightUomId" db:"weight_uom_id"`
	UnitUomId   string              `json:"unitUomId" db:"unit_uom_id"`
	CreatedAt   time.Time           `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time           `json:"updatedAt" db:"updated_at"`
}

type ProductPrice struct {
	Id            uuid.UUID       `json:"id" db:"id"`
	ProductId     int64           `json:"productId" db:"product_id"`
	CurrencyUomId string          `json:"currencyUomId" db:"currency_uom_id"`
	Price         decimal.Decimal `json:"price" db:"price"`
	CreatedBy     string          `json:"createdBy" db:"created_by"`
	EffectiveFrom time.Time       `json:"effectiveFrom" db:"effective_from"`
	ExpiredAt     basic.NullTime  `json:"expiredAt" db:"expired_at"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time       `json:"updatedAt" db:"updated_at"`
}

type InsertionPrice struct {
	ProductId     int64           `json:"productId" db:"product_id"`
	Price         decimal.Decimal `json:"price" db:"price"`
	CurrencyUomId string          `json:"currencyUomId" db:"currency_uom_id"`
	CreatedBy     uuid.UUID       `json:"createdBy" db:"created_by_user_login_id"`
	EffectiveFrom time.Time       `json:"effectiveFrom" db:"effective_from"`
}
