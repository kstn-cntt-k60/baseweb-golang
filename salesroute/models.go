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
