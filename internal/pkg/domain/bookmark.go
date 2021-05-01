package domain

type Bookmark struct {
	Title string
	Content string
	URL string
}

type BookmarkRepository interface {
	Add(bookmark *Bookmark) error
}
