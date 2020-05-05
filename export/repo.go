package export

import (
	"baseweb/order"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/shopspring/decimal"
)

type Repo struct {
	db *sqlx.DB
}

func InitRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

var ErrExported = errors.New("sale_order_item exported")

func (repo *Repo) ExportSaleOrderItem(
	ctx context.Context, saleOrderId int64,
	saleOrderSeq int, effectiveFrom time.Time) error {

	log.Println("ExportSaleOrderItem", saleOrderId,
		saleOrderSeq, effectiveFrom)

	tx, err := repo.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelRepeatableRead,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback()

	exported := false
	query := repo.db.Rebind(`
        select exported from sale_order_item
        where sale_order_id = ? and sale_order_seq = ?
        `)
	err = tx.GetContext(ctx, &exported, query,
		saleOrderId, saleOrderSeq)
	if err != nil {
		return err
	}

	if exported {
		return ErrExported
	}

	type OrderItemInfo struct {
		ProductId   int64           `db:"product_id"`
		WarehouseId uuid.UUID       `db:"original_warehouse_id"`
		Quantity    decimal.Decimal `db:"quantity"`
	}
	info := OrderItemInfo{}

	query = repo.db.Rebind(`
        select pp.product_id, o.original_warehouse_id, oi.quantity
        from product_price pp
        inner join sale_order_item oi on oi.product_price_id = pp.id
        inner join sale_order o on o.id = oi.sale_order_id
        where oi.sale_order_id = ? and oi.sale_order_seq = ?
        `)
	err = tx.GetContext(ctx, &info, query, saleOrderId, saleOrderSeq)
	if err != nil {
		return err
	}

	type InventoryItem struct {
		Id             int64           `db:"id"`
		QuantityOnHand decimal.Decimal `db:"quantity_on_hand"`
	}
	items := make([]InventoryItem, 0)
	query = repo.db.Rebind(`
        select id, quantity_on_hand
        from inventory_item
        where product_id = ? and warehouse_id = ?
        and quantity_on_hand > 0
        `)
	err = tx.SelectContext(ctx, &items, query,
		info.ProductId, info.WarehouseId)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(`
        insert into inventory_item_detail(
        inventory_item_id, exported_quantity,
        effective_from,
        sale_order_id, sale_order_seq)
        values(?, ?, ?, ?, ?)
        `)

	updateQuery := repo.db.Rebind(`
        update inventory_item
        set quantity_on_hand = quantity_on_hand - ?
        where id = ?
        `)

	quantity := info.Quantity
	for _, item := range items {
		// minimum
		var exportedQuantity decimal.Decimal
		if quantity.GreaterThan(item.QuantityOnHand) {
			exportedQuantity = item.QuantityOnHand
		} else {
			exportedQuantity = quantity
		}

		_, err = tx.ExecContext(ctx, query,
			item.Id, exportedQuantity,
			effectiveFrom,
			saleOrderId, saleOrderSeq,
		)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, updateQuery, exportedQuantity, item.Id)
		if err != nil {
			return err
		}

		quantity = quantity.Sub(exportedQuantity)
		if quantity.Equal(decimal.Zero) {
			break
		}
	}

	query = repo.db.Rebind(`
        update warehouse_product_statistics
        set quantity_on_hand = quantity_on_hand - ?
        where warehouse_id = ? and product_id = ?
        `)
	_, err = tx.ExecContext(ctx, query,
		info.Quantity, info.WarehouseId, info.ProductId)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(`
        update sale_order_item
        set exported = TRUE
        where sale_order_id = ? and sale_order_seq = ?
        `)
	_, err = tx.ExecContext(ctx, query,
		saleOrderId, saleOrderSeq)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(`
        update sale_order
        set sale_order_status_id = 3
        where id = ?
        `)
	_, err = tx.ExecContext(ctx, query, saleOrderId)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *Repo) ViewExportableSalesOrder(
	ctx context.Context,
	page, pageSize int,
	sortedBy, sortOrder string,
) (int, []order.SaleOrder, error) {
	log.Println("ViewExportableSalesOrder", page,
		pageSize, sortedBy, sortOrder)

	var count int
	orders := make([]order.SaleOrder, 0)
	var err error

	countQuery := repo.db.Rebind(`
            select count(*) from sale_order
            where sale_order_status_id = 2 or sale_order_status_id = 3
            `)
	err = repo.db.GetContext(ctx, &count, countQuery)
	if err != nil {
		return count, orders, err
	}

	query := `select o.id, c.name as customer,
        fw.name as warehouse, u.username as created_by,
        o.ship_to_address,
        coalesce(fc.name, '') as customer_store,
        o.sale_order_status_id,
        o.created_at, o.updated_at
        from sale_order o
            inner join customer c on c.id = o.customer_id
            inner join user_login u on u.id = o.created_by_user_login_id
            inner join facility_warehouse w on w.id = o.original_warehouse_id
            inner join facility fw on fw.id = w.id
            left join (
                select f.id, f.name from facility f
                inner join facility_customer fc on fc.id = f.id
            ) fc on fc.id = o.ship_to_facility_customer_id
        where sale_order_status_id = 2 or sale_order_status_id = 3
        order by o.%s %s
        offset ? limit ?`
	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)

	err = repo.db.SelectContext(ctx, &orders, query,
		page*pageSize, pageSize)

	return count, orders, err
}

func (repo *Repo) CompleteSalesOrder(ctx context.Context, id int64) error {
	log.Println("CompleteSalesOrder", id)

	query := repo.db.Rebind(`
        update sale_order
        set sale_order_status_id = 4
        where id = ? and sale_order_status_id = 3
        `)
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *Repo) ViewCompletedSalesOrder(
	ctx context.Context,
	page, pageSize int,
	sortedBy, sortOrder string,
) (int, []order.SaleOrder, error) {
	log.Println("ViewCompletedSalesOrder", page,
		pageSize, sortedBy, sortOrder)

	var count int
	orders := make([]order.SaleOrder, 0)
	var err error

	countQuery := repo.db.Rebind(`
            select count(*) from sale_order
            where sale_order_status_id = 4
            `)
	err = repo.db.GetContext(ctx, &count, countQuery)
	if err != nil {
		return count, orders, err
	}

	query := `select o.id, c.name as customer,
        fw.name as warehouse, u.username as created_by,
        o.ship_to_address,
        coalesce(fc.name, '') as customer_store,
        o.sale_order_status_id,
        o.created_at, o.updated_at
        from sale_order o
            inner join customer c on c.id = o.customer_id
            inner join user_login u on u.id = o.created_by_user_login_id
            inner join facility_warehouse w on w.id = o.original_warehouse_id
            inner join facility fw on fw.id = w.id
            left join (
                select f.id, f.name from facility f
                inner join facility_customer fc on fc.id = f.id
            ) fc on fc.id = o.ship_to_facility_customer_id
        where sale_order_status_id = 4
        order by o.%s %s
        offset ? limit ?`
	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)

	err = repo.db.SelectContext(ctx, &orders, query,
		page*pageSize, pageSize)

	return count, orders, err
}
