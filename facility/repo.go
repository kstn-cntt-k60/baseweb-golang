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
	query := "select count(*) from facility_warehouse"
	countWarehouse, err := db.Preparex(query)
	if err != nil {
		log.Panic(err)
	}

	query = `select f.id, f.name, f.address,
        f.created_at, f.updated_at
        from facility f 
        inner join facility_warehouse fw on fw.id = f.id
        order by created_at desc
        limit ? offset ?`
	viewWarehouse, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = `select f.id, f.name, f.address,
        f.created_at, f.updated_at
        from facility f
        inner join facility_warehouse fw on fw.id = f.id`
	selectWarehouse, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = "select count(*) from facility_customer"
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
        inner join customer c on c.id = fc.customer_id`
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

	id, err := uuid.NewUUID()
	if err != nil {
		return err
	}
	warehouse.Id = id

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `insert into facility(
        id, name, facility_type_id, address)
        values (:id, :name, 1, :address)`
	_, err = tx.NamedExecContext(ctx, query, warehouse)
	if err != nil {
		return err
	}

	query = `insert into facility_warehouse(id) values (:id)`
	_, err = tx.NamedExecContext(ctx, query, warehouse)
	if err != nil {
		return err
	}

	return tx.Commit()
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
		query := `select f.id, f.name, f.address,
            f.created_at, f.updated_at
            from facility f
            inner join facility_warehouse fw on fw.id = f.id
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

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := "delete from facility_warehouse where id = ?"
	query = repo.db.Rebind(query)
	_, err = repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	query = "delete from facility where id = ?"
	query = repo.db.Rebind(query)
	_, err = repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
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

func (repo *Repo) GetWarehouse(
	ctx context.Context, id uuid.UUID) (Warehouse, error) {

	log.Println("GetWarehouse", id)

	query := `select f.id, f.name, f.address,
        f.created_at, f.updated_at
        from facility f
        inner join facility_warehouse fw on fw.id = f.id
        where f.id = ?`
	query = repo.db.Rebind(query)

	warehouse := Warehouse{}
	err := repo.db.GetContext(ctx, &warehouse, query, id)

	return warehouse, err
}
