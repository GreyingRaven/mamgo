package api

import (
	"fmt"
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/greyingraven/mamgo/db"
)

var (
    AuthorRe       = regexp.MustCompile(`^/author/*$`)
    AuthorsRe       = regexp.MustCompile(`^/authors/*$`)
    AuthorReWithID = regexp.MustCompile(`^/author/([0-9]+)$`)
	AuthorReWithVID = regexp.MustCompile(`^/author/([a-zA-Z0-9]+)$`)
)

func (v *authorHandler) ListAuthors(w http.ResponseWriter, r *http.Request) {
	// Find all authors in db
	authors, err := db.FindAuthors()
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	// Parse authors to json and return to client
	jsonBytes, err := json.Marshal(authors)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (v *authorHandler) GetAuthor(w http.ResponseWriter, r *http.Request) {
	// Extract resource ID using regex
	matches := AuthorReWithID.FindStringSubmatch(r.URL.Path)
	// matches should be length >=2
	if len(matches) < 2 {
		InternalServerErrorHandler(w, r)
		return
	}

	// Retrieve author from db by ID
	id, _ := strconv.Atoi(matches[1])
	author, err := db.GetAuthorById(id)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	jsonBytes, err := json.Marshal(author)
	if err != nil {
		InternalServerErrorHandler(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonBytes)
}

func (v *authorHandler) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var newAuthor *db.Author
	err := json.NewDecoder(r.Body).Decode(&newAuthor)
	if err != nil {
		http.Error(w, "Error reading author info from body", http.StatusInternalServerError)
		return
	}
	id, err := db.InsertAuthor(newAuthor)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error creating new author: %v", err), http.StatusInternalServerError)
		return
	}
	response := fmt.Sprintf("Created new author with id: %d\n", id)
	fmt.Println(response)
	
	w.WriteHeader(http.StatusOK)
}

func (v *authorHandler) UpdateAuthor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered UpdateAuthor")
}

func (v *authorHandler) DeleteAuthor(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Entered DeleteAuthor")
}

func (v *authorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch{
	case r.Method == http.MethodPost && AuthorRe.MatchString(r.URL.Path):
		v.CreateAuthor(w, r)
		return
	case r.Method == http.MethodGet && AuthorsRe.MatchString(r.URL.Path):
		v.ListAuthors(w, r)
		return
	case r.Method == http.MethodGet && AuthorReWithID.MatchString(r.URL.Path):
		v.GetAuthor(w, r)
		return
	case r.Method == http.MethodPut && AuthorReWithID.MatchString(r.URL.Path):
		v.UpdateAuthor(w, r)
		return
	case r.Method == http.MethodDelete && AuthorReWithID.MatchString(r.URL.Path):
		v.DeleteAuthor(w, r)
		return
	default:
		return
	}
}

type authorHandler struct{}
