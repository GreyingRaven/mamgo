package db

import (
	"fmt"
	
	"github.com/jackc/pgx/v5"
	"github.com/greyingraven/mamgo/pgconn"
)

type Video struct {
	Id        int
	Title     string
	Author_id int
	Src       string
	V_type    string
	V_id      string
}

func InsertVideo(video *Video) int {
	query := `INSERT INTO public.video(title, author_id, src, v_type, v_id) VALUES (@title, @author_id, @src, @type, @v_id) RETURNING id`
	args := pgx.NamedArgs{
		"title": video.Title,
		"author_id": video.Author_id,
		"src": video.Src,
		"type":	video.V_type,
		"v_id": video.V_id,
	}
	newId := pgconn.Insert(query, args)
	return newId
}

func UpdateVideo(video *Video) (err error) {
	query := `UPDATE public.video SET title = @title, author_id = @author_id, src = @src, v_type = @type, v_id = @v_id WHERE id = @id`
	args := pgx.NamedArgs{
		"id": video.Id,
		"title": video.Title,
		"author_id": video.Author_id,
		"src": video.Src,
		"type":	video.V_type,
		"v_id": video.V_id,
	}
	err = pgconn.Update(query, args)
	if err != nil {
		return err
	}
	return nil
}

func GetVideoById(id int) (video Video, err error) {
	query := fmt.Sprintf("SELECT * FROM public.video WHERE id=%d", id)
	rows, err := pgconn.GetMany(query)
	if err != nil {
		return video, err
	}
	videos, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Video])
	if err != nil {
		return video, err
	}
	return videos[0], nil	
}
