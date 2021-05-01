package http

import (
	"fmt"
	"net/http"
)

func NewServer(
	port int,
	bookmarksHandler *BookmarksHandler,
) *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/", bookmarksHandler)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler:           mux,
	}

	return server
}
