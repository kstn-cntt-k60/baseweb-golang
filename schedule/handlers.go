package schedule

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

func (root *Root) ViewStoreCityHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

	var err error
	var city string

	city = query.Get("city")

	list, err := root.repo.ViewStoreIdByCity(ctx, city)
	if err != nil {
		return err
	}

	type Response struct {
		StoreIdList []uuid.UUID `json:"storeIdList"`
	}

	res := Response{
		StoreIdList: list,
	}

	return json.NewEncoder(w).Encode(res)
}
