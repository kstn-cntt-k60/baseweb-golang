package facility

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB

	countWarehouse  *sqlx.Stmt
	viewWarehouse   *sqlx.Stmt
	selectWarehouse *sqlx.Stmt

	countCustomerStore  *sqlx.Stmt
	viewCustomerStore   *sqlx.Stmt
	selectCustomerStore *sqlx.Stmt
}

func InitRepo(db *sqlx.DB) *Repo {
	query := "select count(*) from facility where facility_type_id = 1"
	countWarehouse, err := db.Preparex(query)
	if err != nil {
		log.Panic(err)
	}

	query = `select id, name, address, created_at, updated_at
        from facility where facility_type_id = 1
        order by created_at desc
        limit ? offset ?`
	viewWarehouse, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = `select id, name, address, created_at, updated_at
        from facility where facility_type_id = 1`
	selectWarehouse, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = "select count(*) from facility where facility_type_id = 2"
	countCustomerStore, err := db.Preparex(query)
	if err != nil {
		log.Panic(err)
	}

	query = `select f.id, f.name, f.address,
        c.name as customer_name,
        f.created_at, f.updated_at
        from facility f 
        inner join facility_customer fc on fc.id = f.id
        inner join customer c on c.id = fc.customer_id
        where f.facility_type_id = 2
        order by f.created_at desc
        limit ? offset ?`
	viewCustomerStore, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = `select f.id, f.name, f.address,
        c.name as customer_name,
        f.created_at, f.updated_at
        from facility f 
        inner join facility_customer fc on fc.id = f.id
        inner join customer c on c.id = fc.customer_id
        where f.facility_type_id = 2`
	selectCustomerStore, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	return &Repo{
		db:                  db,
		countWarehouse:      countWarehouse,
		viewWarehouse:       viewWarehouse,
		selectWarehouse:     selectWarehouse,
		countCustomerStore:  countCustomerStore,
		viewCustomerStore:   viewCustomerStore,
		selectCustomerStore: selectCustomerStore,
	}
}

func (repo *Repo) InsertWarehouse(
	ctx context.Context, warehouse Warehouse) error {

	log.Println("InsertWarehouse", warehouse.Name, warehouse.Address)

	query := `insert into facility(
        name, facility_type_id, address)
        values (:name, 1, :address)`

	_, err := repo.db.NamedExecContext(ctx, query, warehouse)
	return err
}

func (repo *Repo) ViewWarehouse(
	ctx context.Context, page, pageSize uint,
	sortedBy, sortOrder string) (uint, []Warehouse, error) {

	log.Println("ViewWarehouse", page, pageSize,
		sortedBy, sortOrder)

	var count uint
	result := make([]Warehouse, 0)

	err := repo.countWarehouse.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "created_at" && sortOrder == "desc" {
		err = repo.viewWarehouse.SelectContext(ctx,
			&result, pageSize, page*pageSize)
		return count, result, err
	} else {
		query := `select id, name, address, created_at, updated_at
            from facility where facility_type_id = 1
            order by %s %s
            limit ? offset ?`
		query = fmt.Sprintf(query, sortedBy, sortOrder)
		log.Println("[SQL]", query)
		query = repo.db.Rebind(query)

		err = repo.db.SelectContext(ctx, &result,
			query, pageSize, page*pageSize)
		return count, result, err
	}
}

func (repo *Repo) SelectWarehouse(
	ctx context.Context) (uint, []Warehouse, error) {

	log.Println("SelectWarehouse")

	var count uint
	result := make([]Warehouse, 0)

	err := repo.countWarehouse.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	err = repo.selectWarehouse.SelectContext(ctx, &result)
	return count, result, err
}

func (repo *Repo) UpdateWarehouse(
	ctx context.Context, warehouse Warehouse) error {

	log.Println("UpdateWarehouse", warehouse.Id,
		warehouse.Name, warehouse.Address)

	query := `update facility set name = :name,
        address = :address where id = :id`
	_, err := repo.db.NamedExecContext(ctx, query, warehouse)
	return err
}

func (repo *Repo) DeleteWarehouse(
	ctx context.Context, id uuid.UUID) error {

	log.Println("DeleteWarehouse", id)

	query := "delete from facility where id = ?"
	query = repo.db.Rebind(query)

	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *Repo) ViewCustomerStore(
	ctx context.Context, page, pageSize uint,
	sortedBy, sortOrder string) (uint, []CustomerStore, error) {

	log.Println("ViewCustomerStore", page, pageSize,
		sortedBy, sortOrder)

	var count uint
	result := make([]CustomerStore, 0)

	err := repo.countCustomerStore.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "created_at" && sortOrder == "desc" {
		err = repo.viewCustomerStore.SelectContext(ctx,
			&result, pageSize, page*pageSize)
		return count, result, err
	} else {

		query := `select f.id, f.name, f.address,
            c.name as customer_name,
            f.created_at, f.updated_at
            from facility f 
            inner join facility_customer fc on fc.id = f.id
            inner join customer c on c.id = fc.customer_id
            where f.facility_type_id = 2
            order by f.%s %s
            limit ? offset ?`
		query = fmt.Sprintf(query, sortedBy, sortOrder)
		log.Println("[SQL]", query)
		query = repo.db.Rebind(query)

		err = repo.db.SelectContext(ctx, &result,
			query, pageSize, page*pageSize)
		return count, result, err
	}
}

func (repo *Repo) SelectCustomerStore(
	ctx context.Context) (uint, []CustomerStore, error) {

	log.Println("SelectCustomerStore")

	var count uint
	result := make([]CustomerStore, 0)

	err := repo.countCustomerStore.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	err = repo.selectCustomerStore.SelectContext(ctx, &result)
	return count, result, err
}

func (repo *Repo) SelectSimpleCustomer(
	ctx context.Context) ([]SimpleCustomer, error) {

	log.Println("SelectSimpleCustomer")

	result := make([]SimpleCustomer, 0)
	query := "select id, name from customer"
	return result, repo.db.SelectContext(ctx, &result, query)
}

func (repo *Repo) InsertCustomerStore(
	ctx context.Context, store InsertStore) error {

	log.Println("InsertCustomerStore", store.Name,
		store.Address, store.CustomerId)

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	store.Id = id

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `insert into facility(
        id, name, facility_type_id, address)
        values (:id, :name, 2, :address)`

	_, err = tx.NamedExecContext(ctx, query, store)
	if err != nil {
		return err
	}

	query = `insert into facility_customer(
        id, customer_id)
        values (:id, :customer_id)`
	_, err = tx.NamedExecContext(ctx, query, store)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *Repo) UpdateCustomerStore(
	ctx context.Context, store CustomerStore) error {

	log.Println("UpdateCustomerStore", store.Id,
		store.Name, store.Address)

	query := `update facility set name = :name,
        address = :address where id = :id`
	_, err := repo.db.NamedExecContext(ctx, query, store)
	return err
}

func (repo *Repo) DeleteCustomerStore(
	ctx context.Context, id uuid.UUID) error {

	log.Println("DeleteCustomerStore", id)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "delete from facility_customer where id = ?"
	query = repo.db.Rebind(query)

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	query = "delete from facility where id = ?"
	query = repo.db.Rebind(query)

	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
