package facility

import (
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

func (root *Root) AddWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	warehouse := Warehouse{}
	err := json.NewDecoder(r.Body).Decode(&warehouse)
	if err != nil {
		return err
	}

	err = root.repo.InsertWarehouse(ctx, warehouse)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) ViewWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	queries := r.URL.Query()

	var err error

	page, err := strconv.Atoi(queries.Get("page"))
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.Atoi(queries.Get("pageSize"))
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
	var warehouses []Warehouse

	if search == "" {
		count, warehouses, err = root.repo.ViewWarehouse(ctx,
			uint(page), uint(pageSize), sortedBy, sortOrder)
		if err != nil {
			return err
		}
	} else {
		count, warehouses, err = root.repo.SelectWarehouse(ctx)
		if err != nil {
			return err
		}
		count, warehouses = FuzzySearchWarehouse(warehouses,
			uint(page), uint(pageSize), search)
	}

	type Response struct {
		WarehouseList  []Warehouse `json:"warehouseList"`
		WarehouseCount uint        `json:"warehouseCount"`
	}

	res := Response{
		WarehouseList:  warehouses,
		WarehouseCount: count,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) UpdateWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	warehouse := Warehouse{}
	err := json.NewDecoder(r.Body).Decode(&warehouse)
	if err != nil {
		return err
	}

	err = root.repo.UpdateWarehouse(ctx, warehouse)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) DeleteWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	warehouse := Warehouse{}
	err := json.NewDecoder(r.Body).Decode(&warehouse)
	if err != nil {
		return err
	}

	err = root.repo.DeleteWarehouse(ctx, warehouse.Id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}
