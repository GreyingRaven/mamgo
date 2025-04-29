package api

import (
	"fmt"
	"net/http"
	
	"github.com/greyingraven/mamgo/cfg"
)

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Home page placeholder"))
}

func InternalServerErrorHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Write([]byte("500 Internal Server Error"))
}

func MamGoHandler() {
	// Load configuration
	conf := cfg.GetConfig()
	// Create a request multiplexer
	// to dispatch request to the matching handlers
	mux := http.NewServeMux()
	// Register routes and handlers
	mux.Handle("/", &homeHandler{})
	mux.Handle("/video", &videoHandler{})
	mux.Handle("/videos", &videoHandler{})
	mux.Handle("/video/", &videoHandler{})
	mux.Handle("/author", &authorHandler{})
	mux.Handle("/authors", &authorHandler{})
	mux.Handle("/author/", &authorHandler{})
	// Run the server
	port := conf.Port
	servePort := fmt.Sprintf(":%v", port)
	http.ListenAndServe(servePort, mux)
}

type homeHandler struct{}
