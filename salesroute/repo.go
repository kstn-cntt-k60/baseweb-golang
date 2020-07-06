package salesroute

import (
	"baseweb/kmeans"
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
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `insert into salesman(id, created_by_user_login_id)
	values (:id, :created_by_user_login_id)`

	_, err = tx.NamedExecContext(ctx, query, salesman)
	if err != nil {
		return err
	}

	query = `insert into user_login_security_group(user_login_id, security_group_id)
	values (?, 8)`
	query = repo.db.Rebind(query)
	_, err = tx.ExecContext(ctx, query, salesman.Id)
	if err != nil {
		return err
	}

	return tx.Commit()
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
	p.last_name || ' ' || p.middle_name || ' ' || p.first_name as salesman_name,
	u2.username as created_by,
	s.created_at, s.updated_at
	from salesman s
	inner join user_login u1 on u1.id = s.id
	inner join user_login u2 on u2.id = s.created_by_user_login_id
	inner join person p on p.id = u1.person_id
	order by %s %s
	limit ? offset ?`

	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)
	err = repo.db.SelectContext(ctx, &result, query, pageSize, page*pageSize)
	return count, result, err
}

func (repo *Repo) ViewSalesmanWithFullName(ctx context.Context,
	sortedBy, sortOrder string, page, pageSize int,
	fullName string,
) (int, []ClientSalesman, error) {

	log.Println("ViewSalesman", page, pageSize,
		sortedBy, sortOrder, fullName)

	var count int = 0
	result := make([]ClientSalesman, 0)

	query := repo.db.Rebind(`
        select count(*) from salesman s, person p, user_login ul
        where p.id = ul.person_id and s.id = ul.id and full_name_tsvector @@ plainto_tsquery(vn_unaccent(?))
        `)
	err := repo.db.GetContext(ctx, &count, query, fullName)
	if err != nil {
		return count, result, err
	}

	query = `	select s.id, u1.username, 
	p.last_name || ' ' || p.middle_name || ' ' || p.first_name as salesman_name,
	u2.username as created_by,
	s.created_at, s.updated_at
	from salesman s
	inner join user_login u1 on u1.id = s.id
	inner join user_login u2 on u2.id = s.created_by_user_login_id
	inner join person p on p.id = u1.person_id
	where p.full_name_tsvector @@ plainto_tsquery(vn_unaccent(?))
	order by %s %s
	limit ? offset ?`
	query = fmt.Sprintf(query, sortedBy, sortOrder)
	query = repo.db.Rebind(query)

	err = repo.db.SelectContext(ctx, &result,
		repo.db.Rebind(query), fullName, pageSize, page*pageSize)

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
	ctx context.Context, planningId int, salesmanId uuid.UUID,
	customerStores []Store,
) error {
	log.Println("InsertSchedule", planningId, salesmanId, customerStores)

	for _, store := range customerStores {
		query := repo.db.Rebind(
			`insert into sales_route_detail(
			config_id, planning_period_id, customer_store_id, salesman_id)
			values (?, ?, ?, ?)
			`)
		_, err := repo.db.ExecContext(ctx, query, store.ConfigId, planningId, store.StoreId, salesmanId)
		if err != nil {
			return err
		}
	}

	return nil
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
	pp.from_date, pp.thru_date, p.last_name || ' ' || p.middle_name || ' ' || p.first_name as salesman_name,
	facility."name" as store_name,
	facility.address, customer."name" as customer_name,
	srd.created_at, srd.updated_at
	from sales_route_detail srd
	inner join sales_route_planning_period pp on srd.planning_period_id = pp.id
	inner join sales_route_config c on srd.config_id = c.id
	inner join sales_route_config_day d on d.config_id = c.id
	inner join salesman on srd.salesman_id = salesman.id
	inner join user_login u on salesman.id = u.id
	inner join facility_customer fc on fc.id = srd.customer_store_id
	inner join customer customer on customer.id = fc.customer_id
	inner join facility on facility.id = fc.id
	inner join person p on p.id = u.person_id
	group by srd.id, c.repeat_week,  u.username, pp.from_date, pp.thru_date, p.last_name,
	p.middle_name, p.first_name, 
	facility."name", facility.address, customer."name"
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

	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := `delete from salesman_checkin_history
		where sales_route_detail_id = ?`
	query = repo.db.Rebind(query)
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	query = "delete from sales_route_detail where id = ?"
	query = repo.db.Rebind(query)
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return tx.Commit()
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

func (repo *Repo) ViewClustering(ctx context.Context, nCluster int, city string,
) ([]ClientNeighbor, error) {
	log.Println("View CustomerStore Kmeans", nCluster)

	result := make([]ClientNeighbor, 0)
	if city == "hanoi" {
		city = "%Hà Nội%"
	} else if city == "hcm" {
		city = "%Hồ Chí Minh%"
	}

	type DBPoint struct {
		Id           uuid.UUID `db:"id"`
		StoreName    string    `db:"store_name"`
		Address      string    `db:"address"`
		CustomerName string    `db:"customer_name"`
		Lat          float32   `db:"latitude"`
		Long         float32   `db:"longitude"`
	}

	var dbPoints []DBPoint

	if city == "" {
		query := `select f.id, f.name store_name, f.address, c."name" customer_name, 
	fc.latitude, fc.longitude 
	from facility as f, facility_customer as fc, 
	customer as c
	where f.id = fc.id
	and fc.customer_id = c.id`
		err := repo.db.SelectContext(ctx, &dbPoints, query)
		if err != nil {
			return result, err
		}
	} else {
		query := `select f.id, f.name store_name, f.address, c."name" customer_name, 
	fc.latitude, fc.longitude 
	from facility as f, facility_customer as fc, 
	customer as c
	where f.id = fc.id
	and fc.customer_id = c.id and f.address like ?`
		query = repo.db.Rebind(query)
		err := repo.db.SelectContext(ctx, &dbPoints, query, city)
		if err != nil {
			return result, err
		}
	}

	points := make([]kmeans.Point, 0)
	for _, p := range dbPoints {
		points = append(points, kmeans.Point{Lat: p.Lat, Long: p.Long})
	}

	_, clusterCustomerStores := kmeans.Kmeans(points, nCluster)

	for index, nb := range clusterCustomerStores {
		result = append(result, ClientNeighbor{
			Id: dbPoints[index].Id, Index: nb.Index,
			StoreName:    dbPoints[index].StoreName,
			Address:      dbPoints[index].Address,
			CustomerName: dbPoints[index].CustomerName,
			Lat:          nb.Point.Lat, Long: nb.Point.Long,
		})
	}

	log.Println("COMPLETE", result)

	return result, nil
}

func (repo *Repo) ViewStoreOfSalesman(ctx context.Context, salesmanId uuid.UUID,
) ([]StoreOfSalesman, error) {
	log.Println("ViewStoreOfSalesman", salesmanId)

	result := make([]StoreOfSalesman, 0)

	query := `select srd.customer_store_id, fc.latitude, fc.longitude, 
		f.name as customer_store_name, c.name as customer_name
		from sales_route_detail srd, facility_customer fc, facility f, customer c
		where srd.customer_store_id = fc.id and f.id = fc.id 
		and c.id = fc.customer_id and 
		srd.salesman_id = ?`

	query = fmt.Sprintf(query)
	query = repo.db.Rebind(query)
	err := repo.db.SelectContext(ctx, &result, query, salesmanId)
	return result, err
}

func (repo *Repo) GetPairStoreSalesman(ctx context.Context,
) ([]PairStoreSalesmanId, error) {
	log.Println("GetPairStoreSalesman")

	result := make([]PairStoreSalesmanId, 0)

	query := `select srd.customer_store_id, srd.salesman_id from sales_route_detail srd `

	query = fmt.Sprintf(query)
	query = repo.db.Rebind(query)
	err := repo.db.SelectContext(ctx, &result, query)
	return result, err
}
