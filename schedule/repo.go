package schedule

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repo struct {
	db *sqlx.DB
}

func InitRepo(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (repo *Repo) ViewStoreIdByCity(
	ctx context.Context, city string,
) ([]uuid.UUID, error) {
	log.Println("ViewStoreCity", city)

	result := make([]uuid.UUID, 0)

	if city == "hanoi" {
		city = "%Hà Nội%"
	} else if city == "hcm" {
		city = "%Hồ Chí Minh%"
	} else {
		return result, errors.New("unsupported city")
	}

	query := repo.db.Rebind(
		`select f.id from facility f 
		where f.address like ?`)

	err := repo.db.SelectContext(ctx, &result, query, city)
	return result, err
}
