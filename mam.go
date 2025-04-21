package main

import (
	"fmt"
	"os"
	
	"github.com/greyingraven/mamgo/cfg"
	"github.com/greyingraven/mamgo/pgconn"
	"github.com/greyingraven/mamgo/db"
)

func main() {
	fmt.Println("Loading configuration file")
	conf, err := cfg.LoadConfig("cfg/mamgo.ini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config file: %v\n", err)
	}
	fmt.Fprintf(os.Stdout, "Main ToDo: %s\n", conf.Todo)
	pgconn.StartConnection()
	fmt.Println(db.GetVideo())
}
