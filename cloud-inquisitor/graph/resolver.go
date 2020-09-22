package graph

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

import (
	"errors"
	"strings"
	"time"

	"github.com/RiotGames/cloud-inquisitor/cloud-inquisitor/graph/model"
	"github.com/RiotGames/cloud-inquisitor/cloud-inquisitor/secrets/vault"
	"github.com/RiotGames/cloud-inquisitor/cloud-inquisitor/settings"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

/*func init() {
	//do all database related setup
	db, err := NewDBConnection()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if !db.HasTable(&model.Account{}) {
		db.CreateTable(&model.Account{})
	}
	if !db.HasTable(&model.Zone{}) {
		db.CreateTable(&model.Zone{})
	}
	if !db.HasTable(&model.Record{}) {
		db.CreateTable(&model.Record{})
	}
	if !db.HasTable(&model.Value{}) {
		db.CreateTable(&model.Value{})
	}
}*/

func NewDBConnection() (*gorm.DB, error) {
	if settings.IsSet("mysql") {
		var connectionString string
		var err error = nil
		switch settings.GetString("mysql.connection_string.type") {
		case "vault":
			connectionString, err = vault.GetString(settings.GetString("mysql.connection_string.value"), "value")
			if err != nil {
				return nil, err
			}
		case "string":
			connectionString = settings.GetString("mysql.connection_string.value")
		default:
			return nil, errors.New("unknown datastore/graph configuration type")
		}

		if connectionString == "" {
			return nil, errors.New("datastore/graph connection is empty")
		}

		if !strings.HasSuffix(connectionString, "?parseTime=true") {
			connectionString = connectionString + "?parseTime=true"
		}

		return gorm.Open("mysql", connectionString)
	}

	return nil, errors.New("no mysql connection config provided")
}

func CreateTables() error {
	db, err := NewDBConnection()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return err
	}
	defer db.Close()

	if !db.HasTable(&model.Account{}) {
		db.CreateTable(&model.Account{})
	}
	if !db.HasTable(&model.Zone{}) {
		db.CreateTable(&model.Zone{})
	}
	if !db.HasTable(&model.Record{}) {
		db.CreateTable(&model.Record{})
	}
	if !db.HasTable(&model.Value{}) {
		db.CreateTable(&model.Value{})
	}

	if !db.HasTable(&model.Distribution{}) {
		db.CreateTable(&model.Distribution{})
	}

	if !db.HasTable(&model.Origin{}) {
		db.CreateTable(&model.Origin{})
	}
	if !db.HasTable(&model.OriginGroup{}) {
		db.CreateTable(&model.OriginGroup{})
	}

	if !db.HasTable(&model.ElasticbeanstalkEnvironment{}) {
		db.CreateTable(&model.ElasticbeanstalkEnvironment{})
	}

	if !db.HasTable(&model.S3{}) {
		db.CreateTable(&model.S3{})
	}

	return nil
}

func DropTables() error {
	db, err := NewDBConnection()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return err
	}
	defer db.Close()

	err = db.DropTableIfExists(
		&model.Account{},
		&model.Zone{},
		&model.Record{},
		&model.Value{},
		&model.Distribution{},
		&model.Origin{},
		&model.OriginGroup{},
		&model.ElasticbeanstalkEnvironment{},
		&model.S3{},
	).Error

	if err != nil {
		return err
	}

	return nil
}

func MigrateTables() error {
	db, err := NewDBConnection()
	if err != nil {
		if db != nil {
			db.Close()
		}
		return err
	}
	defer db.Close()

	err = db.AutoMigrate(
		&model.Account{},
		&model.Zone{},
		&model.Record{},
		&model.Value{},
		&model.Distribution{},
		&model.Origin{},
		&model.OriginGroup{},
		&model.ElasticbeanstalkEnvironment{},
		&model.S3{},
	).Error

	if err != nil {
		return err
	}

	return nil
}

type Resolver struct {
	DB *gorm.DB
}

func NewResolver() (*Resolver, error) {
	db, err := NewDBConnection()
	if err != nil {
		db.Close()
		return nil, err
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(10)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(100)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(time.Hour)

	return &Resolver{DB: db}, nil
}

func (r *Resolver) Close() error {
	return r.DB.Close()
}
