package u

import (
	"reflect"
	"strconv"
	"strings"
	"unicode"
)

func Check(err error) {
	if err != nil {
		panic(err)
	}
}

func Ptr[T any](t T) *T {
	return &t
}

func Map[T any, U any](ts []T, f func(T) U) []U {
	us := make([]U, len(ts))
	for i, t := range ts {
		us[i] = f(t)
	}
	return us
}

func FilterBy(filterStruct any) (string, []any) {
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
