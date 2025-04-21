package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/greyingraven/mamgo/pgsql"
)

type Video struct {
	ID int
	Path string
	Author int
	Src string
	Type string
}

func InsertVideo(video Video) (err error) {
	query := `INSERT INTO public.video(path, author_id, src, type) VALUES (@path, @author_id, @src, @type)`
	args := pgx.NamedArgs{
		"path": video.Path,
		//"author_id": video.Author,
		"src": video.Src,
		"type":	video.Type,
	}
	_, err := pgsql.Insert(query, args)
	if err != nil {
		return err
	}
	return nil
}

func GetVideo() (path string, err error) {
	var path string
	query := `SELECT path FROM public.video WHERE id=3`
	_, err := pgsql.GetOne(query).Scan(&path)
	if err != nil {
		return nil, err
	}
	return path, nil
	
}
