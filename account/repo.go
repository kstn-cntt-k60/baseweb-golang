package account

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db            *sqlx.DB
	countPerson   *sqlx.Stmt
	viewPerson    *sqlx.Stmt
	countCustomer *sqlx.Stmt
	viewCustomer  *sqlx.Stmt
}

func InitRepo(db *sqlx.DB) *Repo {
	countPerson, err := db.Preparex(`select count(*) from person`)
	if err != nil {
		log.Panicln(err)
	}

	query := `select person.id as id, 
        first_name, middle_name, last_name,
        gender_id, birth_date,
        person.created_at as created_at,
        person.updated_at as updated_at,
        party.description as description
        from person
        inner join party on party.id = person.id
        order by person.created_at desc
        limit ? offset ?`
	viewPerson, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	countCustomer, err := db.Preparex(
		`select count(*) from customer`)
	if err != nil {
		log.Panicln(err)
	}

	query = `select c.id, c.name,
        c.created_at, c.updated_at,
        p.description
        from customer c
        inner join party p on p.id = c.id
        order by c.created_at desc
        limit ? offset ?`
	viewCustomer, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	return &Repo{
		db:            db,
		countPerson:   countPerson,
		viewPerson:    viewPerson,
		countCustomer: countCustomer,
		viewCustomer:  viewCustomer,
	}
}

func (repo *Repo) InsertParty(ctx context.Context,
	tx *sqlx.Tx, partyTypeId int16, description string,
	userLoginId uuid.UUID) (uuid.UUID, error) {

	log.Println("InsertParty", partyTypeId, description)

	id, err := uuid.NewUUID()
	if err != nil {
		return id, err
	}

	query := `insert into party(
                id, party_type_id, description,
                created_by_user_login_id,
                updated_by_user_login_id)
            values (?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, repo.db.Rebind(query),
		id, partyTypeId, description, userLoginId, userLoginId)

	return id, err
}

func (repo *Repo) InsertPerson(ctx context.Context,
	tx *sqlx.Tx, person Person) error {

	log.Println("InsertPerson", person.FirstName,
		person.MiddleName, person.LastName)

	_, err := tx.NamedExecContext(ctx,
		`insert into person(
        id, first_name, middle_name, last_name,
        gender_id, birth_date)
        values (:id, :first_name, :middle_name,
        :last_name, :gender_id, :birth_date)`,
		person)

	return err
}

func (repo *Repo) InsertCustomer(ctx context.Context,
	tx *sqlx.Tx, id uuid.UUID, customerName string) error {

	log.Println("InsertCustomer", id, customerName)

	query := `insert into customer(id, name) values (?, ?)`
	_, err := tx.ExecContext(ctx,
		repo.db.Rebind(query), id, customerName)

	return err
}

func (repo *Repo) GetParty(ctx context.Context,
	id uuid.UUID) (Party, error) {

	log.Println("GetParty", id)

	party := Party{}
	query := "select * from party where id = ?"
	err := repo.db.GetContext(ctx, &party, repo.db.Rebind(query), id)
	return party, err
}

func (repo *Repo) GetPerson(ctx context.Context,
	id uuid.UUID) (Person, error) {

	log.Println("GetPerson", id)

	query := "select * from person where id = ?"
	person := Person{}
	err := repo.db.GetContext(ctx, &person, repo.db.Rebind(query), id)
	return person, err
}

func (repo *Repo) GetCustomer(ctx context.Context,
	id uuid.UUID) (Customer, error) {

	log.Println("GetCustomer", id)

	c := Customer{}
	query := "select * from customer where id = ?"
	err := repo.db.GetContext(ctx, &c, repo.db.Rebind(query), id)
	return c, err
}

func (repo *Repo) ViewPerson(ctx context.Context,
	page uint, pageSize uint,
	sortedBy, sortOrder string) (uint, []ClientPerson, error) {

	log.Println("ViewPerson", page, pageSize, sortedBy, sortOrder)

	var count uint = 0
	result := make([]ClientPerson, 0)

	err := repo.countPerson.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "created_at" && sortOrder == "desc" {
		err = repo.viewPerson.SelectContext(
			ctx, &result, pageSize, page*pageSize)
		return count, result, err
	} else {
		query := fmt.Sprintf(`select person.id,
                first_name, middle_name, last_name,
                gender_id, birth_date,
                person.created_at, person.updated_at,
                party.description
                from person
                inner join party on party.id = person.id
                order by person.%s %s
                limit ? offset ?`, sortedBy, sortOrder)
		log.Println("[SQL]", query)

		err = repo.db.SelectContext(ctx, &result,
			repo.db.Rebind(query), pageSize, page*pageSize)

		return count, result, err
	}
}

func (repo *Repo) ViewCustomer(ctx context.Context,
	page uint, pageSize uint,
	sortedBy, sortOrder string) (uint, []ClientCustomer, error) {

	log.Println("ViewCustomer", page, pageSize, sortedBy, sortOrder)

	var count uint = 0
	result := make([]ClientCustomer, 0)

	err := repo.countCustomer.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "created_at" && sortOrder == "desc" {
		err = repo.viewCustomer.SelectContext(
			ctx, &result, pageSize, page*pageSize)
		return count, result, err
	} else {
		query := fmt.Sprintf(
			`select c.id, c.name,
            c.created_at, c.updated_at,
            p.description
            from customer c
            inner join party p on p.id = c.id
            order by c.%s %s
            limit ? offset ?`,
			sortedBy, sortOrder)

		log.Println("[SQL]", query)

		err = repo.db.SelectContext(ctx, &result,
			repo.db.Rebind(query), pageSize, page*pageSize)
		return count, result, err
	}
}

func (repo *Repo) UpdatePerson(
	ctx context.Context, person ClientPerson) error {

	log.Println("UpdatePerson", person.Id, person.FirstName,
		person.MiddleName, person.LastName, person.GenderId,
		person.BirthDate, person.Description)

	_, err := repo.db.NamedExecContext(ctx,
		`update person set first_name = :first_name,
        middle_name = :middle_name, last_name = :last_name,
        gender_id = :gender_id, birth_date = :birth_date
        where id = :id`, person)
	if err != nil {
		return err
	}

	_, err = repo.db.NamedExecContext(ctx,
		`update party set description = :description where id = :id`, person)

	return err
}

func (repo *Repo) DeletePerson(ctx context.Context, id uuid.UUID) error {
	log.Println("DeletePerson", id)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "delete from person where id = ?"
	_, err = tx.Exec(repo.db.Rebind(query), id)
	if err != nil {
		return err
	}

	query = "delete from party where id = ?"
	_, err = tx.Exec(repo.db.Rebind(query), id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *Repo) UpdateCustomer(
	ctx context.Context, customer ClientCustomer) error {

	log.Println("UpdateCustomer", customer.Id,
		customer.Description, customer.Name)

	query := `update customer set name = :name where id = :id`
	_, err := repo.db.NamedExecContext(ctx, repo.db.Rebind(query), customer)
	if err != nil {
		return err
	}

	query = `update party set description = :description where id = :id`
	_, err = repo.db.NamedExecContext(ctx, repo.db.Rebind(query), customer)

	return err
}

func (repo *Repo) DeleteCustomer(ctx context.Context, id uuid.UUID) error {
	log.Println("DeleteCustomer", id)

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "delete from customer where id = ?"
	_, err = tx.Exec(repo.db.Rebind(query), id)
	if err != nil {
		return err
	}

	query = "delete from party where id = ?"
	_, err = tx.Exec(repo.db.Rebind(query), id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
