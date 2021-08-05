package config

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

var DB  *gorm.DB


// Config represents configuration
type Config struct {
	Database Database `yaml:"database"`
}
type Database struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	UserName string `yaml:"user_name"`
	DBName   string `yaml:"db_name"`
	Password string `yaml:"password"`
}

//readFile for reading config.yml file
func readFile(cfg *Config) {
	f, err := os.Open("../../config/config.yml")
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

// InitialiseDB ...assign connection to global *gorm.DB variable DB
func InitialiseDB() error {
	dbString := DbURL(BuildConfig())
	var err error
	DB, err = gorm.Open("mysql", dbString)
	if err != nil {
		return err
	}
	return nil
}
