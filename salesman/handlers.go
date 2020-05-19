package salesman

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

func (root *Root) ViewScheduleHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	var err error

	var page int
	page, err = strconv.Atoi(query.Get("page"))
	if err != nil {
		return err
	}

	var pageSize int
	pageSize, err = strconv.Atoi(query.Get("pageSize"))
	if err != nil {
		return err
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

	count, list, err := root.repo.ViewSchedule(ctx, sortedBy, sortOrder, page, pageSize, userLogin.Id)
	if err != nil {
		return err
	}

	type Response struct {
		ClientSchedule []ClientSchedule `json:"scheduleList"`
		ScheduleCount  int              `json:"scheduleCount"`
	}
	res := Response{
		ClientSchedule: list,
		ScheduleCount:  count,
	}
	return json.NewEncoder(w).Encode(res)
}

func (root *Root) AddCheckinHandler(
	w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	type Request struct {
		DetailId int `json:"detailId"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode((&req))
	if err != nil {
		return err
	}

	err = root.repo.InsertCheckin(ctx, req.DetailId)
	if err != nil {
		return err
	}
	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) ViewCheckinHistoryHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	var err error

	var page int
	page, err = strconv.Atoi(query.Get("page"))
	if err != nil {
		return err
	}

	var pageSize int
	pageSize, err = strconv.Atoi(query.Get("pageSize"))
	if err != nil {
		return err
	}

	sortedBy := "checkin_time"
	sortedByQuery := query.Get("sortedBy")
	if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	} else if sortedByQuery == "username" {
		sortedBy = "username"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	count, list, err := root.repo.ViewCheckinHistory(ctx, sortedBy, sortOrder, page, pageSize, userLogin.Id)
	if err != nil {
		return err
	}

	type Response struct {
		ClientHistory []ClientCheckinHistory `json:"checkinList"`
		CheckinCount  int                    `json:"checkinCount"`
	}
	res := Response{
		ClientHistory: list,
		CheckinCount:  count,
	}
	return json.NewEncoder(w).Encode(res)
}
