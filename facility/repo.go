package facility

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db              *sqlx.DB
	countWarehouse  *sqlx.Stmt
	viewWarehouse   *sqlx.Stmt
	selectWarehouse *sqlx.Stmt
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

	return &Repo{
		db:              db,
		countWarehouse:  countWarehouse,
		viewWarehouse:   viewWarehouse,
		selectWarehouse: selectWarehouse,
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

	log.Println("ViewWarehouse", page, pageSize)

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
