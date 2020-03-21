package account

import (
	"time"

	"github.com/google/uuid"
)

type Party struct {
	Id          uuid.UUID `json:"uuid" db:"id"`
	TypeId      int16     `json:"partyTypeId" db:"party_type_id"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	CreatedBy   uuid.UUID `json:"createdBy" db:"created_by_user_login_id"`
	UpdatedBy   uuid.UUID `json:"updatedBy" db:"updated_by_user_login_id"`
}

type Person struct {
	Id         uuid.UUID `json:"id" db:"id"`
	FirstName  string    `json:"firstName" db:"first_name"`
	MiddleName string    `json:"middleName" db:"middle_name"`
	LastName   string    `json:"lastName" db:"last_name"`
	BirthDate  string    `json:"birthDate" db:"birth_date"`
	GenderId   int16     `json:"genderId" db:"gender_id"`
	CreatedAt  time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt  time.Time `json:"updatedAt" db:"updated_at"`
}

type Customer struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

type ClientPerson struct {
	Id          uuid.UUID `json:"id" db:"id"`
	FirstName   string    `json:"firstName" db:"first_name"`
	MiddleName  string    `json:"middleName" db:"middle_name"`
	LastName    string    `json:"lastName" db:"last_name"`
	BirthDate   string    `json:"birthDate" db:"birth_date"`
	GenderId    int16     `json:"genderId" db:"gender_id"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Description string    `json:"description" db:"description"`
}

type ClientCustomer struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
	Description string    `json:"description" db:"description"`
}

type SimplePerson struct {
	Id         uuid.UUID `json:"id" db:"id"`
	FirstName  string    `json:"firstName" db:"first_name"`
	MiddleName string    `json:"middleName" db:"middle_name"`
	LastName   string    `json:"lastName" db:"last_name"`
	BirthDate  string    `json:"birthDate" db:"birth_date"`
	GenderId   int16     `json:"genderId" db:"gender_id"`
}

type UserLogin struct {
	Id        uuid.UUID `json:"id" db:"id"`
	Username  string    `json:"username" db:"username"`
	Password  string    `json:"password" db:"password"`
	PersonId  uuid.UUID `json:"personId" db:"person_id"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}
