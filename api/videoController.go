package api

import (
	"fmt"
	"encoding/json"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/greyingraven/mamgo/db"
	"github.com/greyingraven/mamgo/cfg"
)

var (
    VideoRe       = regexp.MustCompile(`^/video/*$`)
    VideosRe       = regexp.MustCompile(`^/videos/*$`)
    VideoReWithID = regexp.MustCompile(`^/video/([0-9]+)$`)
	VideoReWithVID = regexp.MustCompile(`^/video/([a-zA-Z0-9]+)$`)
)

func (v *videoHandler) ListVideos(w http.ResponseWriter, r *http.Request) {
	// Find all videos in db
	videos, err := db.FindVideos()
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Parse videos to json and return to client
	jsonBytes, err := json.Marshal(videos)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (v *videoHandler) GetVideo(w http.ResponseWriter, r *http.Request) {
	// Extract resource ID using regex
	matches := VideoReWithID.FindStringSubmatch(r.URL.Path)
	// matches should be length >=2
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// Retrieve video from db by ID
	id, _ := strconv.Atoi(matches[1])
	video, err := db.GetVideoById(id)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(video)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (v *videoHandler) ServeVideo(w http.ResponseWriter, r *http.Request) {
	// TODO: Check stream video test and replicate serving a video found by v_id
	// Path for video should be ROOT/videos/{author_id}/{v_id}.mp4
	// Extract resource ID using regex
	matches := VideoReWithVID.FindStringSubmatch(r.URL.Path)
	// matches should be length >=2
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// Retrieve author_id from db by v_id
	v_id := matches[1]
	vid, err := db.GetVideoByVId(v_id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Video not found in db: %v", err), http.StatusInternalServerError)
		return 
	}
	a_id := vid.Author_id
	videoPath := fmt.Sprintf("%v/%v/%v.mp4", cfg.Cfg.Path.Videos, strconv.Itoa(a_id), v_id)
	f, err := os.Open(videoPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Video not found: %v", err), http.StatusInternalServerError)
		return 
	}
	fmt.Printf("Video file: %v\n", f)
	fi, err := f.Stat()
	if err != nil {
		http.Error(w, fmt.Sprintf("Can't get video info: %v", err), http.StatusInternalServerError)
		return 
	}
	totalSize := fi.Size()
	
	//	f, totalSize := stream.FindVideo(strconv.Itoa(author_id), strconv.Itoa(v_id))
	
	rangeHeader := r.Header.Get("Range")
	var start, buffSize int64
	buffSize = 1024*1024 - 1  
	
	if rangeHeader == "" {
		// Default to first 1MB
		start = 0
	} else {
		// Parse the Range header: "bytes=start-end"
		rangeParts := strings.TrimPrefix(rangeHeader, "bytes=")
		rangeValues := strings.Split(rangeParts, "-")
		var err error
		
		// Get start byte
		startValue, err := strconv.Atoi(rangeValues[0])
		if err != nil {
			http.Error(w, "Invalid start byte", http.StatusBadRequest)
			return
		}
		start = int64(startValue)
		
		// Get end byte or set to default
		if len(rangeValues) > 1 && rangeValues[1] != "" {
			end, err := strconv.Atoi(rangeValues[1])
			if err != nil {
				http.Error(w, "Invalid end byte", http.StatusBadRequest)
				return
			}
			buffSize = int64(end) - start
		} 
	}
	
	// Ensure end is within the total video size
	if start + buffSize >= totalSize {
		buffSize = int64(totalSize) - start
	}
	
	// Fetch the video data
	_, err = f.Seek(start, 0)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error seek in video file: %v", err), http.StatusInternalServerError)
		return
	}
	buf := make([]byte, buffSize)
	videoData, err := f.Read(buf)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error retrieving video: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers and serve the video chunk
	w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, start+int64(videoData), totalSize))
	w.Header().Set("Accept-Ranges", "bytes")
	w.Header().Set("Content-Length", fmt.Sprintf("%d", videoData))
	w.Header().Set("Content-Type", "video/mp4")
	w.WriteHeader(http.StatusPartialContent)
	_, err = w.Write(buf[:videoData])
	if err != nil {
		http.Error(w, "Error streaming video", http.StatusInternalServerError)
	}
}

func (v *videoHandler) CreateVideo(w http.ResponseWriter, r *http.Request) {
	var newVideo *db.Video
	err := json.NewDecoder(r.Body).Decode(&newVideo)
	if err != nil {
		http.Error(w, "Error reading video info from body", http.StatusInternalServerError)
		return
	}
	id := db.InsertVideo(newVideo)
	response := fmt.Sprintf("Created new video with id: %d\n", id)
	fmt.Println(response)
	
	w.WriteHeader(http.StatusOK)
	//	w.Write(response)
}

func (v *videoHandler) UpdateVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered UpdateVideo")
}

func (v *videoHandler) DeleteVideo(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered DeleteVideo")
}

func (v *videoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch{
	case r.Method == http.MethodPost && VideoRe.MatchString(r.URL.Path):
		v.CreateVideo(w, r)
		return
	case r.Method == http.MethodGet && VideosRe.MatchString(r.URL.Path):
		v.ListVideos(w, r)
		return
	case r.Method == http.MethodGet && VideoReWithID.MatchString(r.URL.Path):
		v.GetVideo(w, r)
		return
	case r.Method == http.MethodGet && VideoReWithVID.MatchString(r.URL.Path):
		v.ServeVideo(w, r)
		return
	case r.Method == http.MethodPut && VideoReWithID.MatchString(r.URL.Path):
		v.UpdateVideo(w, r)
		return
	case r.Method == http.MethodDelete && VideoReWithID.MatchString(r.URL.Path):
		v.DeleteVideo(w, r)
		return
	default:
		return
	}
}

type videoHandler struct{}
