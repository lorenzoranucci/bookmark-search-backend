package http

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/domain"
)

func NewBookmarksHandler(
	bookmarkFactory *domain.BookmarkFactory,
	bookmarkRepository domain.BookmarkRepository,
) *BookmarksHandler {
	return &BookmarksHandler{bookmarkFactory: bookmarkFactory, bookmarkRepository: bookmarkRepository}
}

type BookmarksHandler struct {
	bookmarkFactory *domain.BookmarkFactory
	bookmarkRepository domain.BookmarkRepository
}

type BookmarkPostRequest struct {
	URL string `json:"url"`
	Title *string `json:"title,omitempty"`
}

func (b *BookmarksHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		b.handlePost(writer, request)
		return
	case http.MethodGet:
		b.handleGet(writer, request)
		return
	}
}

func (b *BookmarksHandler) handlePost(writer http.ResponseWriter, request *http.Request) {
	body, err := ioutil.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	bookmarkPostRequest := &BookmarkPostRequest{}
	err = json.Unmarshal(body, bookmarkPostRequest)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	domainBookmark, err := b.bookmarkFactory.CreateBookmark(bookmarkPostRequest.URL, bookmarkPostRequest.Title)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = b.bookmarkRepository.Add(domainBookmark)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	writer.WriteHeader(http.StatusCreated)
}

func (b *BookmarksHandler) handleGet(writer http.ResponseWriter, request *http.Request) {

}


