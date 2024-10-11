package database

import (
	"gorm.io/gorm"
	"main/database/sqlite"
)

// TODO: Implement Mysql database

type dbOpts interface {
	Config() *gorm.Config
	Dsn() string
}

type DB interface {
	Type() string
	Conn() *gorm.DB
}

type DbInitError struct {
	message string
}

func (err *DbInitError) Error() string {
	return err.message
}

func Factory(dbType string, opts dbOpts) (DB, error) {
	if dbType == "sqlite" {
		sqliteOpts, ok := opts.(*sqlite.Opts)
		if !ok {
			return nil, &DbInitError{message: "Invalid options type for SQLite"}
		}
		dbSqlite := sqlite.NewSqlite(sqliteOpts)
		return dbSqlite, nil
	}
	return nil, &DbInitError{message: "Database type not supported"}
}

func Migrate(db DB, model interface{}) error {
	err := db.Conn().AutoMigrate(model)
	return err
}
