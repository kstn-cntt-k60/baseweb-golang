package salesroute

import (
	"context"
	"fmt"
	"log"

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

func (repo *Repo) InsertSalesman(ctx context.Context, salesman Salesman) error {
	log.Println("InsertSalesman", salesman.Id, salesman.CreatedBy)

	query := `insert into salesman(id, created_by_user_login_id)
	values (:id, :created_by_user_login_id)`

	_, err := repo.db.NamedExecContext(ctx, query, salesman)
	return err
}

func (repo *Repo) ViewSalesman(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int,
) (int, []ClientSalesman, error) {
	log.Println("ViewSalesman", sortedBy, sortOrder, page, pageSize)

	count := 0
	result := make([]ClientSalesman, 0)

	query := `select count(*) from salesman`
	err := repo.db.GetContext(ctx, &count, query)
	if err != nil {
		return count, result, err
	}

	query = `select s.id, u1.username,
	u2.username as created_by,
	s.created_at, s.updated_at
	from salesman s
	inner join user_login u1 on u1.id = s.id
	inner join user_login u2 on u2.id = s.created_by_user_login_id
	order by %s %s
	limit ? offset ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, pageSize, page*pageSize)
	return count, result, err
}

func (repo *Repo) ViewPlanningPeriod(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int,
) (int, []ClientPlanningPeriod, error) {
	log.Println("View Planning Period", sortedBy, sortOrder, page, pageSize)

	count := 0
	result := make([]ClientPlanningPeriod, 0)

	query := `select count(*) from sales_route_planning_period`
	err := repo.db.GetContext(ctx, &count, query)
	if err != nil {
		return count, result, err
	}

	query = `select s.id, s.from_date, s.thru_date, u.username as created_by, 
	s.created_at, s.updated_at
	from sales_route_planning_period s
	inner join user_login u on u.id = s.created_by_user_login_id
	order by %s %s
	limit ? offset ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, pageSize, page*pageSize)
	return count, result, err
}

func (repo *Repo) GetPlanningPeriod(
	ctx context.Context, id int) (ClientPlanningPeriod, error) {

	log.Println("Get Planning Period", id)

	query := `select s.id, s.from_date, s.thru_date, u.username as created_by, 
        s.created_at, s.updated_at
        from sales_route_planning_period s
        inner join user_login u on u.id = s.created_by_user_login_id
		where s.id = ?`

	query = repo.db.Rebind(query)

	planningPeriod := ClientPlanningPeriod{}
	err := repo.db.GetContext(ctx, &planningPeriod, query, id)

	return planningPeriod, err
}

func (repo *Repo) InsertPlanningPeriod(
	ctx context.Context, planningPeriod PlanningPeriod) error {

	log.Println("InsertPlanningPeriod", planningPeriod.FromDate, planningPeriod.ThruDate)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `insert into sales_route_planning_period(
        from_date, thru_date, created_by_user_login_id)
        values (:from_date, :thru_date, :created_by_user_login_id)`
	_, err = tx.NamedExecContext(ctx, query, planningPeriod)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *Repo) UpdatePlanningPeriod(
	ctx context.Context, planningPeriod PlanningPeriod) error {

	log.Println("UpdatePlanningPeriod", planningPeriod.Id,
		planningPeriod.FromDate, planningPeriod.ThruDate, planningPeriod.CreatedBy, planningPeriod.UpdatedAt)

	query := `update sales_route_planning_period set from_date = :from_date,
	thru_date = :thru_date where id = :id`
	_, err := repo.db.NamedExecContext(ctx, query, planningPeriod)
	return err
}

func (repo *Repo) DeletePlanningPeriod(
	ctx context.Context, id int) error {

	log.Println("DeletePlanningPeriod", id)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := "delete from sales_route_planning_period where id = ?"
	query = repo.db.Rebind(query)
	_, err = repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}
