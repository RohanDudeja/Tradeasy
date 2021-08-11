package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

// Config represents configuration
type Config struct {
	Database      Database      `yaml:"database"`
	Server        Server        `yaml:"server"`
	StockExchange StockExchange `yaml:"stock_exchange"`
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
	config = &cfg
	//fmt.Printf("%+v", cfg)
	return &cfg

}

var config *Config

func GetConfig() Config {
	return *config
}
