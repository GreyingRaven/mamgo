package cfg

import (
	"fmt"
	"os"

	"github.com/gookit/ini/v2"
)

var Cfg *Config

type Config struct {
	Todo string
	Root string
	Path Path
	Db   Db
}

type Path struct {
	Videos string
	Music  string
	Images string
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
	ini.WithOptions(ini.ParseEnv, ini.ParseVar)
	err := ini.LoadExists(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config file: %v\n", err)
		return nil, err
	}
	path := &Path{
		Videos: ini.String("path.video"),
		Music: ini.String("path.music"),
		Images: ini.String("path.image"),
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
		Path: *path,
		Db: *db,
	}
	Cfg = config
	return config, nil
}
