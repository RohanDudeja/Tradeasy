package config

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/kelseyhightower/envconfig"
)

type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user_name"`
	DBName   string `yaml:"db_name"`
	Password string `yaml:"password"`
}

func DbURL(config Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Database.UserName,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
	)
}
