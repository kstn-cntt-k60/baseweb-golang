package product

import (
	"baseweb/security"
	"encoding/json"
	"net/http"
	"strconv"
)

type Root struct {
	repo *Repo
}

func InitRoot(repo *Repo) *Root {
	return &Root{
		repo: repo,
	}
}

type OkResponse struct {
	Status string `json:"status"`
}

var okResponse = OkResponse{
	Status: "ok",
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

	return json.NewEncoder(w).Encode(okResponse)
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

	if search == "" {
		count, products, err = root.repo.ViewProduct(
			ctx, uint(page), uint(pageSize), sortedBy, sortOrder)
		if err != nil {
			return err
		}
	} else {
		count, products, err = root.repo.SelectProduct(ctx)
		if err != nil {
			return err
		}
		count, products = FuzzySearchProduct(products,
			uint(page), uint(pageSize), search)
	}

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

	return json.NewEncoder(w).Encode(okResponse)
}
