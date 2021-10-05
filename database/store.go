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

	if err := json.Unmarshal(v, &dst); err != nil {
		return fmt.Errorf("json.Unmarshal: %w", err)
	}

	return nil
}

func (s *store) FindAll(ctx context.Context, dst interface{}, filters ...Filter) error {
	if s.db() == nil {
		return fmt.Errorf("database was not initialized")
	}

	// TODO: Check if dst is actually a slice
	// TODO: Check if dst elements are of the same type of s.model

	bucketName := s.model.Kind()
	myType := reflect.TypeOf(s.model)

	slice := reflect.ValueOf(dst).Elem()

	err := s.db().View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()

		for k, bytes := c.First(); k != nil; k, bytes = c.Next() {
			m := reflect.New(myType)
			err := json.Unmarshal(bytes, m.Interface())
			// TODO: apply filter
			if err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}
			slice = reflect.Append(slice, m.Elem())
		}

		return nil
	})
	if err != nil {
		return fmt.Errorf("Database.View: %w", err)
	}

	reflect.ValueOf(dst).Elem().Set(slice)

	return nil
}
