package account

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/google/uuid"
)

type Repo struct {
	db            *sql.DB
	countPerson   *sql.Stmt
	viewPerson    *sql.Stmt
	countCustomer *sql.Stmt
	viewCustomer  *sql.Stmt
}

func InitRepo(db *sql.DB) *Repo {
	countPerson, err := db.Prepare(`select count(*) from person`)
	if err != nil {
		log.Panicln(err)
	}

	viewPerson, err := db.Prepare(
		`select person.id, 
        person.first_name, person.middle_name, person.last_name,
        person.gender_id, person.birth_date,
        person.created_at, person.updated_at,
        party.description
        from person
        inner join party on party.id = person.id
        order by person.created_at desc
        limit $1 offset $2`)
	if err != nil {
		log.Panicln(err)
	}

	countCustomer, err := db.Prepare(
		`select count(*) from customer`)
	if err != nil {
		log.Panicln(err)
	}

	viewCustomer, err := db.Prepare(
		`select c.id, c.name,
        c.created_at, c.updated_at,
        p.description
        from customer c
        inner join party p on p.id = c.id
        order by c.created_at desc
        limit $1 offset $2`)
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
	tx *sql.Tx, partyTypeId int16, description string,
	userLoginId uuid.UUID) (uuid.UUID, error) {

	log.Println("InsertParty", partyTypeId, description)

	id, err := uuid.NewUUID()
	if err != nil {
		log.Panicln(err)
	}

	_, err = tx.ExecContext(ctx,
		`insert into party(
            id, party_type_id, description,
            created_by_user_login_id,
            updated_by_user_login_id)
        values ($1, $2, $3, $4, $4)`,
		id, partyTypeId, description, userLoginId)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return id, err
	}
	if err != nil {
		log.Panicln(err)
	}

	return id, nil
}

func (repo *Repo) InsertPerson(ctx context.Context,
	tx *sql.Tx, id uuid.UUID,
	firstName, middleName, lastName string,
	genderId int16,
	birthDate string) error {

	log.Println("InsertPerson", firstName, middleName, lastName)

	_, err := tx.ExecContext(ctx,
		`insert into person(
        id, first_name, middle_name, last_name,
        gender_id, birth_date)
        values ($1, $2, $3, $4, $5, $6)`,
		id, firstName, middleName, lastName,
		genderId, birthDate)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}

	return nil
}

func (repo *Repo) InsertCustomer(ctx context.Context,
	tx *sql.Tx, id uuid.UUID, customerName string) error {

	log.Println("InsertCustomer", id, customerName)

	_, err := tx.ExecContext(ctx,
		`insert into customer(id, name) values ($1, $2)`,
		id, customerName)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}

	return nil
}

func (repo *Repo) GetParty(ctx context.Context,
	id uuid.UUID) (Party, error) {

	log.Println("GetParty", id)

	row := repo.db.QueryRowContext(ctx,
		`select id, party_type_id, description,
            created_at, updated_at,
            created_by_user_login_id,
            updated_by_user_login_id
        from party where id = $1`, id)

	party := Party{}
	err := row.Scan(&party.Id, &party.TypeId, &party.Description,
		&party.CreatedAt, &party.UpdatedAt,
		&party.CreatedBy, &party.UpdatedBy)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return party, err
	}
	if err != nil {
		log.Panicln(err)
	}

	return party, nil
}

func (repo *Repo) GetPerson(ctx context.Context,
	id uuid.UUID) (Person, error) {

	log.Println("GetPerson", id)

	row := repo.db.QueryRowContext(ctx,
		`select id, first_name, middle_name, last_name,
            gender_id, birth_date, created_at, updated_at
        from person where id = $1`, id)

	person := Person{}
	err := row.Scan(&person.Id,
		&person.FirstName, &person.MiddleName, &person.LastName,
		&person.GenderId, &person.BirthDate,
		&person.CreatedAt, &person.UpdatedAt)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return person, err
	}
	if err != nil {
		log.Panicln(err)
	}

	return person, nil
}

func (repo *Repo) GetCustomer(ctx context.Context,
	id uuid.UUID) (Customer, error) {

	log.Println("GetCustomer", id)

	row := repo.db.QueryRowContext(ctx,
		`select id, name, created_at, updated_at
        from customer where id = $1`, id)

	c := Customer{}
	err := row.Scan(&c.Id, &c.Name, &c.CreatedAt, &c.UpdatedAt)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return c, err
	}
	if err != nil {
		log.Panicln(err)
	}

	return c, nil
}

