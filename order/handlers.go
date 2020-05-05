package order

import (
	"baseweb/basic"
	"baseweb/security"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Root struct {
	repo *Repo
}

func InitRoot(repo *Repo) *Root {
	return &Root{
		repo: repo,
	}
}

func (root *Root) ViewCustomerStoreByCustomerHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

	customerId, err := uuid.Parse(query.Get("customerId"))
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
	if sortedByQuery == "name" {
		sortedBy = "name"
	} else if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	search := query.Get("searchText")

	var count int
	var stores []CustomerStore

	if search == "" {
		count, stores, err = root.repo.ViewCustomerStoreByCustomer(ctx,
			customerId, page, pageSize, sortedBy, sortOrder)
		if err != nil {
			return err
		}
	} else {
		stores, err = root.repo.SelectCustomerStoreByCustomer(ctx, customerId)
		if err != nil {
			return err
		}
		count, stores = FuzzySearchCustomerStore(stores, page, pageSize, search)
	}

	type Response struct {
		StoreCount int             `json:"storeCount"`
		StoreList  []CustomerStore `json:"storeList"`
	}

	res := Response{
		StoreCount: count,
		StoreList:  stores,
	}

	return json.NewEncoder(w).Encode(res)
}

var bodyError = errors.New("body constraints violation")

func (root *Root) AddOrderHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	type Request struct {
		CustomerId      uuid.UUID       `json:"customerId"`
		WarehouseId     uuid.UUID       `json:"warehouseId"`
		Products        []ClientProduct `json:"products"`
		Address         string          `json:"address"`
		CustomerStoreId *uuid.UUID      `json:"customerStoreId"`
	}

	req := Request{
		Products: make([]ClientProduct, 0),
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	zero := decimal.Zero
	for _, p := range req.Products {
		if p.Quantity.LessThanOrEqual(zero) {
			return bodyError
		}
	}

	err = root.repo.AddOrder(ctx, req.CustomerId,
		req.WarehouseId, req.Products, req.Address,
		req.CustomerStoreId, userLogin.Id)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) ViewProductInfoByWarehouseHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

	var err error

	var page int
	page, err = strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	var pageSize int
	pageSize, err = strconv.Atoi(query.Get("pageSize"))
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

	warehouseId, err := uuid.Parse(query.Get("warehouseId"))
	if err != nil {
		return err
	}

	search := query.Get("searchText")

	var count int
	var products []ProductInfo

	if search == "" {
		count, products, err = root.repo.ViewProductInfoByWarehouse(ctx,
			warehouseId, page, pageSize, sortedBy, sortOrder)
		if err != nil {
			return err
		}
	} else {
		products, err = root.repo.SelectProductInfoByWarehouse(ctx, warehouseId)
		if err != nil {
			return err
		}
		count, products = FuzzySearchProductInfo(products, page, pageSize, search)
	}

	type Response struct {
		ProductInfoList  []ProductInfo `json:"productInfoList"`
		ProductInfoCount int           `json:"productInfoCount"`
	}

	res := Response{
		ProductInfoList:  products,
		ProductInfoCount: count,
	}

	return json.NewEncoder(w).Encode(res)

}

func (root *Root) ViewSaleOrderHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

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

	statusId, err := strconv.Atoi(query.Get("statusId"))
	if err != nil || statusId < 0 || statusId > 5 {
		statusId = 0
	}

	var count int
	var orders []SaleOrder

	count, orders, err = root.repo.ViewSaleOrder(ctx,
		page, pageSize, sortedBy, sortOrder, statusId)
	if err != nil {
		return err
	}

	type Response struct {
		OrderCount int         `json:"orderCount"`
		OrderList  []SaleOrder `json:"orderList"`
	}

	res := Response{
		OrderCount: count,
		OrderList:  orders,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) ViewSingleSaleOrderHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

	saleOrderId, err := strconv.ParseInt(query.Get("saleOrderId"), 10, 64)
	if err != nil {
		return err
	}

	order, items, err := root.repo.GetSaleOrder(ctx, saleOrderId)
	if err != nil {
		return err
	}

	type Response struct {
		Order      SaleOrder       `json:"order"`
		OrderItems []SaleOrderItem `json:"orderItems"`
	}

	res := Response{
		Order:      order,
		OrderItems: items,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) AcceptSalesOrderHandler(
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

	err = root.repo.AcceptSalesOrder(ctx, req.Id)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) CancelSalesOrderHandler(
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

	err = root.repo.CancelSalesOrder(ctx, req.Id)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}
