package sqlite

import (
	"database/sql"
	"reflect"
)

func rowsToStructs(rows *sql.Rows, dest any) error {
	destv := reflect.ValueOf(dest).Elem()
	args := make([]any, destv.Type().Elem().NumField())

	for rows.Next() {
		rowp := reflect.New(destv.Type().Elem())
		rowv := rowp.Elem()

		for i := 0; i < rowv.NumField(); i++ {
			args[i] = rowv.Field(i).Addr().Interface()
		}

		if err := rows.Scan(args...); err != nil {
			return err
		}
		destv.Set(reflect.Append(destv, rowv))
	}
	return nil
}
