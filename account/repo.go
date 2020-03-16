package account

import (
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
)

type Repo struct {
	db *sql.DB
}

func InitRepo(db *sql.DB) *Repo {
	return &Repo{
		db: db,
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
		tx.Rollback()
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
		tx.Rollback()
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
		tx.Rollback()
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
