package salesroute

import (
	"baseweb/security"
	"encoding/json"
	"errors"
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

type OkResponse struct {
	Status string `json:"status"`
}

var okResponse = OkResponse{
	Status: "ok",
}

func (root *Root) AddSalesmanHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	userLogin := ctx.Value("userLogin").(security.UserLogin)

	salesman := Salesman{}

	err := json.NewDecoder(r.Body).Decode(&salesman)
	if err != nil {
		return err
	}

	salesman.CreatedBy = userLogin.Id
	err = root.repo.InsertSalesman(ctx, salesman)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) ViewSalesmanHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

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
	} else if sortedByQuery == "username" {
		sortedBy = "username"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	count, list, err := root.repo.ViewSalesman(ctx, sortedBy, sortOrder, page, pageSize)
	if err != nil {
		return err
	}

	type Response struct {
		SalesmanList  []ClientSalesman `json:"salesmanList"`
		SalesmanCount int              `json:"salesmanCount"`
	}
	res := Response{
		SalesmanList:  list,
		SalesmanCount: count,
	}
	return json.NewEncoder(w).Encode(res)

}

func (root *Root) ViewPlanningPeriodHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

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
	} else if sortedByQuery == "username" {
		sortedBy = "username"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	count, list, err := root.repo.ViewPlanningPeriod(ctx, sortedBy, sortOrder, page, pageSize)
	if err != nil {
		return err
	}

	type Response struct {
		PlanningList  []ClientPlanningPeriod `json:"planningList"`
		PlanningCount int                    `json:"planningCount"`
	}
	res := Response{
		PlanningList:  list,
		PlanningCount: count,
	}
	return json.NewEncoder(w).Encode(res)
}

func (root *Root) GetPlanningPeriodHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	vars := mux.Vars(r)

	planningId, err := strconv.Atoi(vars["planningId"])
	if err != nil {
		return err
	}

	warehouse, err := root.repo.GetPlanningPeriod(ctx, planningId)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(warehouse)
}

func (root *Root) AddPlanningPeriodHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	planningPeriod := PlanningPeriod{}
	err := json.NewDecoder(r.Body).Decode(&planningPeriod)
	if err != nil {
		return err
	}

	planningPeriod.CreatedBy = userLogin.Id
	err = root.repo.InsertPlanningPeriod(ctx, planningPeriod)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) UpdatePlanningPeriodHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	planningPeriod := PlanningPeriod{}
	err := json.NewDecoder(r.Body).Decode(&planningPeriod)
	if err != nil {
		return err
	}

	err = root.repo.UpdatePlanningPeriod(ctx, planningPeriod)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) DeletePlanningPeriodHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	planningPeriod := PlanningPeriod{}
	err := json.NewDecoder(r.Body).Decode(&planningPeriod)
	if err != nil {
		return err
	}

	err = root.repo.DeletePlanningPeriod(ctx, planningPeriod.Id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) ViewConfigHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

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
	} else if sortedByQuery == "username" {
		sortedBy = "username"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	count, list, err := root.repo.ViewConfig(ctx, sortedBy, sortOrder, page, pageSize)
	if err != nil {
		return err
	}

	type Response struct {
		ClientSalesRouteConfig []ClientSalesRouteConfig `json:"configList"`
		ConfigCount            int                      `json:"configCount"`
	}
	res := Response{
		ClientSalesRouteConfig: list,
		ConfigCount:            count,
	}
	return json.NewEncoder(w).Encode(res)

}

func (root *Root) AddConfigHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	type Request struct {
		RepeatWeek int   `json:"repeatWeek"`
		DayList    []int `json:"dayList"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	if len(req.DayList) == 0 {
		return errors.New("dayList empty")
	}

	err = root.repo.InsertConfig(ctx, req.RepeatWeek, req.DayList, userLogin.Id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) UpdateConfigHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id         int   `json:"id" `
		RepeatWeek int   `json:"repeatWeek"`
		ToBeInsert []int `json:"toBeInsert" `
		ToBeDelete []int `json:"toBeDelete"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.UpdateConfig(ctx, req.Id, req.RepeatWeek, req.ToBeInsert, req.ToBeDelete)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) DeleteConfigHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type DeleteConfig struct {
		Id int `json:"id"`
	}

	config := DeleteConfig{}
	err := json.NewDecoder(r.Body).Decode(&config)
	if err != nil {
		return err
	}

	err = root.repo.DeleteConfig(ctx, config.Id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) AddScheduleHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		PlanningId int       `json:"planningId"`
		CustomerId uuid.UUID `json:"customerId"`
		SalesmanId uuid.UUID `json:"salesmanId"`
		ConfigId   int       `json:"configId"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.InsertSchedule(ctx, req.PlanningId, req.CustomerId, req.SalesmanId, req.ConfigId)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) ViewScheduleHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

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
	} else if sortedByQuery == "username" {
		sortedBy = "username"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	count, list, err := root.repo.ViewSchedule(ctx, sortedBy, sortOrder, page, pageSize)
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

func (root *Root) DeleteScheduleHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id int `json:"id"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.DeleteSchedule(ctx, req.Id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}

func (root *Root) GetScheduleHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	vars := mux.Vars(r)

	scheduleId, err := strconv.Atoi(vars["scheduleId"])
	if err != nil {
		return err
	}

	warehouse, err := root.repo.GetSchedule(ctx, scheduleId)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(warehouse)
}

func (root *Root) ViewUserLoginHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	query := r.URL.Query()

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
	} else if sortedByQuery == "username" {
		sortedBy = "username"
	}

	sortOrder := "desc"
	sortOrderQuery := query.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	count, list, err := root.repo.ViewUserLogin(ctx, sortedBy, sortOrder, page, pageSize)
	if err != nil {
		return err
	}

	type Response struct {
		UserLoginList  []ClientUserLogin `json:"userLoginList"`
		UserLoginCount int               `json:"userLoginCount"`
	}
	res := Response{
		UserLoginList:  list,
		UserLoginCount: count,
	}
	return json.NewEncoder(w).Encode(res)

}

func (root *Root) DeleteSalesmanHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id string `json:"id"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.DeleteSalesman(ctx, req.Id)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(okResponse)
}
