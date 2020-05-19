package salesroute

import (
	"time"

	"github.com/google/uuid"
)

type Salesman struct {
	Id        uuid.UUID `json:"id" db:"id"`
	CreatedBy uuid.UUID `json:"createdBy" db:"created_by_user_login_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ClientSalesman struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	CreatedBy string    `json:"createdBy" db:"created_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type ClientUserLogin struct {
	Id         uuid.UUID `json:"id" db:"id"`
	Username   string    `json:"username" db:"username"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
	FirstName  string    `json:"firstName" db:"first_name"`
	MiddleName string    `json:"middleName" db:"middle_name"`
	LastName   string    `json:"lastName" db:"last_name"`
	BirthDate  string    `json:"birthDate" db:"birth_date"`
	GenderId   int16     `json:"genderId" db:"gender_id"`
}

type ClientPlanningPeriod struct {
	Id        int       `json:"id" db:"id"`
	FromDate  time.Time `json:"fromDate" db:"from_date"`
	ThruDate  time.Time `json:"thruDate" db:"thru_date"`
	CreatedBy string    `json:"createdBy" db:"created_by"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type PlanningPeriod struct {
	Id        int       `json:"id" db:"id"`
	FromDate  time.Time `json:"fromDate" db:"from_date"`
	ThruDate  time.Time `json:"thruDate" db:"thru_date"`
	CreatedBy uuid.UUID `db:"created_by_user_login_id"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type SalesRouteConfig struct {
	Id         int `json:"id" db:"id"`
	RepeatWeek int `json:"repeatWeek" db:"repeat_week"`
	// Day []  `json:"days"`
	CreatedBy uuid.UUID `db:"created_by_user_login_id"`
}

type ClientSalesRouteConfig struct {
	Id         int       `json:"id" db:"id"`
	RepeatWeek int       `json:"repeatWeek" db:"repeat_week"`
	DayList    string    `json:"dayList" db:"day_list"`
	CreatedBy  string    `json:"createdBy" db:"created_by"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type ClientSchedule struct {
	Id           int       `json:"id" db:"id"`
	PlanningId   int       `json:"planningId" db:"planning_id"`
	CustomerName string    `json:"customerName" db:"customer_name"`
	SalesmanName string    `json:"salesmanName" db:"salesman_name"`
	ConfigId     int       `json:"configId" db:"config_id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	FromDate     time.Time `json:"fromDate" db:"from_date"`
	ThruDate     time.Time `json:"thruDate" db:"thru_date"`
	RepeatWeek   int       `json:"repeatWeek" db:"repeat_week"`
	DayList      string    `json:"dayList" db:"day_list"`
}

type ScheduleDetail struct {
	Id           int       `json:"id" db:"id"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	PlanningId   int       `json:"planningId" db:"planning_id"`
	FromDate     time.Time `json:"fromDate" db:"from_date"`
	ThruDate     time.Time `json:"thruDate" db:"thru_date"`
	CustomerName string    `json:"customerName" db:"customer_name"`
	SalesmanName string    `json:"salesmanName" db:"salesman_name"`
	ConfigId     int       `json:"configId" db:"config_id"`
	RepeatWeek   int       `json:"repeatWeek" db:"repeat_week"`
	DayList      string    `json:"dayList" db:"day_list"`
}
