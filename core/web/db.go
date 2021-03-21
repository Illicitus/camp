package web

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Model interface {
	IsGormModel() bool
}

type DB struct {
	Conn   *gorm.DB
	Models []interface{}
}

func NewDbConnection(dialect, connectionInfo string, mode bool) (*DB, error) {
	db, err := gorm.Open(dialect, connectionInfo)
	if err != nil {
		return nil, err
	}
	return &DB{Conn: db}, nil
}

func (db *DB) detectModels() error {
	if len(db.Models) == 0 {
		return errors.New("Models list can't be empty. Please connect db to app first.")
	} else {
		for _, m := range db.Models {
			// TODO check orm type
			fmt.Println(m)
		}
	}

	return nil
}

func (db *DB) DestructiveReset() error {

	if err := db.detectModels(); err != nil {
		return err
	}

	if err := db.Conn.DropTableIfExists(db.Models...).Error; err != nil {
		return err
	}
	return db.AutoMigrate()
}

func (db *DB) AutoMigrate() error {
	if err := db.detectModels(); err != nil {
		return err
	}
	return db.Conn.AutoMigrate(db.Models...).Error
}

// Close user service db connection
func (db *DB) Close() error {
	return db.Conn.Close()
}
