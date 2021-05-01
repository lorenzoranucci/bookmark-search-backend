package elasticsearch

import (
	"context"

	"github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/domain"
	"github.com/olivere/elastic/v7"
)

func NewBookmarkRepository(elasticClient *elastic.Client) *BookmarkRepository {
	return &BookmarkRepository{ElasticClient: elasticClient}
}

type BookmarkRepository struct {
	ElasticClient *elastic.Client
}

type Bookmark struct {
	UID string `json:"uid"`
	URL string `json:"url"`
	Content string `json:"body"`
	Title string `json:"title"`
}

func (pr *BookmarkRepository) Add(bookmark *domain.Bookmark) error {
	request := elastic.NewBulkIndexRequest().Index("page").
		Doc(mapDomainBookmarkWithElasticsearchBookmark(bookmark)).
		Id(bookmark.UID)

	bulkResponse, err := pr.ElasticClient.
		Bulk().
		Add(request).
		Do(context.Background())

	if err != nil {
		return err
	}

	if bulkResponse.Errors == false {
		return nil
	}

	for _, items := range bulkResponse.Items {
		for _, item := range items {
			if item.Error != nil {
				return err
			}
		}
	}

	return nil
}

func mapDomainBookmarkWithElasticsearchBookmark(domainBookmark *domain.Bookmark) Bookmark {
	return Bookmark{
		UID: domainBookmark.UID,
		URL:  domainBookmark.URL,
		Content: domainBookmark.Content,
		Title: domainBookmark.Title,
	}
}
