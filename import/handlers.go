package importProduct

import (
	"encoding/json"
	"net/http"

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

	count, products, err := root.repo.ViewProductByWarehouse(
		ctx, warehouseId)
	if err != nil {
		return err
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
