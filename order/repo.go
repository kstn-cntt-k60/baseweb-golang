package order

import (
	"context"
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

func (repo *Repo) ViewCustomerStoreByCustomer(
	ctx context.Context,
	customerId uuid.UUID,
	page, pageSize int,
	sortedBy, sortOrder string) (int, []CustomerStore, error) {

	log.Println("ViewCustomerStoreByCustomer", customerId,
		page, pageSize, sortedBy, sortOrder)

	var count int
	result := make([]CustomerStore, 0)

	query := `select count(*) from facility f
        inner join facility_customer fc on fc.id = f.id
        where fc.customer_id = ?`
	query = repo.db.Rebind(query)
	err := repo.db.GetContext(ctx, &count, query, customerId)
	if err != nil {
		return count, result, err
	}

	query = `select f.id, f.name, c.name as customer_name,
        fc.customer_id, f.address, f.created_at, f.updated_at
        from facility f
            inner join facility_customer fc on fc.id = f.id
            inner join customer c on fc.customer_id = c.id
        where fc.customer_id = ?
        order by f.%s %s
        offset ? limit ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, customerId,
		page*pageSize, pageSize)

	return count, result, err
}

func (repo *Repo) SelectCustomerStoreByCustomer(
	ctx context.Context,
	customerId uuid.UUID) ([]CustomerStore, error) {

	log.Println("SelectCustomerStoreByCustomer", customerId)

	result := make([]CustomerStore, 0)

	query := `select f.id, f.name, c.name as customer_name,
        fc.customer_id, f.address, f.created_at, f.updated_at
        from facility f
            inner join facility_customer fc on fc.id = f.id
            inner join customer c on fc.customer_id = c.id
        where fc.customer_id = ?`

	query = repo.db.Rebind(query)
	err := repo.db.SelectContext(ctx, &result, query, customerId)
	return result, err
}

var quantityAvailableErr error = errors.New("quantity available exceeded")

func (repo *Repo) AddOrder(
	ctx context.Context,
	customerId, warehouseId uuid.UUID,
	products []ClientProduct,
	address string,
	customerStoreId *uuid.UUID,
	userLoginId uuid.UUID) error {

	log.Println("AddOrder", customerId, warehouseId, products,
		address, customerStoreId)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `insert into sale_order(
        customer_id, original_warehouse_id,
        created_by_user_login_id, ship_to_address, 
        ship_to_facility_customer_id,
        sale_order_status_id)
        values (?, ?, ?, ?, ?, 1) returning id`
	query = repo.db.Rebind(query)

	var orderId int64
	err = tx.GetContext(ctx, &orderId, query,
		customerId, warehouseId, userLoginId, address, customerStoreId)
	if err != nil {
		return err
	}

	now := time.Now()

	priceQuery := `select id from product_price
        where product_id = ?
            and effective_from <= ? 
            and (expired_at is null or ? < expired_at)`
	priceQuery = repo.db.Rebind(priceQuery)

	query = `insert into sale_order_item(
        sale_order_id, sale_order_seq,
        product_price_id, quantity)
        values (?, ?, ?, ?)`
	query = repo.db.Rebind(query)

	availableQuery := `select quantity_available
        from warehouse_product_statistics
        where product_id = ? and warehouse_id = ?`
	availableQuery = repo.db.Rebind(availableQuery)

	updateAvailableQuery := `update warehouse_product_statistics
        set quantity_available = ?
        where product_id = ? and warehouse_id = ?`
	updateAvailableQuery = repo.db.Rebind(updateAvailableQuery)

	for index, product := range products {
		var priceId uuid.UUID
		err = tx.GetContext(ctx, &priceId, priceQuery, product.Id, now, now)
		if err != nil {
			return err
		}

		var quantityAvailable decimal.Decimal
		err = tx.GetContext(ctx, &quantityAvailable, availableQuery,
			product.Id, warehouseId)
		if err != nil {
			return err
		}

		if product.Quantity.GreaterThan(quantityAvailable) {
			return quantityAvailableErr
		}

		_, err = tx.ExecContext(ctx, updateAvailableQuery,
			quantityAvailable.Sub(product.Quantity),
			product.Id, warehouseId)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, query, orderId, index, priceId, product.Quantity)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *Repo) ViewProductInfoByWarehouse(
	ctx context.Context, warehouseId uuid.UUID,
	page, pageSize int,
	sortedBy, sortOrder string) (int, []ProductInfo, error) {

	log.Println("ViewProductInfoByWarehouse", warehouseId,
		page, pageSize, sortedBy, sortOrder)

	var count int
	result := make([]ProductInfo, 0)
	now := time.Now()

	query := `select count(p.id)
        from product p
           inner join warehouse_product_statistics s on s.product_id = p.id
           inner join product_price pp on pp.product_id = p.id
        where
           s.warehouse_id = ?
           and pp.effective_from <= ?
           and (pp.expired_at is null or ? < pp.expired_at)`
	query = repo.db.Rebind(query)
	err := repo.db.GetContext(ctx, &count, query, warehouseId, now, now)
	if err != nil {
		return count, result, err
	}

	query = `select p.id, p.name, u.username as created_by,
        p.weight, p.weight_uom_id, p.unit_uom_id,
        p.created_at, p.updated_at,
        pp.price, pp.currency_uom_id, pp.effective_from,
        s.quantity_available
        from product p
           inner join warehouse_product_statistics s on s.product_id = p.id
           inner join product_price pp on pp.product_id = p.id
           inner join user_login u on u.id = p.created_by_user_login_id
        where
           s.warehouse_id = ?
           and pp.effective_from <= ?
           and (pp.expired_at is null or ? < pp.expired_at)
        order by p.%s %s
        offset ? limit ?`
	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query,
		warehouseId, now, now, page*pageSize, pageSize)

	return count, result, err
}

func (repo *Repo) SelectProductInfoByWarehouse(
	ctx context.Context,
	warehouseId uuid.UUID) ([]ProductInfo, error) {

	log.Println("SelectProductInfoByWarehouse", warehouseId)

	result := make([]ProductInfo, 0)
	now := time.Now()

	query := `select p.id, p.name, u.username as created_by,
        p.weight, p.weight_uom_id, p.unit_uom_id,
        p.created_at, p.updated_at,
        pp.price, pp.currency_uom_id, pp.effective_from,
        s.quantity_available
        from product p
           inner join warehouse_product_statistics s on s.product_id = p.id
           inner join product_price pp on pp.product_id = p.id
           inner join user_login u on u.id = p.created_by_user_login_id
        where
           s.warehouse_id = ?
           and pp.effective_from <= ?
           and (pp.expired_at is null or ? < pp.expired_at)`
	query = repo.db.Rebind(query)
	err := repo.db.SelectContext(ctx, &result, query,
		warehouseId, now, now)

	return result, err
}
