package sabnzbd

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func asString(i interface{}) string {
	switch reflect.TypeOf(i) {
	case reflect.TypeOf(string("")):
		return i.(string)
	case reflect.TypeOf(int(0)):
		return strconv.Itoa(i.(int))
	case reflect.TypeOf(int32(0)):
		return strconv.Itoa(int(i.(int32)))
	case reflect.TypeOf(int64(0)):
		return strconv.Itoa(int(i.(int64)))
	}

	return fmt.Sprintf("<%s>", reflect.TypeOf(i))
}

func InjectQuery(query url.Values, s interface{}) error {
	typeOf := reflect.TypeOf(s)

	if typeOf.Kind() != reflect.Struct {
		return fmt.Errorf("invalid parameter kind: %s", typeOf.Kind())
	}

	valueOf := reflect.ValueOf(s)

	for i := 0; i < typeOf.NumField(); i++ {
		field := typeOf.Field(i)

		tag := field.Tag.Get("query_name")
		pieces := strings.Split(tag, ",")

		if pieces[0] != "-" {
			fieldName := field.Name
			if pieces[0] != "" {
				fieldName = pieces[0]
			}

			value := valueOf.Field(i).Interface()
			query.Set(fieldName, asString(value))
		}
	}

	return nil
}

func InjectInUrl(url *url.URL, s interface{}) error {
	query := url.Query()
	err := InjectQuery(query, s)
	if err != nil {
		return fmt.Errorf("InjectQuery: %w", err)
	}
	url.RawQuery = query.Encode()

	return nil
}
