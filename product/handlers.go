package product

import (
	"baseweb/basic"
	"baseweb/security"
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type Root struct {
	repo *Repo
}

func InitRoot(repo *Repo) *Root {
	return &Root{
		repo: repo,
	}
}

func (root *Root) AddProductHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	product := Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		return err
	}
	product.CreatedBy = userLogin.Id

	err = root.repo.InsertProduct(ctx, product)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) ViewProductHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	queries := r.URL.Query()

	var err error

	var page int
	page, err = strconv.Atoi(queries.Get("page"))
	if err != nil {
		page = 0
	}

	var pageSize int
	pageSize, err = strconv.Atoi(queries.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedBy := "created_at"
	sortedByQuery := queries.Get("sortedBy")
	if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	} else if sortedByQuery == "name" {
		sortedBy = "name"
	}

	sortOrder := "desc"
	sortOrderQuery := queries.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	search := queries.Get("query")

	var count uint
	var products []ClientProduct

	count, products, err = ViewProductFuzzy(ctx, root.repo,
		uint(page), uint(pageSize),
		sortedBy, sortOrder, search)

	type Response struct {
		ProductList  []ClientProduct `json:"productList"`
		ProductCount uint            `json:"productCount"`
	}

	res := Response{
		ProductList:  products,
		ProductCount: count,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) UpdateProductHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	product := Product{}
	err := json.NewDecoder(r.Body).Decode(&product)
	if err != nil {
		return err
	}

	err = root.repo.UpdateProduct(ctx, product)
	if err != nil {
		return err
	}

	p, err := root.repo.GetProduct(ctx, product.Id)

	return json.NewEncoder(w).Encode(p)
}

func (root *Root) DeleteProductHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id int64 `json:"id"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.DeleteProduct(ctx, req.Id)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) ViewProductPricingHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	queries := r.URL.Query()

	var err error

	var page int
	page, err = strconv.Atoi(queries.Get("page"))
	if err != nil {
		page = 0
	}

	var pageSize int
	pageSize, err = strconv.Atoi(queries.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedBy := "created_at"
	sortedByQuery := queries.Get("sortedBy")
	if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	} else if sortedByQuery == "name" {
		sortedBy = "name"
	}

	sortOrder := "desc"
	sortOrderQuery := queries.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	search := queries.Get("query")

	var count uint
	var products []ClientProduct
	count, products, err = ViewProductFuzzy(ctx, root.repo,
		uint(page), uint(pageSize),
		sortedBy, sortOrder, search)

	idList := make([]int64, len(products))
	for i := 0; i < len(products); i++ {
		idList[i] = products[i].Id
	}

	prices, err := root.repo.SelectProductPriceFromIdList(
		ctx, idList, time.Now())
	if err != nil {
		return err
	}

	type Response struct {
		ProductList  []ClientProduct `json:"productList"`
		ProductCount uint            `json:"productCount"`
		PriceList    []ProductPrice  `json:"priceList"`
	}

	res := Response{
		ProductList:  products,
		ProductCount: count,
		PriceList:    prices,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) ViewProductPriceHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	queries := r.URL.Query()

	var err error

	productId, err := strconv.Atoi(queries.Get("id"))
	if err != nil {
		return err
	}

	var page int
	page, err = strconv.Atoi(queries.Get("page"))
	if err != nil {
		page = 0
	}

	var pageSize int
	pageSize, err = strconv.Atoi(queries.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedBy := "created_at"
	sortedByQuery := queries.Get("sortedBy")
	if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	} else if sortedByQuery == "effectiveFrom" {
		sortedBy = "effective_from"
	} else if sortedByQuery == "expiredAt" {
		sortedBy = "expired_at"
	}

	sortOrder := "desc"
	sortOrderQuery := queries.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	count, prices, err := root.repo.ViewProductPrice(ctx,
		int64(productId), uint(page), uint(pageSize), sortedBy, sortOrder)

	type Response struct {
		PriceList  []ProductPrice `json:"priceList"`
		PriceCount uint           `json:"priceCount"`
	}

	res := Response{
		PriceList:  prices,
		PriceCount: count,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) AddProductPriceHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	userLogin := ctx.Value("userLogin").(security.UserLogin)

	price := InsertionPrice{}
	json.NewDecoder(r.Body).Decode(&price)
	price.CreatedBy = userLogin.Id

	err := root.repo.InsertProductPrice(ctx, price)
	if err != nil {
		return err
	}

	type Response struct {
		ProductId int64 `json:"productId"`
	}

	res := Response{
		ProductId: price.ProductId,
	}

	return json.NewEncoder(w).Encode(res)
}
