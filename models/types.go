package models

import (
	"database/sql/driver"
)

type AnyStruct interface {
	Value() (driver.Value, error)
	Scan(src any) error
}
