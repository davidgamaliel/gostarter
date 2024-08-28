package database

import (
	"fmt"

	dbtemplate "github.com/bitzero/gostarter/internal/database/template"
)

type Driver string

const (
	Postgres Driver = "postgres"
	MySQL    Driver = "mysql"
	Monggo   Driver = "monggo"
)

type DBDriverTemplater interface {
	Service() []byte
	Env() []byte
}

type DBDriver struct {
	Driver      Driver
	PackageName []string
	Template    DBDriverTemplater
}

func (d Driver) String() string {
	return string(d)
}

var AllowedDBDrivers = []string{string(MySQL), string(Monggo), string(Postgres)}

var (
	mysqlDriver    = []string{"github.com/go-sql-driver/mysql"}
	postgresDriver = []string{"github.com/lib/pq"}
	sqliteDriver   = []string{"github.com/mattn/go-sqlite3"}
	redisDriver    = []string{"github.com/redis/go-redis/v9"}
	mongoDriver    = []string{"go.mongodb.org/mongo-driver"}
)

func (d *Driver) Set(value string) error {
	for _, driver := range AllowedDBDrivers {
		if driver == value {
			*d = Driver(value)
			return nil
		}
	}

	return fmt.Errorf("invalid database driver: %s", value)
}

func GetDBDriverMap() map[string]DBDriver {
	return map[string]DBDriver{
		"mysql": {
			Driver:      MySQL,
			PackageName: mysqlDriver,
			Template:    dbtemplate.MysqlTemplate{},
		},
		"postgres": {
			Driver:      Postgres,
			PackageName: postgresDriver,
			Template:    dbtemplate.PostgresTemplate{},
		},
		"monggo": {
			Driver:      Monggo,
			PackageName: mongoDriver,
			Template:    dbtemplate.MongoTemplate{},
		},
	}
}
