package api

import (
	"fmt"
	"net/http"

	"github.com/rs/cors"
	"github.com/greyingraven/mamgo/cfg"
)

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    w.Write([]byte("Home Placeholder"))
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 Internal Server Error"))
}

func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func MamGoHandler() {
	// Load configuration
	conf := cfg.GetConfig()
	// Create a request multiplexer
	// to dispatch request to the matching handlers
	mux := http.NewServeMux()
	// Register routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/v1/video", &videoHandler{})
	mux.Handle("/v1/videos", &videoHandler{})
	mux.Handle("/v1/video/", &videoHandler{})
	mux.Handle("/v1/author", &authorHandler{})
	mux.Handle("/v1/authors", &authorHandler{})
	mux.Handle("/v1/author/", &authorHandler{})
	// Run the server
	port := conf.Port
	servePort := fmt.Sprintf(":%v", port)
	handler := cors.Default().Handler(mux)
	http.ListenAndServe(servePort, handler)
}

type homeHandler struct{}
