package database

import (
	"context"
	"github.com/pkg/errors"
	bolt "go.etcd.io/bbolt"
)

// TODO: remove this from here and use what the model says
const MovieBucketName = "movie"

type DB struct {
	Database *bolt.DB
}

type Store interface {
	FindByID(ctx context.Context, dst interface{}, id string) error
	FindAll(ctx context.Context, dst interface{}, filters ...Filter) error
	//FindOne(ctx context.Context, dst interface{}, filters ...Filter) error
}

func Open(path string) (*DB, error) {
	db, err := bolt.Open(path, 0660, nil)
	if err != nil {
		return nil, errors.Wrap(err, "bolt.Open")
	}

	err = db.Update(func(tx *bolt.Tx) error {
		// TODO: the database package should not be creating stuff on its own
		_, err := tx.CreateBucketIfNotExists([]byte(MovieBucketName))
		if err != nil {
			return errors.Wrap(err, "create bucket")
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "db.Update")
	}

	return &DB{db}, nil
}

func (db *DB) Close() error {
	return db.Database.Close()
}

func (db *DB) CreateBucket(bucketName string) error {
	err := db.Database.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(bucketName))
		if err != nil {
			return errors.Wrap(err, "tx.CreateBucketIfNotExists")
		}
		return nil
	})

	if err != nil {
		return errors.Wrap(err, "db.Database.Update")
	}

	return nil
}

func (db *DB) Store(bucket string, key []byte, data []byte) error {
	if db.Database == nil {
		return errors.New("database was not initialized")
	}
	err := db.Database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, data)
		if err != nil {
			return errors.Wrap(err, "bucket.Put")
		}

		return nil
	})

	if err != nil {
		return errors.Wrap(err, "db.Update")
	}

	return nil
}

func (db *DB) Read(bucket []byte, key []byte) ([]byte, error) {
	if db.Database == nil {
		return nil, errors.New("database was not initialized")
	}

	var retVal []byte

	err := db.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		if b == nil {
			return errors.Errorf("cannot find bucket: %s", bucket)
		}
		retVal = b.Get(key)
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "Database.View")
	}

	return retVal, nil
}

func (db *DB) BucketKeys(bucketName string) ([][]byte, error) {
	if db.Database == nil {
		return nil, errors.New("database was not initialized")
	}

	var resp [][]byte

	err := db.Database.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(bucketName))

		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			resp = append(resp, k)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return resp, nil
}
