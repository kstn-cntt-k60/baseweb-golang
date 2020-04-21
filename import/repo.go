package importProduct

import (
	"context"
	"database/sql"
	"fmt"
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
        p.unit_uom_id, coalesce(s.quantity_total, 0) as quantity_total,
        s.updated_at from product p
        left join (
            select quantity_total, product_id, updated_at
            from warehouse_product_statistics 
            where warehouse_id = ?
        ) s on s.product_id = p.id
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
	ctx context.Context, warehouseId uuid.UUID,
	page, pageSize int,
	sortedBy, sortOrder string) (uint, []Product, error) {

	log.Println("ViewProductByWarehouse", warehouseId,
		page, pageSize, sortedBy, sortOrder)

	var count uint
	result := make([]Product, 0)
	var err error

	err = repo.countProduct.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "updated_at" && sortOrder == "desc" {
		err = repo.viewProductByWarehouse.SelectContext(
			ctx, &result, warehouseId, pageSize, page*pageSize)
	} else {
		query := `select p.id, p.name,
            p.weight, p.weight_uom_id,
            p.unit_uom_id, coalesce(s.quantity_total, 0) as quantity_total,
            s.updated_at from product p
            left join (
                select quantity_total, product_id, updated_at
                from warehouse_product_statistics 
                where warehouse_id = ?
            ) s on s.product_id = p.id
            order by %s %s nulls last
            limit ? offset ?`
		query = fmt.Sprintf(query, sortedBy, sortOrder)
		query = repo.db.Rebind(query)

		err = repo.db.SelectContext(ctx, &result,
			query, warehouseId, pageSize, page*pageSize)
	}
	return count, result, err
}

func (repo *Repo) InsertInventoryItem(
	ctx context.Context, item InventoryItem) error {

	log.Println("InsertInventoryItem", item.ProductId, item.WarehouseId,
		item.Quantity, item.UnitCost, item.CurrencyUomId)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `insert into inventory_item(product_id,
    warehouse_id, quantity, unit_cost, currency_uom_id)
    values (:product_id, :warehouse_id, :quantity,
        :unit_cost, :currency_uom_id)`
	_, err = tx.NamedExecContext(ctx, query, item)
	if err != nil {
		return err
	}

	stat := WarehouseProductStatistics{}
	query = `select * from warehouse_product_statistics
        where warehouse_id = ? and product_id = ?`
	query = repo.db.Rebind(query)
	err = tx.GetContext(ctx, &stat,
		query, item.WarehouseId, item.ProductId)

	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if err == sql.ErrNoRows {
		stat.WarehouseId = item.WarehouseId
		stat.ProductId = item.ProductId
		stat.InventoryItemCount = 1
		stat.QuantityTotal = item.Quantity
		stat.QuantityOnHand = item.Quantity
		stat.QuantityAvailable = item.Quantity

		query = `insert into warehouse_product_statistics(
            warehouse_id, product_id,
            inventory_item_count, quantity_total,
            quantity_on_hand, quantity_available)
            values (:warehouse_id, :product_id,
            :inventory_item_count, :quantity_total,
            :quantity_on_hand, :quantity_available)`
		_, err = tx.NamedExecContext(ctx, query, stat)
		if err != nil {
			return err
		}
	} else {
		stat.InventoryItemCount += 1
		stat.QuantityTotal = stat.QuantityTotal.Add(item.Quantity)
		stat.QuantityOnHand = stat.QuantityOnHand.Add(item.Quantity)
		stat.QuantityAvailable = stat.QuantityAvailable.Add(item.Quantity)

		query = `update warehouse_product_statistics
            set inventory_item_count = :inventory_item_count,
                quantity_total = :quantity_total,
                quantity_on_hand = :quantity_on_hand,
                quantity_available = :quantity_available
            where warehouse_id = :warehouse_id and product_id = :product_id`
		_, err = tx.NamedExecContext(ctx, query, stat)
		if err != nil {
			return err
		}
	}

	return tx.Commit()

}

func (repo *Repo) SelectProductByWarehouse(
	ctx context.Context, warehouseId uuid.UUID) ([]Product, error) {

	log.Println("SelectProductByWarehouse", warehouseId)

	query := `select p.id, p.name,
        p.weight, p.weight_uom_id,
        p.unit_uom_id, coalesce(s.quantity_total, 0) as quantity_total,
        s.updated_at from product p
        left join (
            select quantity_total, product_id, updated_at
            from warehouse_product_statistics 
            where warehouse_id = ?
        ) s on s.product_id = p.id`
	query = repo.db.Rebind(query)

	result := make([]Product, 0)
	err := repo.db.SelectContext(ctx, &result, query, warehouseId)

	return result, err
}

func (repo *Repo) ViewInventoryItemByWarehouse(
	ctx context.Context,
	warehouseId uuid.UUID,
	page, pageSize int,
	sortedBy, sortOrder string) (int, []InventoryItem, error) {

	log.Println("ViewInventoryItemByWarehouse", warehouseId,
		page, pageSize, sortedBy, sortOrder)

	var count int
	result := make([]InventoryItem, 0)

	query := `select count(*) from inventory_item where warehouse_id = ?`
	query = repo.db.Rebind(query)
	err := repo.db.GetContext(ctx, &count, query, warehouseId)
	if err != nil {
		return count, result, err
	}

	query = `select i.id, i.product_id,
        p.name as product_name,
        i.warehouse_id, i.quantity, i.unit_cost,
        i.currency_uom_id,
        i.created_at, i.updated_at
        from inventory_item i 
        inner join product p
            on p.id = i.product_id
        where i.warehouse_id = ?
        order by i.%s %s
        limit ? offset ?`
	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)

	err = repo.db.SelectContext(ctx, &result, query,
		warehouseId, pageSize, page*pageSize)

	return count, result, err
}

func (repo *Repo) ViewInventoryItemByProduct(
	ctx context.Context,
	warehouseId uuid.UUID,
	productId int,
	page, pageSize int,
	sortedBy, sortOrder string) (int, []InventoryItem, error) {

	log.Println("ViewInventoryItemByProduct",
		warehouseId, productId,
		page, pageSize, sortedBy, sortOrder)

	var count int
	result := make([]InventoryItem, 0)

	query := `select count(*) from inventory_item
        where warehouse_id = ? and product_id = ?`
	query = repo.db.Rebind(query)
	err := repo.db.GetContext(ctx, &count, query, warehouseId, productId)
	if err != nil {
		return count, result, err
	}

	query = `select i.id, i.product_id,
        p.name as product_name,
        i.warehouse_id, i.quantity, i.unit_cost,
        i.currency_uom_id,
        i.created_at, i.updated_at
        from inventory_item i 
        inner join product p
            on p.id = i.product_id
        where i.warehouse_id = ? and product_id = ?
        order by i.%s %s
        limit ? offset ?`
	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)

	err = repo.db.SelectContext(ctx, &result, query,
		warehouseId, productId, pageSize, page*pageSize)

	return count, result, err
}