func (repo *Repo) ViewPerson(ctx context.Context,
	page uint, pageSize uint,
	sortedBy, sortOrder string) (uint, []ClientPerson, error) {

	log.Println("ViewPerson", page, pageSize)

	var count uint = 0
	result := make([]ClientPerson, 0)

	row := repo.countPerson.QueryRowContext(ctx)
	err := row.Scan(&count)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return count, result, err
	}
	if err != nil {
		log.Panicln(err)
	}

	var rows *sql.Rows
	if sortedBy == "created_at" && sortOrder == "desc" {
		rows, err = repo.viewPerson.QueryContext(
			ctx, pageSize, page*pageSize)
	} else {
		query := fmt.Sprintf(`select person.id, 
                person.first_name, person.middle_name, person.last_name,
                person.gender_id, person.birth_date,
                person.created_at, person.updated_at,
                party.description
                from person
                inner join party on party.id = person.id
                order by person.%s %s
                limit $1 offset $2`, sortedBy, sortOrder)

		log.Println("[SQL]", query)

		rows, err = repo.db.QueryContext(
			ctx, query, pageSize, page*pageSize)
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		return count, result, err
	}
	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		p := ClientPerson{}
		err = rows.Scan(&p.Id,
			&p.FirstName, &p.MiddleName, &p.LastName,
			&p.GenderId, &p.BirthDate,
			&p.CreatedAt, &p.UpdatedAt,
			&p.Description)
		if err == context.Canceled || err == context.DeadlineExceeded {
			return count, result, err
		}
		if err != nil {
			log.Panicln(err)
		}

		result = append(result, p)
	}

	return count, result, nil
}

func (repo *Repo) ViewCustomer(ctx context.Context,
	page uint, pageSize uint,
	sortedBy, sortOrder string) (uint, []ClientCustomer, error) {

	log.Println("ViewCustomer", page, pageSize)

	var count uint = 0
	result := make([]ClientCustomer, 0)

	row := repo.countCustomer.QueryRowContext(ctx)
	err := row.Scan(&count)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return count, result, err
	}
	if err != nil {
		log.Panicln(err)
	}

	var rows *sql.Rows
	if sortedBy == "created_at" && sortOrder == "desc" {
		rows, err = repo.viewCustomer.QueryContext(
			ctx, pageSize, page*pageSize)
	} else {
		query := fmt.Sprintf(
			`select c.id, c.name,
            c.created_at, c.updated_at,
            p.description
            from customer c
            inner join party p on p.id = c.id
            order by c.%s %s
            limit $1 offset $2`,
			sortedBy, sortOrder)

		log.Println("[SQL]", query)

		rows, err = repo.db.QueryContext(
			ctx, query, pageSize, page*pageSize)
	}
	if err == context.Canceled || err == context.DeadlineExceeded {
		return count, result, err
	}
	if err != nil {
		log.Panicln(err)
	}

	for rows.Next() {
		c := ClientCustomer{}
		err = rows.Scan(&c.Id, &c.Name,
			&c.CreatedAt, &c.UpdatedAt, &c.Description)
		if err == context.Canceled || err == context.DeadlineExceeded {
			return count, result, err
		}
		if err != nil {
			log.Panicln(err)
		}

		result = append(result, c)
	}

	return count, result, nil
}

func (repo *Repo) UpdatePerson(
	ctx context.Context, id uuid.UUID,
	firstName, middleName, lastName string,
	genderId int16, birthDate string,
	description string) error {

	log.Println("UpdatePerson", id, firstName, middleName,
		lastName, genderId, birthDate, description)

	_, err := repo.db.ExecContext(ctx,
		`update person set first_name = $2,
        middle_name = $3, last_name = $4,
        gender_id = $5, birth_date = $6
        where id = $1`,
		id, firstName, middleName,
		lastName, genderId, birthDate)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}

	_, err = repo.db.ExecContext(ctx,
		`update party set description = $2 where id = $1`, id, description)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}

	return nil
}

func (repo *Repo) DeletePerson(ctx context.Context, id uuid.UUID) error {
	log.Println("DeletePerson", id)

	tx, err := repo.db.BeginTx(ctx, nil)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("delete from person where id = $1", id)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}

	_, err = tx.Exec("delete from party where id = $1", id)
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}

	err = tx.Commit()
	if err == context.Canceled || err == context.DeadlineExceeded {
		return err
	}
	if err != nil {
		log.Panicln(err)
	}

	return nil
}
