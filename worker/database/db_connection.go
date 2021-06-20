package database

import (
	"database/sql"
	"fmt"

	"github.com/caarlos0/env/v6"

	// blank import for mysql driver

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DbConn *gorm.DB
)

type Config struct {
	DbHost     string `env:"DB_HOST"`
	DbName     string `env:"MYSQL_DATABASE"`
	DbUser     string `env:"MYSQL_USER"`
	DbPassword string `env:"MYSQL_PASSWORD"`
	DbConn     *sql.DB
}

func Connect() error {
	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("%+v", err)
	}
	dsn := cfg.DbUser + ":" + cfg.DbPassword + "@" + cfg.DbHost + "/" + cfg.
		DbName + "?parseTime=true&charset=utf8"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	if err != nil {
		panic("failed to connect database")
	}

	DbConn = db

	return nil
}
