package export

import (
	"baseweb/basic"
	"baseweb/order"
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

func (root *Root) ExportSaleOrderItemHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		SaleOrderId   int64     `json:"saleOrderId" db:"sale_order_id"`
		SaleOrderSeq  int       `json:"saleOrderSeq" db:"sale_order_seq"`
		EffectiveFrom time.Time `json:"effectiveFrom" db:"effective_from"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.ExportSaleOrderItem(ctx,
		req.SaleOrderId, req.SaleOrderSeq, req.EffectiveFrom)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) ViewExportableSalesOrderHandler(
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

	var count int
	var orders []order.SaleOrder

	count, orders, err = root.repo.ViewExportableSalesOrder(ctx,
		page, pageSize, sortedBy, sortOrder)
	if err != nil {
		return err
	}

	type Response struct {
		OrderCount int               `json:"orderCount"`
		OrderList  []order.SaleOrder `json:"orderList"`
	}

	res := Response{
		OrderCount: count,
		OrderList:  orders,
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) CompleteSalesOrderHandler(
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

	err = root.repo.CompleteSalesOrder(ctx, req.Id)
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) ViewCompletedSalesOrderHandler(
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

	var count int
	var orders []order.SaleOrder

	count, orders, err = root.repo.ViewCompletedSalesOrder(ctx,
		page, pageSize, sortedBy, sortOrder)
	if err != nil {
		return err
	}

	type Response struct {
		OrderCount int               `json:"orderCount"`
		OrderList  []order.SaleOrder `json:"orderList"`
	}

	res := Response{
		OrderCount: count,
		OrderList:  orders,
	}

	return json.NewEncoder(w).Encode(res)
}
