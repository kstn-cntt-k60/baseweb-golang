package facility

import (
	"baseweb/basic"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Root struct {
	repo *Repo
}

func InitRoot(repo *Repo) *Root {
	return &Root{
		repo: repo,
	}
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

	return basic.ReturnOk(w)
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

	return basic.ReturnOk(w)
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

	return basic.ReturnOk(w)
}

func (root *Root) ViewCustomerStoreHandler(
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
	var customerStores []CustomerStore

	if search == "" {
		count, customerStores, err = root.repo.ViewCustomerStore(ctx,
			uint(page), uint(pageSize), sortedBy, sortOrder)
		if err != nil {
			return err
		}
	} else {
		count, customerStores, err = root.repo.SelectCustomerStore(ctx)
		if err != nil {
			return err
		}
		count, customerStores = FuzzySearchCustomerStore(customerStores,
			uint(page), uint(pageSize), search)
	}

	type Response struct {
		StoreList  []CustomerStore `json:"storeList"`
		StoreCount uint            `json:"storeCount"`
	}

	res := Response{
		StoreList:  customerStores,
		StoreCount: count,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) QuerySimpleCustomerHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	query := r.URL.Query().Get("query")

	customers, err := root.repo.SelectSimpleCustomer(ctx)
	if err != nil {
		return err
	}

	customers = FuzzySearchSimpleCustomer(customers, query)

	return json.NewEncoder(w).Encode(customers)
}

func (root *Root) AddCustomerStoreHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	store := InsertStore{}
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		return err
	}

	err = root.repo.InsertCustomerStore(ctx, store)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) UpdateCustomerStoreHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	store := CustomerStore{}
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		return err
	}

	err = root.repo.UpdateCustomerStore(ctx, store)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) DeleteCustomerStoreHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	store := CustomerStore{}
	err := json.NewDecoder(r.Body).Decode(&store)
	if err != nil {
		return err
	}

	err = root.repo.DeleteCustomerStore(ctx, store.Id)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) GetWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	vars := mux.Vars(r)

	warehouseId, err := uuid.Parse(vars["warehouseId"])
	if err != nil {
		return err
	}

	warehouse, err := root.repo.GetWarehouse(ctx, warehouseId)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(warehouse)
}
