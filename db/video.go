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

type Video struct {
	Id        int
	Title     string
	Author_id int
	Src       string
	V_type    string
	V_id      string
}

type VideoSearch struct {
	Title     string
	Author_id int
	V_type    string
	V_id      string
}

type VideoStream struct {
	Author_id int
}

func InsertVideo(video *Video) (int, error) {
	conf := cfg.GetConfig()
	query := `INSERT INTO public.video(title, author_id, src, v_type, v_id) VALUES (@title, @author_id, @src, @type, @v_id) RETURNING id`
	args := pgx.NamedArgs{
		"title": video.Title,
		"author_id": video.Author_id,
		"src": video.Src,
		"type":	video.V_type,
		"v_id": video.V_id,
	}
	// Create author folders in videos if it doesn't exist
	path := filepath.Join(conf.Path.Videos, strconv.Itoa(video.Author_id))
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		return 0, err 
	}
	// Move video to author folder
	videoSrc := fmt.Sprintf("%v/%v", conf.Path.Videos, video.Src)
	videoPath := fmt.Sprintf("%v/%v/%v.mp4", conf.Path.Videos, video.Author_id, video.V_id)
	err = os.Rename(videoSrc, videoPath)
	if err != nil {
		return 0, err
	}
	newId := pgconn.Insert(query, args)
	return newId, nil
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

func FindVideos() (videos []VideoSearch, err error) {
	query := fmt.Sprintf("SELECT title, author_id, v_type, v_id FROM public.video")
	rows, err := pgconn.GetMany(query)
	if err != nil {
		return videos, err
	}
	videos, err = pgx.CollectRows(rows, pgx.RowToStructByName[VideoSearch])
	if err != nil {
		return videos, err
	}
	return videos, nil	
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

func GetVideoByVId(v_id string) (video VideoStream, err error) {
	query := fmt.Sprintf("SELECT author_id FROM public.video WHERE v_id='%s'", v_id)
	rows, err := pgconn.GetMany(query)
	if err != nil {
		return video, err
	}
	videos, err := pgx.CollectRows(rows, pgx.RowToStructByPos[VideoStream])
	if err != nil {
		return video, err
	}
	return videos[0], nil	
}
