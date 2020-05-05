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
	TotalQuantity decimal.Decimal     `json:"totalQuantity" db:"quantity_total"`
	UpdatedAt     basic.NullTime      `json:"updatedAt" db:"updated_at"`
}

type InventoryItem struct {
	Id             int64           `json:"id" db:"id"`
	ProductId      int64           `json:"productId" db:"product_id"`
	ProductName    string          `json:"productName" db:"product_name"`
	WarehouseId    uuid.UUID       `json:"warehouseId" db:"warehouse_id"`
	Quantity       decimal.Decimal `json:"quantity" db:"quantity"`
	QuantityOnHand decimal.Decimal `json:"quantityOnHand" db:"quantity_on_hand"`
	UnitCost       decimal.Decimal `json:"unitCost" db:"unit_cost"`
	CurrencyUomId  string          `json:"currencyUomId" db:"currency_uom_id"`
	CreatedAt      time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt      time.Time       `json:"updatedAt" db:"updated_at"`
}

type WarehouseProductStatistics struct {
	WarehouseId        uuid.UUID       `json:"warehouseId" db:"warehouse_id"`
	ProductId          int64           `json:"productId" db:"product_id"`
	InventoryItemCount int64           `json:"inventoryItemCount" db:"inventory_item_count"`
	QuantityTotal      decimal.Decimal `json:"quantityTotal" db:"quantity_total"`
	QuantityOnHand     decimal.Decimal `json:"quantityOnHand" db:"quantity_on_hand"`
	QuantityAvailable  decimal.Decimal `json:"quantityAvailable" db:"quantity_available"`
	CreatedAt          time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt          time.Time       `json:"updatedAt" db:"updated_at"`
}
