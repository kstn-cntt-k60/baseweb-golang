package account

import (
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type Repo struct {
	db              *sqlx.DB
	countPerson     *sqlx.Stmt
	viewPerson      *sqlx.Stmt
	countCustomer   *sqlx.Stmt
	viewCustomer    *sqlx.Stmt
	countUserLogin  *sqlx.Stmt
	viewUserLogin   *sqlx.Stmt
	selectUserLogin *sqlx.Stmt
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

	query = "select count(*) from user_login"
	countUserLogin, err := db.Preparex(query)
	if err != nil {
		log.Panicln(err)
	}

	query = `select u.id, u.username,
        u.created_at, u.updated_at,
        p.first_name, p.middle_name, p.last_name,
        p.birth_date, p.gender_id
        from user_login u
        inner join person p on u.person_id = p.id
        order by u.created_at desc
        limit ? offset ?`
	viewUserLogin, err := db.Preparex(db.Rebind(query))
	if err != nil {
		log.Panicln(err)
	}

	query = `select u.id, u.username,
        u.created_at, u.updated_at,
        p.first_name, p.middle_name, p.last_name,
        p.birth_date, p.gender_id
        from user_login u
        inner join person p on u.person_id = p.id`
	selectUserLogin, err := db.Preparex(query)
	if err != nil {
		log.Panicln(err)
	}

	return &Repo{
		db:              db,
		countPerson:     countPerson,
		viewPerson:      viewPerson,
		countCustomer:   countCustomer,
		viewCustomer:    viewCustomer,
		countUserLogin:  countUserLogin,
		viewUserLogin:   viewUserLogin,
		selectUserLogin: selectUserLogin,
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

func (repo *Repo) SelectSimplePerson(
	ctx context.Context) ([]SimplePerson, error) {

	log.Println("SelectSimplePerson")

	query := `select id, first_name, middle_name, last_name,
            birth_date, gender_id from person`

	result := make([]SimplePerson, 0)
	return result, repo.db.SelectContext(ctx, &result, query)
}

func (repo *Repo) InsertUserLogin(
	ctx context.Context, userLogin UserLogin) error {

	log.Println("InsertUserLogin", userLogin.Username,
		userLogin.Password, userLogin.PersonId)

	hash, err := bcrypt.GenerateFromPassword([]byte(userLogin.Password), 10)
	if err != nil {
		return err
	}
	userLogin.Password = string(hash)

	query := `insert into user_login(username, password, person_id)
            values (:username, :password, :person_id)`

	_, err = repo.db.NamedExecContext(ctx, query, userLogin)
	return err
}

func (repo *Repo) ViewUserLogin(
	ctx context.Context, page, pageSize uint,
	sortedBy, sortOrder string) (uint, []ClientUserLogin, error) {

	log.Println("ViewUserLogin", page, pageSize)

	var count uint
	result := make([]ClientUserLogin, 0)

	err := repo.countUserLogin.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	if sortedBy == "created_at" && sortOrder == "desc" {
		err = repo.viewUserLogin.SelectContext(ctx,
			&result, pageSize, page*pageSize)
		return count, result, err
	} else {
		query := fmt.Sprintf(`select u.id, u.username,
            u.created_at, u.updated_at,
            p.first_name, p.middle_name, p.last_name,
            p.birth_date, p.gender_id
            from user_login u
            inner join person p on u.person_id = p.id
            order by u.%s %s
            limit ? offset ?`, sortedBy, sortOrder)
		log.Println("[SQL]", query)

		err = repo.db.SelectContext(ctx, &result,
			repo.db.Rebind(query), pageSize, page*pageSize)
		return count, result, err
	}
}

func (repo *Repo) SelectUserLogin(
	ctx context.Context) (uint, []ClientUserLogin, error) {

	log.Println("SelectUserLogin")

	var count uint
	result := make([]ClientUserLogin, 0)

	err := repo.countUserLogin.GetContext(ctx, &count)
	if err != nil {
		return count, result, err
	}

	return count, result, repo.selectUserLogin.SelectContext(ctx, &result)
}

func (repo *Repo) UpdateUserLogin(
	ctx context.Context, userLogin UserLogin) error {

	log.Println("UpdateUserLogin", userLogin.Id,
		userLogin.Username, userLogin.Password)

	if userLogin.Password == "" {
		query := `update user_login
            set username = :username where id = :id`
		_, err := repo.db.NamedExecContext(ctx, query, userLogin)
		return err
	} else {
		query := `update user_login
            set username = :username, password = :password
            where id = :id`

		hash, err := bcrypt.GenerateFromPassword(
			[]byte(userLogin.Password), 10)
		if err != nil {
			return err
		}

		userLogin.Password = string(hash)
		_, err = repo.db.NamedExecContext(ctx, query, userLogin)
		return err
	}
}

func (repo *Repo) DeleteUserLogin(
	ctx context.Context, id uuid.UUID) error {

	log.Println("DeleteUserLogin", id)

	query := "delete from user_login where id = ?"
	_, err := repo.db.ExecContext(ctx, repo.db.Rebind(query), id)

	return err
}
