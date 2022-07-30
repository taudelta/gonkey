package fixtures

import (
	"database/sql"
	"strings"

	_ "github.com/lib/pq"

	"github.com/lamoda/gonkey/fixtures/aerospike"
	"github.com/lamoda/gonkey/fixtures/mysql"
	"github.com/lamoda/gonkey/fixtures/postgres"
	aerospikeClient "github.com/lamoda/gonkey/storage/aerospike"
)

type DbType int

const (
	Postgres DbType = iota
	Mysql
	Aerospike
	CustomLoader
)

const (
	PostgresParam     = "postgres"
	MysqlParam        = "mysql"
	AerospikeParam    = "aerospike"
	CustomLoaderParam = "custom"
)

type Config struct {
	DB        *sql.DB
	Aerospike *aerospikeClient.Client
	DbType    DbType
	Location  string
	Debug     bool
	Loaders   []Loader
}

type Loader interface {
	Load(names []string) error
}

func NewLoader(cfg *Config) Loader {

	var loader Loader

	location := strings.TrimRight(cfg.Location, "/")

	switch cfg.DbType {
	case Postgres:
		loader = postgres.New(
			cfg.DB,
			location,
			cfg.Debug,
		)
	case Mysql:
		loader = mysql.New(
			cfg.DB,
			location,
			cfg.Debug,
		)
	case Aerospike:
		loader = aerospike.New(
			cfg.Aerospike,
			location,
			cfg.Debug,
		)
	default:
		if len(cfg.Loaders) > 0 {
			return cfg.Loaders[0]
		}
		panic("unknown db type")
	}

	return loader
}

func FetchDbType(dbType string) DbType {
	switch dbType {
	case PostgresParam:
		return Postgres
	case MysqlParam:
		return Mysql
	case AerospikeParam:
		return Aerospike
	case CustomLoaderParam:
		return CustomLoader
	default:
		panic("unknown db type param")
	}
}
