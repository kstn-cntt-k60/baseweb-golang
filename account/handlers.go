package account

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"baseweb/basic"
	"baseweb/security"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Root struct {
	repo *Repo
}

func InitRoot(repo *Repo) *Root {
	return &Root{
		repo: repo,
	}
}

type AddPartyRequest struct {
	PartyTypeId  int16  `json:"partyTypeId"`
	Description  string `json:"description"`
	FirstName    string `json:"firstName"`
	MiddleName   string `json:"middleName"`
	LastName     string `json:"lastName"`
	BirthDate    string `json:"birthDate"`
	GenderId     int16  `json:"genderId"`
	CustomerName string `json:"customerName"`
}

func (root *Root) addPersonHandler(
	w http.ResponseWriter,
	ctx context.Context, tx *sqlx.Tx,
	request AddPartyRequest, id uuid.UUID,
) error {

	p := Person{
		Id:         id,
		FirstName:  request.FirstName,
		MiddleName: request.MiddleName,
		LastName:   request.LastName,
		GenderId:   request.GenderId,
		BirthDate:  request.BirthDate,
	}

	err := root.repo.InsertPerson(ctx, tx, p)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) addCustomerHandler(
	w http.ResponseWriter,
	ctx context.Context, tx *sqlx.Tx,
	request AddPartyRequest, id uuid.UUID,
) error {

	err := root.repo.InsertCustomer(ctx, tx, id, request.CustomerName)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return basic.ReturnOk(w)
}

func (root *Root) AddPartyHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	request := AddPartyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return err
	}

	tx, err := root.repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	id, err := root.repo.InsertParty(ctx, tx,
		request.PartyTypeId, request.Description, userLogin.Id)
	if err != nil {
		return err
	}

	if request.PartyTypeId == 1 {
		return root.addPersonHandler(w, ctx, tx, request, id)
	}
	if request.PartyTypeId == 2 {
		return root.addCustomerHandler(w, ctx, tx, request, id)
	}

	return errors.New("PartyTypeId not supported")
}

func (root *Root) ViewPersonHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	var err error

	query := r.URL.Query()

	page, err := strconv.Atoi(query.Get("page"))
	if err != nil {
		page = 0
	}

	pageSize, err := strconv.Atoi(query.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedByQuery := query.Get("sortedBy")
	sortedBy := "created_at"
	if sortedByQuery == "firstName" {
		sortedBy = "first_name"
	} else if sortedByQuery == "createdAt" {
		sortedBy = "created_at"
	} else if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	} else if sortedByQuery == "birthDate" {
		sortedBy = "birth_date"
	}

	sortOrderQuery := query.Get("sortOrder")
	sortOrder := "desc"
	if sortOrderQuery == "desc" {
		sortOrder = "desc"
	} else if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	searchText := query.Get("searchText")

	var count uint
	var personList []ClientPerson
	if searchText == "" {
		count, personList, err = root.repo.ViewPerson(
			ctx, uint(page), uint(pageSize), sortedBy, sortOrder)
	} else {
		count, personList, err = root.repo.ViewPersonWithFullName(
			ctx, uint(page), uint(pageSize), sortedBy, sortOrder, searchText)
	}
	if err != nil {
		return err
	}

	type Response struct {
		Count      uint           `json:"count"`
		PersonList []ClientPerson `json:"personList"`
	}

	response := Response{
		Count:      count,
		PersonList: personList,
	}
	return json.NewEncoder(w).Encode(response)
}

func (root *Root) ViewCustomerHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	var err error

	query := r.URL.Query()

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

	sortedByQuery := query.Get("sortedBy")
	sortedBy := "created_at"
	if sortedByQuery == "name" {
		sortedBy = "name"
	} else if sortedByQuery == "createdAt" {
		sortedBy = "created_at"
	} else if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	}

	sortOrderQuery := query.Get("sortOrder")
	sortOrder := "desc"
	if sortOrderQuery == "desc" {
		sortOrder = "desc"
	} else if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	searchText := query.Get("searchText")

	var count int
	var customerList []ClientCustomer
	if searchText == "" {
		count, customerList, err = root.repo.ViewCustomer(
			ctx, page, pageSize, sortedBy, sortOrder)
		if err != nil {
			return err
		}
	} else {
		count, customerList, err = root.repo.ViewCustomerWithName(
			ctx, page, pageSize, sortedBy, sortOrder, searchText)
		if err != nil {
			return err
		}
	}

	type Response struct {
		Count        int              `json:"count"`
		CustomerList []ClientCustomer `json:"customerList"`
	}

	response := Response{
		Count:        count,
		CustomerList: customerList,
	}
	return json.NewEncoder(w).Encode(response)
}

