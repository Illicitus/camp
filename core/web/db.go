package web

import (
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
)

type Model interface {
	IsGormModel()
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
	db.LogMode(!mode)
	switch !mode == true {
	case true:
		fmt.Printf("Log mode: is on... \n")
	default:
		fmt.Printf("Log mode: is off... \n")
	}
	return &DB{Conn: db}, nil
}

func (db *DB) detectModels() error {
	if len(db.Models) == 0 {
		return errors.New("models list can't be empty. Please connect db to app first")
	} else {
		for _, m := range db.Models {
			_, ok := m.(Model)
			if !ok {
				return errors.New(fmt.Sprintf("%v is not gorm model.", m))
			}
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
