package domain

import (
	"io"
	"io/ioutil"
)

type BookmarkContentCollector interface {
	CollectText(resourceBody io.ReadCloser) (io.ReadCloser, error)
}

type ResourceBodyProvider interface {
	GetResourceBody(url string) (io.ReadCloser, error)
}

func NewBookmarkFactory(resourceBodyProvider ResourceBodyProvider, bookmarkContentCollector BookmarkContentCollector) *BookmarkFactory {
	return &BookmarkFactory{resourceBodyProvider: resourceBodyProvider, bookmarkContentCollector: bookmarkContentCollector}
}

type BookmarkFactory struct {
	resourceBodyProvider ResourceBodyProvider
	bookmarkContentCollector BookmarkContentCollector
}

func (bf *BookmarkFactory) CreateBookmark(url string, title *string) (*Bookmark, error) {
	resourceBody, err := bf.resourceBodyProvider.GetResourceBody(url)
	if err != nil {
		return nil, err
	}
	defer resourceBody.Close()

	resourceText, err := bf.bookmarkContentCollector.CollectText(resourceBody)
	if err != nil {
		return nil, err
	}

	all, err := ioutil.ReadAll(resourceText)
	if err != nil {
		return nil, err
	}
	defer resourceText.Close()

	finalTitle := url
	if title != nil {
		finalTitle = *title
	}

	return &Bookmark{
		Title:   finalTitle,
		Content: string(all),
	}, nil
}

