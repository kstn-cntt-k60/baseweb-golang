package account

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"baseweb/security"

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

func (root *Root) addPersonHandler(w http.ResponseWriter,
	ctx context.Context,
	tx *sql.Tx,
	request AddPartyRequest, id uuid.UUID) {

	err := root.repo.InsertPerson(ctx, tx, id,
		request.FirstName, request.MiddleName,
		request.LastName, request.GenderId, request.BirthDate)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	type Response struct {
		Party  Party  `json:"party"`
		Person Person `json:"person"`
	}

	party, err := root.repo.GetParty(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	person, err := root.repo.GetPerson(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{
		Party:  party,
		Person: person,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln(err)
	}
}

func (root *Root) addCustomerHandler(w http.ResponseWriter,
	ctx context.Context,
	tx *sql.Tx,
	request AddPartyRequest, id uuid.UUID) {

	err := root.repo.InsertCustomer(ctx, tx, id, request.CustomerName)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	err = tx.Commit()
	if err == context.Canceled || err == context.DeadlineExceeded {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		log.Panicln(err)
	}

	type Response struct {
		Party    Party    `json:"party"`
		Customer Customer `json:"customer"`
	}

	party, err := root.repo.GetParty(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	customer, err := root.repo.GetCustomer(ctx, id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response := Response{
		Party:    party,
		Customer: customer,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Panicln(err)
	}
}

func (root *Root) AddPartyHandler(
	w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userLogin := ctx.Value("userLogin").(security.UserLogin)

	request := AddPartyRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("[ERROR]", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	tx, err := root.repo.db.BeginTx(ctx, nil)
	if err == context.Canceled || err == context.DeadlineExceeded {
		tx.Rollback()
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if err != nil {
		tx.Rollback()
		log.Panicln(err)
	}

	id, err := root.repo.InsertParty(ctx, tx,
		request.PartyTypeId, request.Description, userLogin.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		tx.Rollback()
		return
	}

	if request.PartyTypeId == 1 {
		root.addPersonHandler(w, ctx, tx, request, id)
		return
	}
	if request.PartyTypeId == 2 {
		root.addCustomerHandler(w, ctx, tx, request, id)
		return
	}

	log.Println("[ERROR]", "PartyTypeId not supported")
	w.WriteHeader(http.StatusInternalServerError)
}
