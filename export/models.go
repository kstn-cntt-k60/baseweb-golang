package export

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type InventoryItemDetail struct {
	Id               uuid.UUID       `json:"id" db:"id"`
	InventoryItemId  int64           `json:"inventoryItemId" db:"inventory_item_id"`
	ExportedQuantity decimal.Decimal `json:"exportedQuantity" db:"exported_quantity"`
	EffectiveFrom    time.Time       `json:"effectiveFrom" db:"effective_from"`
	SaleOrderId      int64           `json:"saleOrderId" db:"sale_order_id"`
	SaleOrderSeq     int             `json:"saleOrderSeq" db:"sale_order_seq"`
	CreatedAt        time.Time       `json:"createdAt" db:"created_at"`
	UpdatedAt        time.Time       `json:"updatedAt" db:"updated_at"`
}
