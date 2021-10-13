package database

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/pkg/errors"
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

type Error string

func (e Error) Error() string {
	return string(e)
}

const NotFoundError = Error("not found")

func (s *store) FindByID(ctx context.Context, dst interface{}, id string) error {
	if s.db() == nil {
		return errors.New("database was not initialized")
	}

	bucketName := s.model.Kind()

	v, err := s.session.Read([]byte(bucketName), []byte(id))
	if err != nil {
		return errors.Wrap(err, "db.Read")
	}

	if v == nil {
		return errors.Wrap(NotFoundError, "db.Read")
	}

	if err := json.Unmarshal(v, &dst); err != nil {
		return errors.Wrap(err, "json.Unmarshal")
	}

	return nil
}

func checkFilters(dst interface{}, filters []Filter) (bool, error) {
	if dst == nil {
		return false, errors.New("cannot compare without an object")
	}
	if filters == nil {
		return true, nil
	}

	valueOf := reflect.ValueOf(dst)
	for _, filterMap := range filters {
		for key, comparison := range filterMap {
			value := valueOf.FieldByName(key)
			result, err := comparison.compare(value.Interface())
			if err != nil {
				return false, errors.Wrap(err, "compare")
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
		return errors.New("database was not initialized")
	}

	if dst == nil {
		return errors.New("dst cannot be nil")
	}
	if kind := reflect.TypeOf(dst).Kind(); kind != reflect.Ptr {
		return errors.New(fmt.Sprintf("dst is not a pointer: %s", kind))
	}
	if ptrKind := reflect.TypeOf(dst).Elem().Kind(); ptrKind != reflect.Slice {
		return errors.New(fmt.Sprintf("dst does not point to a slice: %s", ptrKind))
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
				return errors.Wrap(err, "json.Unmarshal")
			}

			shouldInclude, err := checkFilters(m.Elem().Interface(), filters)
			if err != nil {
				return errors.Wrap(err, "checkFilters")
			}

			if shouldInclude {
				slice = reflect.Append(slice, reflect.Indirect(m.Elem()))
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "Database.View")
	}

	reflect.ValueOf(dst).Elem().Set(slice)

	return nil
}
