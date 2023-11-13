package config

import (
	"flag"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

var (
	godotenvLoad     = godotenv.Load
	envconfigProcess = envconfig.Process
)

type GoDotEnv interface {
	Load(filenames ...string) (err error)
}

type EnvConfig interface {
	Process(prefix string, spec interface{}) error
}

type Config struct {
	Port     int     `envconfig:"PORT"`
	Cass     string  `envconfig:"CASSANDRA"`
	CassUser string  `envconfig:"CASSANDRA_USER"`
	CassPass string  `envconfig:"CASSANDRA_PASS"`
	Eps      float64 `envconfig:"EPS"`
}

func Default() Config {
	return Config{
		Port: 50051,
		Cass: "127.0.0.1:9042",
		Eps: 1,
	}
}

func Init() (*Config, error) {
	err := godotenvLoad()
	if err != nil {
		return nil, errors.Wrap(err, "reading .env file")
	}

	config := Default()

	err = envconfigProcess("", &config)
	if err != nil {
		return nil, errors.Wrap(err, "processing env vars")
	}

	flag.IntVar(&config.Port, "port", config.Port, "The server port")
	flag.StringVar(&config.Cass, "cass", config.Cass, "Cassandra hosts")
	flag.StringVar(&config.CassUser, "cass-user", config.Cass, "Cassandra username")
	flag.StringVar(&config.CassPass, "cass-pass", config.Cass, "Cassandra password")
	flag.Float64Var(&config.Eps, "eps", config.Eps, "Distance dimension")

	return &config, nil
}

func (c *Config) Cassandra() []string {
	return strings.Split(c.Cass, ",")
}
