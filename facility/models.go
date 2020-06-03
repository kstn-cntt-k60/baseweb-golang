package facility

import (
	"time"

	"github.com/google/uuid"
)

type Warehouse struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type CustomerStore struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Customer  string    `json:"customer" db:"customer_name"`
	Address   string    `json:"address" db:"address"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

type SimpleCustomer struct {
	Id   uuid.UUID `json:"id" db:"id"`
	Name string    `json:"name" db:"name"`
}

type InsertStore struct {
	Id         uuid.UUID `json:"id" db:"id"`
	Name       string    `json:"name" db:"name"`
	CustomerId string    `json:"customerId" db:"customer_id"`
	Address    string    `json:"address" db:"address"`
	Latitude   float32   `json:"latitude" db:"latitude"`
	Longitude  float32   `json:"longitude" db:"longitude"`
}
