package database

import (
	"encoding/json"
	"fmt"
	bolt "go.etcd.io/bbolt"
)

const MovieBucketName = "movies"

type DB struct {
	Database *bolt.DB
}

func Open(path string) (DB, error) {
	db, err := bolt.Open(path, 0660, nil)
	if err != nil {
		return DB{}, fmt.Errorf("bolt.Open: %w", err)
	}

	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(MovieBucketName))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	if err != nil {
		return DB{}, fmt.Errorf("db.Update: %w", err)
	}

	return DB{db}, nil
}

func (db DB) Store(bucket string, key []byte, data []byte) error {
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

func (db DB) BucketKeys(bucketName string) ([][]byte, error) {
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

func (db DB) Movie(key string) (Movie, error) {
	if db.Database == nil {
		return Movie{}, fmt.Errorf("database was not initialized")
	}

	var movie Movie

	err := db.Database.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(MovieBucketName))
		v := b.Get([]byte(key))
		if v != nil {
			if err := json.Unmarshal(v, &movie); err != nil {
				return fmt.Errorf("json.Unmarshal: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return Movie{}, fmt.Errorf("Database.View: %w", err)
	}

	return movie, nil
}
