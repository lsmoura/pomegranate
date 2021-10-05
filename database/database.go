package database

import (
	"context"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

const MovieBucketName = "movies"

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
		return nil, fmt.Errorf("bolt.Open: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(MovieBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("db.Update: %w", err)
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
			return fmt.Errorf("tx.CreateBucketIfNotExists: %w", err)
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("db.Database.Update: %w", err)
	}

	return nil
}

func (db *DB) Store(bucket string, key []byte, data []byte) error {
	if db.Database == nil {
		return fmt.Errorf("database was not initialized")
	}
	err := db.Database.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		err := b.Put(key, data)
		if err != nil {
			return fmt.Errorf("bucket.Put: %w", err)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("db.Update: %w", err)
	}

	return nil
}

func (db *DB) Read(bucket []byte, key []byte) ([]byte, error) {
	if db.Database == nil {
		return nil, fmt.Errorf("database was not initialized")
	}

	var retVal []byte

	err := db.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket(bucket)
		retVal = b.Get(key)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("Database.View: %w", err)
	}

	return retVal, nil
}

func (db *DB) BucketKeys(bucketName string) ([][]byte, error) {
	if db.Database == nil {
		return nil, fmt.Errorf("database was not initialized")
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
