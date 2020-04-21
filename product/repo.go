package product

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db            *sqlx.DB
	productCount  *sqlx.Stmt
	viewProduct   *sqlx.Stmt
	selectProduct *sqlx.Stmt
	getProduct    *sqlx.Stmt
	priceCount    *sqlx.Stmt
	viewPrice     *sqlx.Stmt
}

func InitRepo(db *sqlx.DB) *Repo {
	query := "select count(*) from product"
	productCount, err := db.Preparex(query)
	if err != nil {
		log.Fatal(err)
	}

	query = `select p.id, p.name, 
        p.weight, p.weight_uom_id, p.unit_uom_id,
        p.description, p.created_at, p.updated_at,
        u.username as created_by
        from product p
        inner join user_login u
            on u.id = p.created_by_user_login_id
        order by p.created_at desc
        limit ? offset ?`
	viewProduct, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Fatal(err)
	}

	query = `select p.id, p.name, 
        p.weight, p.weight_uom_id, p.unit_uom_id,
        p.description, p.created_at, p.updated_at,
        u.username as created_by
        from product p
        inner join user_login u
            on u.id = p.created_by_user_login_id`
	selectProduct, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = `select p.id, p.name, 
        p.weight, p.weight_uom_id, p.unit_uom_id,
        p.description, p.created_at, p.updated_at,
        u.username as created_by
        from product p
        inner join user_login u
            on u.id = p.created_by_user_login_id
        where p.id = ?`
	getProduct, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = "select count(*) from product_price where product_id = ?"
	priceCount, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	query = `select p.id, p.product_id, p.price, 
        p.currency_uom_id, p.effective_from, p.expired_at,
        u.username as created_by,
        p.created_at, p.updated_at
        from product_price p
        inner join user_login u
            on u.id = p.created_by_user_login_id
        where p.product_id = ?
        order by p.created_at desc
        limit ? offset ?`
	viewPrice, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panic(err)
	}

	return &Repo{
		db:            db,
		productCount:  productCount,
		viewProduct:   viewProduct,
		selectProduct: selectProduct,
		getProduct:    getProduct,
		priceCount:    priceCount,
		viewPrice:     viewPrice,
	}
}

func (repo *Repo) InsertProduct(
	ctx context.Context, product Product) error {

	log.Println("InsertProduct", product.Name,
		product.Description, product.Weight,
		product.WeightUomId, product.UnitUomId)

	query := `insert into product(
        name, created_by_user_login_id,
        description, weight,
        weight_uom_id, unit_uom_id)
        values(:name, :created_by_user_login_id,
        :description, :weight,
        :weight_uom_id, :unit_uom_id)`

	_, err := repo.db.NamedExecContext(ctx, query, product)
	return err
}

func (repo *Repo) ViewProduct(
	ctx context.Context, page,
	pageSize uint, sortedBy, sortOrder string) (uint, []ClientProduct, error) {

	log.Println("ViewProduct", page, pageSize, sortedBy, sortOrder)

	var count uint
	result := make([]ClientProduct, 0)

	err := repo.productCount.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "created_at" && sortOrder == "desc" {
		err = repo.viewProduct.SelectContext(ctx, &result, pageSize, page*pageSize)
		return count, result, err
	} else {
		query := fmt.Sprintf(`select p.id, p.name, 
            p.weight, p.weight_uom_id, p.unit_uom_id,
            p.description, p.created_at, p.updated_at,
            u.username as created_by
            from product p
            inner join user_login u
                on u.id = p.created_by_user_login_id
            order by p.%s %s
            limit ? offset ?`, sortedBy, sortOrder)
		log.Println("[SQL", query)

		err = repo.db.SelectContext(ctx, &result,
			repo.db.Rebind(query), pageSize, page*pageSize)
		return count, result, err
	}
}

func (repo *Repo) SelectProduct(
	ctx context.Context) (uint, []ClientProduct, error) {

	log.Println("SelectProduct")

	var count uint
	result := make([]ClientProduct, 0)

	err := repo.productCount.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	err = repo.selectProduct.SelectContext(ctx, &result)
	return count, result, err
}

