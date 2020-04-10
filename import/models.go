package importProduct

import (
	"baseweb/basic"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Product struct {
	Id            int64               `json:"id" db:"id"`
	Name          string              `json:"name" db:"name"`
	Weight        decimal.NullDecimal `json:"weight" db:"weight"`
	WeightUomId   string              `json:"weightUomId" db:"weight_uom_id"`
	UnitUomId     string              `json:"unitUomId" db:"unit_uom_id"`
	TotalQuantity decimal.Decimal     `json:"totalQuantity" db:"total_quantity"`
	UpdatedAt     basic.NullTime      `json:"updatedAt" db:"updated_at"`
}

type InventoryItem struct {
	Id            int64           `json:"id" db:"id"`
	ProductId     int64           `json:"productid" db:"product_id"`
	WarehouseId   uuid.UUID       `json:"warehouseId" db:"facility_id"`
	Quantity      decimal.Decimal `json:"quantity" db:"quantity"`
	UnitCost      decimal.Decimal `json:"unitCost" db:"unit_cost"`
	CurrencyUomId string          `json:"currencyUomId" db:"currency_uom_id"`
	CreatedAt     time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time       `json:"updatedAt" db:"updated_at"`
}
