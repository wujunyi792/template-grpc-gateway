package dtoTool

import (
	"reflect"
	"strings"
)

const (
	FromQueryTag = "search"
)

/*
*
tag: "[column]:[type]"
*/
func TransferSql(q any) map[string]any {
	qType := reflect.TypeOf(q)
	qValue := reflect.ValueOf(q)
	for qType.Kind() == reflect.Pointer {
		qType = qType.Elem()
		qValue = qValue.Elem()
	}
	query := make(map[string]any)
	var tag string
	var ok bool
	for i := 0; i < qType.NumField(); i++ {
		if qValue.Field(i).Kind() != reflect.String || qValue.Field(i).Interface().(string) == "" {
			continue
		}
		tag, ok = "", false
		tag, ok = qType.Field(i).Tag.Lookup(FromQueryTag)
		if !ok {
			continue
		}
		t := strings.Split(tag, ";")
		if len(t) < 1 {
			continue
		}
		if len(t) == 1 {
			query[t[0]+" = ?"] = qValue.Field(i).Interface()
		}
		if len(t) == 2 {
			switch t[1] {
			case "in":
				query[t[0]+" LIKE ?"] = "%" + qValue.Field(i).Interface().(string) + "%"
			}
		}
	}
	return query
}
