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

func TestInit(t *testing.T) {
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

	store := NewStore(db, TestStruct{})
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
		var v []TestStruct
		if err := store.FindAll(context.Background(), &v); err != nil {
			t.Fatalf("store.FindAll: %s", err)
		}
		if len(v) != 3 {
			t.Fatalf("received a different number of elements than expected: %d", len(v))
		}
	}
}
