package salesman

import (
	"context"
	"fmt"
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

func (repo *Repo) ViewSchedule(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int, userLoginId uuid.UUID,
) (int, []ClientSchedule, error) {
	log.Println("ClientSchedule", sortedBy, sortOrder, page, pageSize, userLoginId)

	count := 0
	result := make([]ClientSchedule, 0)

	query := repo.db.Rebind(
		`select count(*) from sales_route_detail where salesman_id = ?`)
	err := repo.db.GetContext(ctx, &count, query, userLoginId)
	if err != nil {
		return count, result, err
	}

	query = `select srd.id, srd.config_id, c.repeat_week, 
	string_agg(d."day"::text, ', ') as day_list, srd.planning_period_id as planning_id,  
	pp.from_date, pp.thru_date, customer."name" as customer_name, u.username as salesman_name,
	srd.created_at, srd.updated_at
	from sales_route_detail srd
	inner join sales_route_planning_period pp on srd.planning_period_id = pp.id
	inner join sales_route_config c on srd.config_id = c.id
	inner join sales_route_config_day d on d.config_id = c.id
	inner join customer on srd.customer_id = customer.id
	inner join salesman on srd.salesman_id = salesman.id
	inner join user_login u on salesman.id = u.id
	where srd.salesman_id = ?
	group by srd.id, c.repeat_week, customer."name", pp.from_date, pp.thru_date, u.username
	order by %s %s
	limit ? offset ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, userLoginId, pageSize, page*pageSize)
	return count, result, err
}

func (repo *Repo) InsertSchedule(
	ctx context.Context, planningId int, customerId uuid.UUID,
	salesmanId uuid.UUID,
) error {
	log.Println("Insert Checkin", planningId, customerId, salesmanId)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := repo.db.Rebind(
		`select id from sales_route_detail 
		where planning_period_id = ? and 
		customer_id = ? and salesman_id = ?
		`)

	log.Println("query", query)

	var detailId int
	err = tx.GetContext(ctx, &detailId, query, planningId, customerId, salesmanId)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(
		`insert into salesman_checkin_history (id) values (?)`)

	_, err = tx.ExecContext(ctx, query, detailId)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func (repo *Repo) InsertCheckin(
	ctx context.Context, detailId int) error {
	log.Println("Insert Checkin", detailId)

	query := repo.db.Rebind(
		`insert into salesman_checkin_history(sales_route_detail_id)
		values (?)
		`)
	_, err := repo.db.ExecContext(ctx, query, detailId)

	return err
}

func (repo *Repo) ViewCheckinHistory(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int, userLoginId uuid.UUID,
) (int, []ClientCheckinHistory, error) {
	log.Println("ViewCheckinHistory", sortedBy, sortOrder, page, pageSize, userLoginId)

	count := 0
	result := make([]ClientCheckinHistory, 0)

	query := repo.db.Rebind(
		`select count(*) 
		from salesman_checkin_history history 
		inner join sales_route_detail detail 
		on detail.id = history.sales_route_detail_id
		where detail.salesman_id = ?`)
	err := repo.db.GetContext(ctx, &count, query, userLoginId)
	if err != nil {
		return count, result, err
	}

	query = `select c.id, c.checkin_time, 
	d.planning_period_id as planning_id, p.from_date, p.thru_date, 
	d.config_id, customer."name" as customer_name
	from salesman_checkin_history c 
	inner join sales_route_detail d on d.id=c.sales_route_detail_id
	inner join sales_route_planning_period p on p.id = d.planning_period_id
	inner join customer on customer.id = d.customer_id
	inner join sales_route_config config on config.id = d.config_id
	where d.salesman_id = ?
	order by %s %s
	limit ? offset ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, userLoginId, pageSize, page*pageSize)
	return count, result, err
}
