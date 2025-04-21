package pgconn

import (
	"context"
	"fmt"
	"os"
	
	"github.com/greyingraven/mamgo/cfg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var db pgxpool.Pool


func Insert(query string, args pgx.NamedArgs) {
	_, err := db.Exec(context.Background(), query, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert row: %w", err)
		return fmt.Errorf("Unable to insert row: %w", err)
	}
	return nil
}

func GetOne(query string) (row pgx.Row, err error) {
	row, err := db.QueryRow(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		return nil, err
	}
	
	return row, nil
}

func GetMany(query string) (rows pgx.Rows, err error) {
	rows, err := db.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return nil, err
	}
	
	return rows, nil	
}

func Test() {
	conf := cfg.GetDbConfig()
	poolCfg, err := pgxpool.ParseConfig(fmt.Sprintf("user=%s password=%s host=%s port=%d dbname=%s", conf.User, conf.Password, conf.Host, conf.Port, conf.Database))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to parse pool configuration: %v\n", err)
		os.Exit(1)
	}
	db, err := pgxpool.NewWithConfig(context.Background(), poolCfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	var name string
	err = db.QueryRow(context.Background(), "select username from users where id=$1", 3).Scan(&name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	
	fmt.Println(name)

	//video := Video{
	//	Path: "/test/path",
	//	Author: 2,
	//	Src: "https://iwara.../iwara_id",
	//	Type: "mmd",
	//}
	//fmt.Sprintf("Inserting video: %s\n", video)
	//InsertVideo(db, context.Background(), video)
	
}
