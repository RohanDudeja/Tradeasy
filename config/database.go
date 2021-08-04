package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"os"
)

// Config represents configuration
type Config struct {
	Database Database `yaml:"database",envconfig:"Database"`
}
type Database struct {
	Host     string `yaml:"host", envconfig:"DB_HOST"`
	Port     int    `yaml:"port", envconfig:"DB_PORT"`
	UserName string `yaml:"user_name", envconfig:"DB_USERNAME"`
	DBName   string `yaml:"db_name", envconfig:"DB_NAME"`
	Password string `yaml:"password", envconfig:"DB_PASSWORD"`
}

//readEnv for reading environment variable
func readEnv(cfg *Config) {
	err := envconfig.Process("", cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
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
	readEnv(&cfg)
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
