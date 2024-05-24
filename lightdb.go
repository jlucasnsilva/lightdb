package lightdb

import (
	"database/sql"
	"reflect"
	"time"
)

type (
	DB struct {
		conn      *sql.DB
		tableName string
	}

	Item[T any] struct {
		ID        int64
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt time.Time
		Data      *T
		Type      string
		OwnerID   int64
		Ownership string
	}

	ItemValue[T any] struct {
		ID        int64
		Field     string
		Value     T
		ItemType  string
		ItemID    int64
		IsIndexed bool
		IsUnique  bool
	}

	ItemRelationship[T any] struct {
		ID        int64
		Name      string
		OwnerType string
		OwnerID   int64
		ItemType  string
		ItemID    int64
		Data      *T
	}
)

func New(fileName, tableName string) (*DB, error) {
	conn, err := sql.Open("sqlite3", fileName)
	if err != nil {
		return nil, err
	}
	return &DB{
		conn:      conn,
		tableName: tableName,
	}, nil
}

func InsertItems[T any](db *DB, items ...Item[T]) error {
	return nil
}

func isValidEntity(e any) bool {
	zero := reflect.Value{}
	elem := reflect.ValueOf(e).Elem()
	timeType := reflect.TypeOf(time.Time{})

	id := elem.FieldByName("ID")
	createdAt := elem.FieldByName("CreatedAt")
	updatedAt := elem.FieldByName("UpdateAt")
	deletedAt := elem.FieldByName("DeletedAt")
	data := elem.FieldByName("Data")
	typeField := elem.FieldByName("Type")
	ownerID := elem.FieldByName("OwnerID")
	ownership := elem.FieldByName("Ownership")

	return id != zero &&
		id.Kind() == reflect.Int64 &&
		createdAt != zero &&
		createdAt.Type() == timeType &&
		updatedAt != zero &&
		updatedAt.Type() == timeType &&
		deletedAt != zero &&
		deletedAt.Type() == timeType &&
		data != zero &&
		typeField != zero &&
		typeField.Kind() != reflect.String &&
		ownerID != zero &&
		ownerID.Kind() == reflect.Int64 &&
		ownership != zero &&
		ownership.Kind() == reflect.String
}
