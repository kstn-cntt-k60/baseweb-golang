package account

import (
	"time"

	"github.com/google/uuid"
)

type Party struct {
	Id          uuid.UUID `json:"uuid"`
	TypeId      int16     `json:"partyTypeId"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	CreatedBy   uuid.UUID `json:"createdBy"`
	UpdatedBy   uuid.UUID `json:"updatedBy"`
}

type Person struct {
	Id         uuid.UUID `json:"id"`
	FirstName  string    `json:"firstName"`
	MiddleName string    `json:"middleName"`
	LastName   string    `json:"lastName"`
	BirthDate  string    `json:"birthDate"`
	GenderId   int16     `json:"genderId"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Customer struct {
	Id        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type ClientPerson struct {
	Id          uuid.UUID `json:"id"`
	FirstName   string    `json:"firstName"`
	MiddleName  string    `json:"middleName"`
	LastName    string    `json:"lastName"`
	BirthDate   string    `json:"birthDate"`
	GenderId    int16     `json:"genderId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	Description string    `json:"description"`
}
