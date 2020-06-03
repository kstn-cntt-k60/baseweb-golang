package salesman

import (
	"time"
)

type ClientSchedule struct {
	Id           int       `json:"id" db:"id"`
	PlanningId   int       `json:"planningId" db:"planning_id"`
	FromDate     time.Time `json:"fromDate" db:"from_date"`
	ThruDate     time.Time `json:"thruDate" db:"thru_date"`
	ConfigId     int       `json:"configId" db:"config_id"`
	RepeatWeek   int       `json:"repeatWeek" db:"repeat_week"`
	DayList      string    `json:"dayList" db:"day_list"`
	StoreName    string    `json:"storeName" db:"store_name"`
	Address      string    `json:"address" db:"address"`
	CustomerName string    `json:"customerName" db:"customer_name"`
	CreatedAt    time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt    time.Time `json:"updatedAt" db:"updated_at"`
	SalesmanName string    `json:"salesmanName" db:"salesman_name"`
}

type ClientCheckinHistory struct {
	Id           int       `json:"id" db:"id"`
	PlanningId   int       `json:"planningId" db:"planning_id"`
	FromDate     time.Time `json:"fromDate" db:"from_date"`
	ThruDate     time.Time `json:"thruDate" db:"thru_date"`
	ConfigId     int       `json:"configId" db:"config_id"`
	StoreName    string    `json:"storeName" db:"store_name"`
	Address      string    `json:"address" db:"address"`
	CustomerName string    `json:"customerName" db:"customer_name"`
	CheckinTime  time.Time `json:"checkinTime" db:"checkin_time"`
}
