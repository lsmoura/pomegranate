package database

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	bolt "go.etcd.io/bbolt"
)

type Key []byte

type Model interface {
	Kind() string

	GetKey() Key
	SetKey(Key)
}

type store struct {
	session *DB
	model   Model
}

func NewStore(sess *DB, model Model) *store {
	return &store{
		session: sess,
		model:   model,
	}
}

func (s *store) db() *bolt.DB {
	return s.session.Database
}

func (s *store) FindByID(ctx context.Context, dst interface{}, id string) error {
	if s.db() == nil {
		return fmt.Errorf("database was not initialized")
	}

	bucketName := s.model.Kind()

	v, err := s.session.Read([]byte(bucketName), []byte(id))
	if err != nil {
		return fmt.Errorf("db.Read: %w", err)
	}

	if v == nil {
		return fmt.Errorf("not found")
	}

	if err := json.Unmarshal(v, &dst); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

func compare(value interface{}, c Comparison) (bool, error) {
	typeOfValue := reflect.TypeOf(value)
	if typeOfValue.Kind() != reflect.TypeOf(c.Value).Kind() {
		return false, fmt.Errorf("cannot compare different types: %s, %s", typeOfValue.Kind(), reflect.TypeOf(c.Value).Kind())
	}

	switch typeOfValue.Kind() {
	case reflect.String:
		if c.Operator == Equal {
			return value.(string) == c.Value.(string), nil
		}
		return false, fmt.Errorf("comparison not implemented: %s, %d", typeOfValue.Kind(), c.Operator)
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
		return false, fmt.Errorf("comparison not implemented: %s, %d", typeOfValue.Kind(), c.Operator)
	}

	return false, fmt.Errorf("comparison not implemented: %s, %d", typeOfValue.Kind(), c.Operator)
}

func checkFilters(dst interface{}, filters []Filter) (bool, error) {
	if dst == nil {
		return false, fmt.Errorf("cannot compare without an object")
	}
	if filters == nil {
		return true, nil
	}

	valueOf := reflect.ValueOf(dst)
	for _, filterMap := range filters {
		for key, comparison := range filterMap {
			value := valueOf.FieldByName(key)
			result, err := compare(value.Interface(), comparison)
			if err != nil {
				return false, fmt.Errorf("compare: %w", err)
			}
			if !result {
				return false, nil
			}
		}
	}

	return true, nil
}

func (s *store) FindAll(ctx context.Context, dst interface{}, filters ...Filter) error {
	if s.db() == nil {
		return fmt.Errorf("database was not initialized")
	}

	if dst == nil {
		return fmt.Errorf("dst cannot be nil")
	}
	if kind := reflect.TypeOf(dst).Kind(); kind != reflect.Ptr {
		return fmt.Errorf("dst is not a pointer: %s", kind)
	}
	if ptrKind := reflect.TypeOf(dst).Elem().Kind(); ptrKind != reflect.Slice {
		return fmt.Errorf("dst does not point to a slice: %s", ptrKind)
	}
	myType := reflect.TypeOf(s.model)

	bucketName := s.model.Kind()
	slice := reflect.ValueOf(dst).Elem()

	err := s.db().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()

		for k, bytes := c.First(); k != nil; k, bytes = c.Next() {
			m := reflect.New(myType)
			err := json.Unmarshal(bytes, m.Interface())
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}

			shouldInclude, err := checkFilters(m.Elem().Interface(), filters)
			if err != nil {
				return fmt.Errorf("checkFilters: %w", err)
			}

			if shouldInclude {
				slice = reflect.Append(slice, m.Elem())
			}
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Database.View: %w", err)
	}

	reflect.ValueOf(dst).Elem().Set(slice)

	return nil
}
