package database

import "context"

type Key []byte

type Model interface {
	Kind() string

	GetKey() Key
	SetKey(Key)
}

type store struct {
	session *DB
	model Model
}

func NewStore(sess *DB, model Model) *store {
	return &store{
		session: sess,
		model:   model,
	}
}

func (s *store) FindOne(ctx context.Context, dst interface{}) error  {
	return nil
}

func (s *store) FindAll(ctx context.Context, dst interface{}) error {
	return nil
}
