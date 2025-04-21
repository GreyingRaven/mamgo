package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/greyingraven/mamgo/pgconn"
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
	err = pgconn.Insert(query, args)
	if err != nil {
		return err
	}
	return nil
}

func GetVideo() (path string, err error) {
	query := `SELECT path FROM public.video WHERE id=3`
	err = pgconn.GetOne(query).Scan(&path)
	if err != nil {
		return ``, err
	}
	return path, nil
	
}
