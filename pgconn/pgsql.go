package pgconn

import (
	"context"
	"fmt"
	"os"
	
	"github.com/greyingraven/mamgo/cfg"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pgconn struct {
	db *pgxpool.Pool
}

var Pool *Pgconn

func Insert(query string, args pgx.NamedArgs) int {
	var id int
	fmt.Printf("Query: %v\n", query)
	fmt.Println(args)
	Pool.db.QueryRow(context.Background(), query, args).Scan(&id)
	fmt.Printf("%v\n", id)
	return id
}

func Update(query string, args pgx.NamedArgs) error {
	return Exec(query, args, "Unable to insert row: %w")
}

func Exec(query string, args pgx.NamedArgs, error string) error {
	_, err := Pool.db.Exec(context.Background(), query, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, error, err)
		return fmt.Errorf("Unable to update row: %w", err)
	}
	return nil
}

func GetOne(query string) (row pgx.Row) {
	row = Pool.db.QueryRow(context.Background(), query)
	return row
}

func GetMany(query string) (rows pgx.Rows, err error) {
	rows, err = Pool.db.Query(context.Background(), query)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Query failed: %v\n", err)
		return nil, err
	}	
	return rows, nil	
}

func StartConnection() {
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
	pool := &Pgconn{
		db: db,
	}

	Pool = pool
	// defer pgconn.db.Close()
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
