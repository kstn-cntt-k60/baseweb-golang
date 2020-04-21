package basic

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Handler func(http.ResponseWriter, *http.Request) error

type NullString struct {
	sql.NullString
}

type NullTime struct {
	sql.NullTime
}

func (v NullString) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.String)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullString) UnmarshalJSON(data []byte) error {
	var x *string
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.String = *x
	} else {
		v.Valid = false
	}
	return nil
}

func (v NullTime) MarshalJSON() ([]byte, error) {
	if v.Valid {
		return json.Marshal(v.Time)
	} else {
		return json.Marshal(nil)
	}
}

func (v *NullTime) UnmarshalJSON(data []byte) error {
	var x *time.Time
	if err := json.Unmarshal(data, &x); err != nil {
		return err
	}
	if x != nil {
		v.Valid = true
		v.Time = *x
	} else {
		v.Valid = false
	}
	return nil
}

func LowerBound(array []int, value int) int {
	first := 0
	last := len(array)

	for first != last {
		mid := (first + last) / 2
		if array[mid] < value {
			first = mid + 1
		} else {
			last = mid
		}
	}

	return first
}
