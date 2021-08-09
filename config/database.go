package config

import (
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

var DB *gorm.DB

// Config represents configuration
type Config struct {
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
	Authentication Authentication `yaml:"authentication"`
}
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user_name"`
	DBName   string `yaml:"db_name"`
	Password string `yaml:"password"`
}

type Server struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type Authentication struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
}
//readFile for reading development.yaml file
func readFile(cfg *Config) {
	f, err := os.Open("./config/development.yaml")
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}

func BuildConfig() *Config {
	var cfg Config
	readFile(&cfg)
	//fmt.Printf("%+v", cfg)
	return &cfg

}
func DbURL(config *Config) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Database.UserName,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DBName,
	)
}
func ServerURL(config *Config) string {
	return fmt.Sprintf(
		"%s:%d",
		config.Server.Host,
		config.Server.Port,
	)
}
func ApiToken(config *Config) (string,string) {
	return config.Authentication.UserName, config.Authentication.Password
}
// InitialiseDB ...assign connection to global *gorm.DB variable DB
func InitialiseDB() error {
	dbString := DbURL(BuildConfig())
	var err error
	DB, err = gorm.Open("mysql", dbString)
	if err != nil {
		return err
	}
	if gin.IsDebugging() {
		DB.LogMode(true)
	}
	return nil
}
