package importProduct

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db                     *sqlx.DB
	countProduct           *sqlx.Stmt
	viewProductByWarehouse *sqlx.Stmt
}

func InitRepo(db *sqlx.DB) *Repo {
	query := `select count(*) from product`
	countProduct, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Fatal(err)
	}

	query = `select p.id, p.name, 
        p.weight, p.weight_uom_id, 
        p.unit_uom_id, coalesce(g.total_quantity, 0) as total_quantity,
        g.updated_at from product p
        left join (
            select i.product_id, sum(i.quantity) as total_quantity, 
            max(i.updated_at) as updated_at from inventory_item i
            where i.facility_id = ?
            group by i.product_id
        ) g on g.product_id = p.id
        order by updated_at desc nulls last
        limit ? offset ?`
	viewProductByWarehouse, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Fatal(err)
	}

	return &Repo{
		db:                     db,
		countProduct:           countProduct,
		viewProductByWarehouse: viewProductByWarehouse,
	}
}

func (repo *Repo) ViewProductByWarehouse(
	ctx context.Context, warehouseId uuid.UUID) (uint, []Product, error) {

	log.Println("ViewProductByWarehouse", warehouseId)

	var count uint
	result := make([]Product, 0)
	var err error

	err = repo.countProduct.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	err = repo.viewProductByWarehouse.SelectContext(
		ctx, &result, warehouseId, 5, 0)

	return count, result, err
}

func (repo *Repo) InsertInventoryItem(
	ctx context.Context, item InventoryItem) error {

	log.Println("InsertInventoryItem", item.ProductId, item.WarehouseId,
		item.Quantity, item.UnitCost, item.CurrencyUomId)

	query := `insert into inventory_item(product_id, 
    facility_id, quantity, unit_cost, currency_uom_id)
    values (:product_id, :facility_id, :quantity, 
        :unit_cost, :currency_uom_id)`

	_, err := repo.db.NamedExecContext(ctx, query, item)

	return err

}
