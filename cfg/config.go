package cfg

import (
	"fmt"
	"os"

	"github.com/gookit/ini/v2"
)

var Cfg *Config

type Config struct {
	Todo string
	Db Db
}

type Db struct {
	Host     string
	Port     int
	Database string
	User     string
	Password string
}

func GetConfig() (*Config) {
	return Cfg
}

func GetDbConfig() (Db) {
	return Cfg.Db
}

func LoadConfig(configFile string) (*Config, error) {
	err := ini.LoadExists(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config file: %v\n", err)
		return nil, err
	}
	db := &Db{
		Host: ini.String("db.host"),
		Port: ini.Int("db.port"),
		Database: ini.String("db.name"),
		User: ini.String("db.user"),
		Password: ini.String("db.password"),
	}
	config := &Config{
		Todo: ini.String("main.todo"),
		Db: *db,
	}
	Cfg = config
	return config, nil
}