func (root *Root) UpdatePersonHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id          uuid.UUID `json:"id"`
		FirstName   string    `json:"firstName"`
		MiddleName  string    `json:"middleName"`
		LastName    string    `json:"lastName"`
		BirthDate   string    `json:"birthDate"`
		GenderId    int16     `json:"genderId"`
		Description string    `json:"description"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	person := ClientPerson{
		Id:          req.Id,
		FirstName:   req.FirstName,
		MiddleName:  req.MiddleName,
		LastName:    req.LastName,
		BirthDate:   req.BirthDate,
		GenderId:    req.GenderId,
		Description: req.Description,
	}
	err = root.repo.UpdatePerson(ctx, person)
	if err != nil {
		return err
	}

	type Response struct {
		Status string `json:"status"`
	}

	res := Response{
		Status: "ok",
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) DeletePersonHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id uuid.UUID `json:"id"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.DeletePerson(ctx, req.Id)
	if err != nil {
		return err
	}

	type Response struct {
		Status string `json:"status"`
	}

	res := Response{
		Status: "ok",
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) UpdateCustomerHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	customer := ClientCustomer{}
	err := json.NewDecoder(r.Body).Decode(&customer)
	if err != nil {
		return err
	}

	err = root.repo.UpdateCustomer(ctx, customer)
	if err != nil {
		return err
	}

	type Response struct {
		Status string `json:"status"`
	}

	res := Response{
		Status: "ok",
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) DeleteCustomerHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id uuid.UUID `json:"id"`
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.DeleteCustomer(ctx, req.Id)
	if err != nil {
		return err
	}

	type Response struct {
		Status string `json:"status"`
	}

	res := Response{
		Status: "ok",
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) QuerySimplePersonHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	fullName := r.URL.Query().Get("query")
	personList, err := root.repo.SelectSimplePersonWithFullName(ctx, fullName)
	if err != nil {
		return err
	}

	return json.NewEncoder(w).Encode(personList)
}

func (root *Root) AddUserLogin(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	userLogin := UserLogin{}
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		return err
	}

	err = root.repo.InsertUserLogin(ctx, userLogin)
	if err != nil {
		return err
	}

	type Response struct {
		Status string `json:"status"`
	}

	res := Response{
		Status: "ok",
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) ViewUserLoginHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()
	queries := r.URL.Query()

	var page int
	page, err := strconv.Atoi(queries.Get("page"))
	if err != nil {
		page = 0
	}

	var pageSize int
	pageSize, err = strconv.Atoi(queries.Get("pageSize"))
	if err != nil {
		pageSize = 10
	}

	sortedBy := "created_at"
	sortedByQuery := queries.Get("sortedBy")
	if sortedByQuery == "username" {
		sortedBy = "username"
	} else if sortedByQuery == "updatedAt" {
		sortedBy = "updated_at"
	}

	sortOrder := "desc"
	sortOrderQuery := queries.Get("sortOrder")
	if sortOrderQuery == "asc" {
		sortOrder = "asc"
	}

	var count uint
	var userLogins []ClientUserLogin

	query := queries.Get("query")
	if query == "" {
		count, userLogins, err = root.repo.ViewUserLogin(
			ctx, uint(page), uint(pageSize), sortedBy, sortOrder)
	} else {
		count, userLogins, err = root.repo.SelectUserLogin(ctx)
		count, userLogins = FuzzySearchUserLogin(userLogins,
			query, uint(page), uint(pageSize))
	}

	if err != nil {
		return err
	}

	type Response struct {
		Count         uint              `json:"count"`
		UserLoginList []ClientUserLogin `json:"userLoginList"`
	}
	response := Response{
		Count:         count,
		UserLoginList: userLogins,
	}
	return json.NewEncoder(w).Encode(response)
}

func (root *Root) UpdateUserLoginHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	userLogin := UserLogin{}
	err := json.NewDecoder(r.Body).Decode(&userLogin)
	if err != nil {
		return err
	}

	err = root.repo.UpdateUserLogin(ctx, userLogin)
	if err != nil {
		return err
	}

	type Response struct {
		Status string `json:"status"`
	}

	res := Response{
		Status: "ok",
	}

	return json.NewEncoder(w).Encode(res)
}

func (root *Root) DeleteUserLoginHandler(
	w http.ResponseWriter, r *http.Request) error {

	ctx := r.Context()

	type Request struct {
		Id uuid.UUID `json:"id"`
	}
	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return err
	}

	err = root.repo.DeleteUserLogin(ctx, req.Id)
	if err != nil {
		return err
	}

	type Response struct {
		Status string `json:"status"`
	}

	res := Response{
		Status: "ok",
	}

	return json.NewEncoder(w).Encode(res)
}
