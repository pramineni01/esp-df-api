package esputils

import (
	"os"

	"github.com/go-sql-driver/mysql"
)

func GetDBConfig() *mysql.Config {
	db_config := mysql.NewConfig()
	db_config.User = os.Getenv("DB_USER")
	db_config.Passwd = os.Getenv("DB_PASSWORD")
	db_config.DBName = os.Getenv("DB_NAME")
	db_config.Addr = os.Getenv("DB_HOST") + ":3306"
	db_config.Net = "tcp"
	db_config.ParseTime = true

	return db_config
}
