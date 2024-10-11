package sqlite

import (
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DbSQLite struct {
	DbType string
	Db     *gorm.DB
}

func (db *DbSQLite) Type() string {
	return db.DbType
}

func (db *DbSQLite) Conn() *gorm.DB {
	return db.Db
}

type Opts struct {
	FileName string
}

func (o *Opts) Config() *gorm.Config {
	return &gorm.Config{}
}

func (o *Opts) Dsn() string {
	return o.FileName
}

func NewSqlite(opts *Opts) *DbSQLite {
	db, err := gorm.Open(sqlite.Open(opts.Dsn()), opts.Config())
	if err != nil {
		panic(err)
	}

	return &DbSQLite{
		DbType: "sqlite",
		Db:     db,
	}
}
