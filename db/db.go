package db

import (
	"context"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (commandTag pgconn.CommandTag, err error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

func filterBy(filterStruct any) (string, []any) {
	var filters []string
	var params []any
	val := reflect.ValueOf(filterStruct).Elem()

	for i := 0; i < val.NumField(); i++ {
		value := val.Field(i)
		if value.IsNil() {
			continue
		}

		fieldType := val.Type().Field(i)

		sqlName := fieldType.Tag.Get("db")
		if sqlName == "-" {
			continue
		}
		if sqlName == "" {
			sqlName = string(unicode.ToLower(rune(fieldType.Name[0]))) + fieldType.Name[1:]
		}

		operator := fieldType.Tag.Get("op")
		if operator == "" {
			operator = "="
		}

		filters = append(filters, ``+sqlName+` `+operator+" $"+strconv.Itoa(len(filters)+1))
		params = append(params, value.Interface())
	}

	if len(filters) == 0 {
		return "", params
	}

	return "WHERE " + strings.Join(filters, " AND "), params
}
