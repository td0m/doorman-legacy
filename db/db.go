package db

import (
	"context"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/jackc/pgx/v5/pgxpool"
)

var pg *pgxpool.Pool

func Init(ctx context.Context) error {
	var err error
	pg, err = pgxpool.New(ctx, "")
	if err != nil {
		return fmt.Errorf("pgxpool.New failed: %w", err)
	}

	if err := pg.Ping(ctx); err != nil {
		return fmt.Errorf("pg.Ping failed: %w", err)
	}
	return nil
}

func Close() error {
	pg.Close()
	return nil
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
