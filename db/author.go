package db

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	
	"github.com/jackc/pgx/v5"
	"github.com/greyingraven/mamgo/pgconn"
	"github.com/greyingraven/mamgo/cfg"
)

type Author struct {
	Id int
	Name string
	Link string
	Aka []string
}

func InsertAuthor(author *Author) (int, error) {
	conf := cfg.GetConfig()
	query := `INSERT INTO public.author(name, link, aka) VALUES (@name, @link, @aka) RETURNING id`
	args := pgx.NamedArgs{
		"name": author.Name,
		"link": author.Link,
		"aka": author.Aka,
	}
	newId := pgconn.Insert(query, args)
	fmt.Printf("%v\n", newId)
	
	// Create author folders in image to store profile.avif

	path := filepath.Join(conf.Path.Images, strconv.Itoa(newId))
	fmt.Printf("%v\n", path)	
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return 0, err 
	}
	return newId, nil 
}

func UpdateAuthor(author *Author) (err error) {
	query := `UPDATE public.author SET name = @name, link = @link, aka = @aka WHERE id = @id`
	args := pgx.NamedArgs{
		"id": author.Id,
		"name": author.Name,
		"link": author.Link,
		"aka": author.Aka,
	}
	err = pgconn.Update(query, args)
	if err != nil {
		return err
	}
	return nil
}

func GetAuthorById(id int) (author Author, err error) {
	query := fmt.Sprintf("SELECT * FROM public.author WHERE id=%d", id)
	rows, err := pgconn.GetMany(query)
	if err != nil {
		return author, err
	}
	authors, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Author])
	if err != nil {
		return author, err
	}
	return authors[0], nil
}
