package importProduct

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Root struct {
	repo *Repo
}

type OkResponse struct {
	Status string `json:"status"`
}

var okResponse = OkResponse{
	Status: "ok",
}

func InitRoot(repo *Repo) *Root {
	return &Root{
		repo: repo,
	}
}

func (root *Root) ViewProductByWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	query := r.URL.Query()
	warehouseId, err := uuid.Parse(query.Get("warehouseId"))
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedBy := "updated_at"
	sortedByQuery := query.Get("sortedBy")
	if sortedByQuery == "name" {
		sortedBy = "name"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	search := query.Get("query")

	var count uint
	var products []Product

	if search == "" {
		count, products, err = root.repo.ViewProductByWarehouse(
			ctx, warehouseId, page, pageSize, sortedBy, sortOrder)
		if err != nil {
			return err
		}
	} else {
		products, err = root.repo.SelectProductByWarehouse(ctx, warehouseId)
		if err != nil {
			return err
		}
		count, products = FuzzySearchProduct(products, page, pageSize, search)
	}

	type Response struct {
		Count       uint      `json:"count"`
		ProductList []Product `json:"productList"`
	}

	res := Response{
		Count:       count,
		ProductList: products,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) AddInventoryItemHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	item := InventoryItem{}
	err := json.NewDecoder(r.Body).Decode(&item)
	if err != nil {
		return err
	}

	err = root.repo.InsertInventoryItem(ctx, item)
	if err != nil {
		return err
	}

	type Response struct {
		WarehouseId uuid.UUID `json:"warehouseId"`
	}

	res := Response{
		WarehouseId: item.WarehouseId,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) ViewInventoryByWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	query := r.URL.Query()
	warehouseId, err := uuid.Parse(query.Get("warehouseId"))
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedBy := "created_at"
	sortedByQuery := query.Get("sortedBy")
	if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	} else if sortedByQuery == "name" {
		sortedBy = "name"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	var count int
	var items []InventoryItem

	count, items, err = root.repo.ViewInventoryItemByWarehouse(
		ctx, warehouseId, page, pageSize, sortedBy, sortOrder)
	if err != nil {
		return err
	}

	type Response struct {
		InventoryCount int             `json:"inventoryCount"`
		InventoryList  []InventoryItem `json:"inventoryList"`
	}

	res := Response{
		InventoryCount: count,
		InventoryList:  items,
	}

	return json.NewEncoder(w).Encode(res)

}

func (root *Root) ViewInventoryByProductHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	query := r.URL.Query()
	warehouseId, err := uuid.Parse(query.Get("warehouseId"))
	if err != nil {
		return err
	}

	productId, err := strconv.Atoi(query.Get("productId"))
	if err != nil {
		return err
	}

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedBy := "created_at"
	sortedByQuery := query.Get("sortedBy")
	if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	var count int
	var items []InventoryItem

	count, items, err = root.repo.ViewInventoryItemByProduct(
		ctx, warehouseId, productId, page, pageSize, sortedBy, sortOrder)
	if err != nil {
		return err
	}

	type Response struct {
		InventoryCount int             `json:"inventoryCount"`
		InventoryList  []InventoryItem `json:"inventoryList"`
	}

	res := Response{
		InventoryCount: count,
		InventoryList:  items,
	}

	return json.NewEncoder(w).Encode(res)

}
