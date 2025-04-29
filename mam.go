package main

import (
	"fmt"
	"os"
	
	"github.com/greyingraven/mamgo/cfg"
	"github.com/greyingraven/mamgo/pgconn"
	"github.com/greyingraven/mamgo/api"
)


func main() {
	fmt.Println("Loading configuration file")
	conf, err := cfg.LoadConfig("cfg/mamgo.ini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config file: %v\n", err)
	}
	fmt.Fprintf(os.Stdout, "Main ToDo: %s\n", conf.Todo)
	// Start database connection	
	fmt.Println("Starting database connection")
	// TODO: Change to start and close connection when called
	pgconn.StartConnection()
	// Start http server
	fmt.Printf("Starting HttpServer. Listening on port: %v\n", conf.Port)
	api.MamGoHandler()	
}

