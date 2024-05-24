package lightdb

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type (
	DB struct {
		conn      *sql.DB
		tableName string
	}

	Ownership struct {
		OwnerID   int64
		Ownership string
	}
)

type (
	itemModel struct {
		ID        int64
		CreatedAt time.Time
		UpdatedAt time.Time
		DeletedAt time.Time
		Data      []byte
		Type      string
		OwnerID   int64
		Ownership string
	}

	itemValueModel struct {
		ID        int64
		Field     string
		Value     any // TODO
		ItemType  string
		ItemID    int64
		IsIndexed bool
		IsUnique  bool
	}

	itemRelationshipModel struct {
		ID        int64
		Name      string
		OwnerType string
		OwnerID   int64
		ItemType  string
		ItemID    int64
		Data      []byte
	}
)

var (
	ErrInvalidItem = errors.New("invalid item: it should be a pointer to struct")
	ErrMarshalling = errors.New("invalid item: JSON marshalling failed")
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

func (db *DB) Insert(typ string, data any, ownership ...Ownership) (int64, error) {
	return db.InsertCtx(context.Background(), typ, data, ownership...)
}

func (db *DB) InsertCtx(
	ctx context.Context,
	typ string,
	data any,
	ownership ...Ownership,
) (int64, error) {
	if isStruct(data) {
		return -1, ErrInvalidItem
	}

	bs, err := json.Marshal(data)
	if err != nil {
		return -1, ErrMarshalling
	}

	var res sql.Result
	if len(ownership) > 0 {
		res, err = db.insertWithOwnershipQ(ctx, typ, bs, ownership[0])
	} else {
		res, err = db.insertQ(ctx, typ, bs)
	}

	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (db *DB) insertQ(ctx context.Context, typ string, data []byte) (sql.Result, error) {
	now := time.Now()
	q := fmt.Sprintf(`
		INSERT INTO %v (created_at, updated_at, type, data)
		VALUES (?, ?, ?, ?);`,
		db.tableName,
	)
	return db.conn.ExecContext(ctx, q, now, now, typ, data)
}

func (db *DB) insertWithOwnershipQ(
	ctx context.Context,
	typ string,
	data []byte,
	o Ownership,
) (sql.Result, error) {
	now := time.Now()
	q := fmt.Sprintf(`
		INSERT INTO %v (created_at, updated_at, type, data, owner_id, ownership)
		VALUES (?, ?, ?, ?);`,
		db.tableName,
	)
	return db.conn.ExecContext(ctx, q, now, now, typ, data, o.OwnerID, o.Ownership)
}

func isStruct(x any) bool {
	v := reflect.ValueOf(x)
	if v.Kind() == reflect.Ptr {
		return false
	}
	return v.Elem().Kind() == reflect.Struct
}
