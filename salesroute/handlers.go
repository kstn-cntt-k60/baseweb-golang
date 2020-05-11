package salesroute

import (
	"baseweb/security"
	"encoding/json"
	"net/http"
	"strconv"

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
