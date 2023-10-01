package config

type Config struct {
	Port string
}

func New() *Config {
	port := "42069"
	return &Config{
		Port: port,
	}
}
