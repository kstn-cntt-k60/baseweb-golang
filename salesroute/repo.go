package salesroute

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

func (repo *Repo) ViewUserLogin(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int,
) (int, []ClientUserLogin, error) {
	log.Println("ViewUserLogin", sortedBy, sortOrder, page, pageSize)

	count := 0
	result := make([]ClientUserLogin, 0)

	query := `select count(*) from user_login 
	where id not in (select id from salesman);`
	err := repo.db.GetContext(ctx, &count, query)
	if err != nil {
		return count, result, err
	}

	query = `select u.id, u.username, u.created_at, u.updated_at, 
	p.first_name, p.middle_name, p.last_name, p.birth_date, 
	p.gender_id from user_login u 
	inner join person p on p.id = u.person_id
	where u.id not in (select id from salesman)
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

	query := `insert into sales_route_planning_period(
        from_date, thru_date, created_by_user_login_id)
        values (:from_date, :thru_date, :created_by_user_login_id)`
	_, err := repo.db.NamedExecContext(ctx, query, planningPeriod)
	return err
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

	query := "delete from sales_route_planning_period where id = ?"
	query = repo.db.Rebind(query)
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *Repo) DeleteSalesman(
	ctx context.Context, id string) error {

	log.Println("DeleteSalesman", id)
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `delete from salesman_checkin_history
		where sales_route_detail_id in (
		select id from sales_route_detail srd
		where srd.salesman_id = ?
		)`
	query = repo.db.Rebind(query)
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(`delete from sales_route_detail where salesman_id = ?`)
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(`delete from salesman where id = ?`)
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()

}

func (repo *Repo) ViewConfig(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int,
) (int, []ClientSalesRouteConfig, error) {
	log.Println("ViewSalesRouteConfig", sortedBy, sortOrder, page, pageSize)

	count := 0
	result := make([]ClientSalesRouteConfig, 0)

	query := `select count(*) from sales_route_config`
	err := repo.db.GetContext(ctx, &count, query)
	if err != nil {
		return count, result, err
	}

	query = `select c.id, c.repeat_week,
		string_agg(d."day"::text, ', ') as day_list, 
		u.username as created_by ,
		c.created_at, c.updated_at
		from sales_route_config c
		inner join sales_route_config_day d on d.config_id = c.id
		inner join user_login u on u.id = c.created_by_user_login_id
		group by c.id, c.repeat_week, c.created_by_user_login_id, c.updated_at, u.username
		order by %s %s
		limit ? offset ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, pageSize, page*pageSize)
	return count, result, err
}

func (repo *Repo) InsertConfig(
	ctx context.Context, repeatWeek int,
	dayList []int, userLoginId uuid.UUID,
) error {
	log.Println("InsertConfig", repeatWeek, dayList, userLoginId)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := repo.db.Rebind(
		`insert into sales_route_config(
        repeat_week, created_by_user_login_id)
		values (?, ?) returning id
		`)

	var configId int
	err = tx.GetContext(ctx, &configId, query, repeatWeek, userLoginId)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(
		`insert into sales_route_config_day(config_id, day) values(?, ?)`)

	for _, day := range dayList {
		_, err = tx.ExecContext(ctx, query, configId, day)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *Repo) UpdateConfig(
	ctx context.Context, id int, repeatWeek int,
	toBeInsert []int, toBeDelete []int,
) error {

	log.Println("UpdateConfig", id, repeatWeek, toBeInsert, toBeDelete)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := repo.db.Rebind(
		`update sales_route_config set repeat_week = ? where id = ?`)
	_, err = tx.ExecContext(ctx, query, repeatWeek, id)

	query = repo.db.Rebind(
		`insert into sales_route_config_day(config_id, day) values(?, ?)`)

	for _, day := range toBeInsert {
		_, err = tx.ExecContext(ctx, query, id, day)
		if err != nil {
			return err
		}
	}

	query = repo.db.Rebind(
		`delete from sales_route_config_day where config_id = ? and day = ?`)
	for _, day := range toBeDelete {
		_, err = tx.ExecContext(ctx, query, id, day)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (repo *Repo) DeleteConfig(
	ctx context.Context, id int) error {

	log.Println("DeleteConfig", id)

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	query := "delete from sales_route_config_day where config_id = ?"
	query = repo.db.Rebind(query)
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	query = repo.db.Rebind(`delete from sales_route_config where id = ?`)
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (repo *Repo) InsertSchedule(
	ctx context.Context, planningId int, customerId uuid.UUID,
	salesmanId uuid.UUID, configId int,
) error {
	log.Println("InsertSchedule", planningId, customerId, salesmanId, configId)

	query := repo.db.Rebind(
		`insert into sales_route_detail(
        config_id, planning_period_id, customer_id, salesman_id)
		values (?, ?, ?, ?)
		`)
	_, err := repo.db.ExecContext(ctx, query, configId, planningId, customerId, salesmanId)

	return err
}

func (repo *Repo) ViewSchedule(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int,
) (int, []ClientSchedule, error) {
	log.Println("ClientSchedule", sortedBy, sortOrder, page, pageSize)

	count := 0
	result := make([]ClientSchedule, 0)

	query := `select count(*) from sales_route_detail`
	err := repo.db.GetContext(ctx, &count, query)
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
	group by srd.id, c.repeat_week, customer."name", u.username, pp.from_date, pp.thru_date
	order by %s %s
	limit ? offset ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, pageSize, page*pageSize)
	return count, result, err
}

func (repo *Repo) DeleteSchedule(
	ctx context.Context, id int) error {

	log.Println("DeleteSchedule", id)

	query := "delete from sales_route_detail where id = ?"
	query = repo.db.Rebind(query)
	_, err := repo.db.ExecContext(ctx, query, id)
	return err
}

func (repo *Repo) GetSchedule(
	ctx context.Context, id int) (ScheduleDetail, error) {

	log.Println("ScheduleDetail", id)

	query := `select srd.id, srd.config_id, c.repeat_week, 
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
	where srd.id = ?
	group by srd.id, c.repeat_week, customer."name", u.username, pp.from_date, pp.thru_date`

	query = repo.db.Rebind(query)

	scheduleDetail := ScheduleDetail{}
	err := repo.db.GetContext(ctx, &scheduleDetail, query, id)

	return scheduleDetail, err
}
