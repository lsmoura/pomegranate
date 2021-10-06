package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/lsmoura/humantoken"
	"io/ioutil"
	"os"
	"testing"
)

type TestStruct struct {
	ID     string
	Name   string
	Number int32
}

func (t TestStruct) Kind() string {
	return "test"
}

func (t TestStruct) GetKey() Key {
	return []byte(t.ID)
}

func (t TestStruct) SetKey(key Key) {
	t.ID = string(key)
}

func TestDatabase(t *testing.T) {
	tmp, err := ioutil.TempDir("", "dbtest")
	if err != nil {
		t.Fatalf("error creating temp dir: %v", err)
	}

	defer func() {
		if err := os.RemoveAll(tmp); err != nil {
			t.Errorf("error cleaning temp dir: %v", err)
		}
	}()

	dbPath := fmt.Sprintf("%s/database.db", tmp)
	db, err := Open(dbPath)
	if err != nil {
		t.Fatalf("error opening database %s: %v", dbPath, err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			t.Errorf("error closing database: %v", err)
		}
	}()

	bucketName := (TestStruct{}).Kind()
	if err := db.CreateBucket(bucketName); err != nil {
		t.Fatalf("db.CreateBucket: %s", err)
	}

	t.Run("Basic", func(t *testing.T) {
		testBasic(t, db)
	})
	t.Run("StoreParameters", func(t *testing.T) {
		testStoreParameters(t, db, TestStruct{})
	})
	t.Run("Store_FindByID", func(t *testing.T) {
		testStore_FindByID(t, db)
	})
	t.Run("Store_FindAll", func(t *testing.T) {
		testStore_FindAll(t, db)
	})
}

func testBasic(t *testing.T, db *DB) {
	bucketName := (TestStruct{}).Kind()

	k := humantoken.Generate(8, nil)

	{
		v := TestStruct{
			ID:     k,
			Name:   "John Doe",
			Number: 42,
		}
		vBytes, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("json.Marshal: %s", err)
		}

		if err := db.Store(bucketName, []byte(k), vBytes); err != nil {
			t.Fatalf("db.Store: %s", err)
		}
	}

	// Read K
	{
		var v TestStruct
		bytes, err := db.Read([]byte(bucketName), []byte(k))
		if err != nil {
			t.Fatalf("db.Read: %s", err)
		}
		if err := json.Unmarshal(bytes, &v); err != nil {
			t.Fatalf("json.Unmarshal: %s", err)
		}

		if v.Name != "John Doe" || v.Number != 42 {
			t.Fatalf("Received data differs from written data")
		}
	}

	// Write more data
	{
		data := map[string]TestStruct{
			humantoken.Generate(8, nil): {
				Name:   "Grimoire Noir",
				Number: 88,
			},
			humantoken.Generate(8, nil): {
				Name:   "Jack Sparrow",
				Number: 9,
			},
		}
		for k, v := range data {
			v.ID = k
			vBytes, err := json.Marshal(v)
			if err != nil {
				t.Fatalf("json.Marshal: %s", err)
			}
			if err := db.Store(bucketName, []byte(k), vBytes); err != nil {
				t.Fatalf("db.Store: %s", err)
			}
		}
	}
}

func testStoreParameters(t *testing.T, db *DB, model Model) {
	store := NewStore(db, model)

	if err := store.FindAll(context.TODO(), nil); err == nil {
		t.Fatalf("store.FindAll should return an error if destination is nil")
	}

	intValue := int(33)
	if err := store.FindAll(context.TODO(), &intValue); err == nil {
		t.Fatalf("store.FindAll should return an error if destination does not point to a slice")
	}
}

func testStore_FindByID(t *testing.T, db *DB) {
	store := NewStore(db, TestStruct{})
	k := humantoken.Generate(8, nil)

	{
		v := TestStruct{
			ID:     k,
			Name:   "John Doe",
			Number: 42,
		}
		vBytes, err := json.Marshal(v)
		if err != nil {
			t.Fatalf("json.Marshal: %s", err)
		}

		if err := db.Store(v.Kind(), []byte(k), vBytes); err != nil {
			t.Fatalf("db.Store: %s", err)
		}
	}
	{
		var v TestStruct
		if err := store.FindByID(context.Background(), &v, k); err != nil {
			t.Fatalf("store.FindByID: %s", err)
		}
		if v.Name != "John Doe" || v.Number != 42 {
			t.Fatalf("Received data differs from written data")
		}
	}

	{
		k := humantoken.Generate(8, nil)
		var v TestStruct
		if err := store.FindByID(context.Background(), &v, k); err == nil {
			t.Fatalf("store.FindByID should return error for non existent key")
		}
	}
}

func testStore_FindAll(t *testing.T, db *DB) {
	store := NewStore(db, TestStruct{})

	var v []TestStruct
	if err := store.FindAll(context.Background(), &v); err != nil {
		t.Fatalf("store.FindAll: %s", err)
	}
}
