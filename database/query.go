package database

import (
	"fmt"
	"github.com/pkg/errors"
	"reflect"
)

type operator int

const (
	LessThan      operator = iota // 0
	LessThanEq                    // 1
	Equal                         // 2
	GreaterThan                   // 3
	GreaterThanEq                 // 4
)

type Filter map[string]Comparison

type Comparison struct {
	Operator operator
	Value    interface{}
}

func (c Comparison) compare(value interface{}) (bool, error) {
	typeOfValue := reflect.TypeOf(value)
	if typeOfValue.Kind() != reflect.TypeOf(c.Value).Kind() {
		return false, errors.New(fmt.Sprintf("cannot compare different types: %s, %s", typeOfValue.Kind(), reflect.TypeOf(c.Value).Kind()))
	}

	switch typeOfValue.Kind() {
	case reflect.String:
		if c.Operator == Equal {
			return value.(string) == c.Value.(string), nil
		}
		return false, errors.New(fmt.Sprintf("comparison not implemented: %s, %d", typeOfValue.Kind(), c.Operator))
	case reflect.Int32:
		switch c.Operator {
		case Equal:
			return c.Value.(int32) == value.(int32), nil
		case LessThan:
			return c.Value.(int32) < value.(int32), nil
		case LessThanEq:
			return c.Value.(int32) <= value.(int32), nil
		case GreaterThan:
			return c.Value.(int32) > value.(int32), nil
		case GreaterThanEq:
			return c.Value.(int32) >= value.(int32), nil
		}
		return false, errors.New(fmt.Sprintf("comparison not implemented: %s, %d", typeOfValue.Kind(), c.Operator))
	}

	return false, errors.New(fmt.Sprintf("comparison not implemented: %s, %d", typeOfValue.Kind(), c.Operator))
}
