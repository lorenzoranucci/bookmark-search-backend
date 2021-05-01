package domain

type Bookmark struct {
	UID string
	Title string
	Content string
	URL string
}

type BookmarkRepository interface {
	Add(bookmark *Bookmark) error
}
