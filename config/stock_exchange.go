package config

type StockExchange struct {
	Authentication Authentication `yaml:"authentication"`
	Url            string         `yaml:"url"`
}

type Authentication struct {
	UserName string `yaml:"user_name"`
	Password string `yaml:"password"`
}
