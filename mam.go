package main

import (
	"fmt"
	"os"
	
	"github.com/greyingraven/mamgo/cfg"
	"github.com/greyingraven/mamgo/pgconn"
	"github.com/greyingraven/mamgo/db"
)

func InsertTest() int {
	author := &db.Author{
		Name: "Tsukimi69",
		Link: "https://www.iwara.tv/profile/Tsukimi69",
		Aka: []string{"empty"},
	}
	authorId, err := db.InsertAuthor(author)
	fmt.Println(authorId)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Insert author failed: %v\n", err)
		return 0
	}
	video := &db.Video{
		Title: "ulm2 - Rabbit Hole - DECO27 [b3sIqbBCEIhCbp][c$ Kangxi][t$ Sex KawaiiStrike].mp4",
		Author_id: authorId,
		Src: "https://www.iwara.tv/video/b3sIqbBCEIhCbp",
		V_type: "mmd",
		V_id: "b3sIqbBCEIhCbp",
	}
	videoId := db.InsertVideo(video)
	fmt.Println(videoId)
	return videoId
}

func UpdateTest() {
	video := &db.Video{
		Id: 16,
		Title: "ulm2 - Rabbit Hole - DECO27 [b3sIqbBCEIhCbp][c$ Kangxi][t$ Sex KawaiiStrike].mp4",
		Author_id: 2,
		Src: "https://www.iwara.tv/video/b3sIqbBCEIhCbp",
		V_type: "mmd",
		V_id: "b3sIqbBCEIhCbp",
	}
	videoId := db.UpdateVideo(video)
	fmt.Println(videoId)
}

func main() {
	videoId := 16
	fmt.Println("Loading configuration file")
	conf, err := cfg.LoadConfig("cfg/mamgo.ini")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to load config file: %v\n", err)
	}
	fmt.Fprintf(os.Stdout, "Main ToDo: %s\n", conf.Todo)
	pgconn.StartConnection()
	//	videoId = InsertTest()
	UpdateTest()
	if videoId == 0 {
		fmt.Fprintf(os.Stderr, "Insert test failed")
		return
	}
	fmt.Println(db.GetVideoById(videoId))
}