func (repo *Repo) UpdateProduct(
	ctx context.Context, product Product) error {

	log.Println("UpdateProduct", product.Id, product.Name,
		product.Description, product.Weight,
		product.WeightUomId, product.UnitUomId)

	query := `update product set name = :name,
        description = :description, weight = :weight, 
        weight_uom_id = :weight_uom_id,
        unit_uom_id = :unit_uom_id where id = :id`
	_, err := repo.db.NamedExecContext(ctx, query, product)
	return err
}

func (repo *Repo) GetProduct(
	ctx context.Context, id int64) (ClientProduct, error) {

	log.Println("GetProduct", id)

	product := ClientProduct{}
	return product, repo.getProduct.GetContext(ctx, &product, id)
}

func (repo *Repo) DeleteProduct(
	ctx context.Context, id int64) error {

	log.Println("DeleteProduct", id)

	query := "delete from product where id = ?"
	_, err := repo.db.ExecContext(ctx, repo.db.Rebind(query), id)
	return err
}

func (repo *Repo) SelectProductPriceFromIdList(
	ctx context.Context, productIdList []int64,
	now time.Time) ([]ProductPrice, error) {

	log.Println("SelectProductPriceFromIdList", productIdList)

	result := make([]ProductPrice, 0)
	if len(productIdList) == 0 {
		return result, nil
	}

	query := `select p.id, p.product_id, p.price, 
        p.currency_uom_id, p.effective_from, p.expired_at,
        u.username as created_by,
        p.created_at, p.updated_at
        from product_price p
        inner join user_login u
            on u.id = p.created_by_user_login_id
        where product_id in (?)
            and effective_from <= ? 
            and (expired_at is null or ? < expired_at)`
	query, args, err := sqlx.In(query, productIdList, now, now)
	if err != nil {
		return result, err
	}
	log.Println(query)

	return result, repo.db.SelectContext(ctx,
		&result, repo.db.Rebind(query), args...)
}

func (repo *Repo) ViewProductPrice(
	ctx context.Context, productId int64,
	page, pageSize uint,
	sortedBy, sortOrder string) (uint, []ProductPrice, error) {

	log.Println("ViewProductPrice", productId,
		page, pageSize, sortedBy, sortOrder)

	var count uint
	result := make([]ProductPrice, 0)

	err := repo.priceCount.GetContext(ctx, &count, productId)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "created_at" && sortOrder == "desc" {
		err = repo.viewPrice.SelectContext(ctx,
			&result, productId, pageSize, page*pageSize)
		return count, result, err
	} else {
		query := `select p.id, p.product_id, p.price, 
            p.currency_uom_id, p.effective_from, p.expired_at,
            u.username as created_by,
            p.created_at, p.updated_at
            from product_price p
            inner join user_login u
                on u.id = p.created_by_user_login_id
            where p.product_id = ?
            order by p.%s %s
            limit ? offset ?`
		query = fmt.Sprintf(query, sortedBy, sortOrder)
		log.Println("[SQL", query)
		query = repo.db.Rebind(query)

		err = repo.db.SelectContext(ctx, &result, query,
			productId, pageSize, page*pageSize)
		return count, result, err
	}
}

func (repo *Repo) InsertProductPrice(
	ctx context.Context, price InsertionPrice) error {

	log.Println("InsertProductPrice", price.ProductId, price.Price,
		price.CurrencyUomId, price.EffectiveFrom, price.CreatedBy)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `update product_price 
        set expired_at = ?
        where (expired_at is null or ? < expired_at)
            and product_id = ?`
	query = repo.db.Rebind(query)
	_, err = tx.ExecContext(ctx, query, price.EffectiveFrom,
		price.EffectiveFrom, price.ProductId)
	if err != nil {
		return err
	}

	query = `insert into product_price(
        product_id, price, currency_uom_id,
        created_by_user_login_id,
        effective_from)
        values (:product_id, :price, :currency_uom_id, 
            :created_by_user_login_id, :effective_from)`
	_, err = tx.NamedExecContext(ctx, query, price)
	if err != nil {
		return err
	}

	return tx.Commit()
}
