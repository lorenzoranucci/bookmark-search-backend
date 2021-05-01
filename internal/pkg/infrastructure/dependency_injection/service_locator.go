package dependency_injection

import (
	"crypto/tls"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/application"
	"github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/domain"
	"github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/infrastructure/elasticsearch"
	iHttp "github.com/lorenzoranucci/bookmark-search-backend/internal/pkg/infrastructure/http"
	"github.com/olivere/elastic/v7"
)

var ServiceLocatorInstance = ServiceLocator{}

type ServiceLocator struct {
	domainBookmarkFactory           *domain.BookmarkFactory
	domainBookmarkRepository        domain.BookmarkRepository
	domainResourceBodyProvider      domain.ResourceBodyProvider
	domainBookmarkContentCollector  domain.BookmarkContentCollector
	resourceBodyProvider            *application.ResourceBodyProvider
	bookmarkContentCollector        *application.BookmarkContentCollector
	httpServer                      *http.Server
	bookmarkHandler                 *iHttp.BookmarksHandler
	elasticsearchBookmarkRepository *elasticsearch.BookmarkRepository
	elasticClient                   *elastic.Client
}

func (s *ServiceLocator) DomainBookmarkFactory() *domain.BookmarkFactory {
	if s.domainBookmarkFactory == nil {
		s.domainBookmarkFactory = domain.NewBookmarkFactory(s.DomainResourceBodyProvider(), s.DomainBookmarkContentCollector())
	}

	return s.domainBookmarkFactory
}

func (s *ServiceLocator) DomainBookmarkRepository() (domain.BookmarkRepository, error) {
	return s.ElasticsearchBookmarkRepository()
}

func (s *ServiceLocator) DomainResourceBodyProvider() domain.ResourceBodyProvider {
	return s.ResourceBodyProvider()
}

func (s *ServiceLocator) DomainBookmarkContentCollector() domain.BookmarkContentCollector {
	return s.BookmarkContentCollector()
}

func (s *ServiceLocator) ResourceBodyProvider() *application.ResourceBodyProvider {
	if s.resourceBodyProvider == nil {
		s.resourceBodyProvider = application.NewResourceBodyProvider(&http.Client{})
	}

	return s.resourceBodyProvider
}

func (s *ServiceLocator) BookmarkContentCollector() *application.BookmarkContentCollector {
	if s.bookmarkContentCollector == nil {
		s.bookmarkContentCollector = &application.BookmarkContentCollector{}
	}

	return s.bookmarkContentCollector
}

func (s *ServiceLocator) HttpServer() (*http.Server, error) {
	if s.httpServer == nil {
		port, err := s.HTTPServerPort()
		if err != nil {
			return nil, err
		}

		handler, err := s.BookmarkHandler()
		if err != nil {
			return nil, err
		}

		s.httpServer = iHttp.NewServer(port, handler)
	}

	return s.httpServer, nil
}

func (s *ServiceLocator) BookmarkHandler() (*iHttp.BookmarksHandler, error) {
	if s.bookmarkHandler == nil {
		repository, err := s.DomainBookmarkRepository()
		if err != nil {
			return nil, err
		}
		s.bookmarkHandler = iHttp.NewBookmarksHandler(s.DomainBookmarkFactory(), repository)
	}

	return s.bookmarkHandler, nil
}

func (s *ServiceLocator) ElasticsearchBookmarkRepository() (*elasticsearch.BookmarkRepository, error) {
	if s.elasticsearchBookmarkRepository == nil {
		client, err := s.ElasticClient()
		if err != nil {
			return nil, err
		}

		s.elasticsearchBookmarkRepository = elasticsearch.NewBookmarkRepository(client)
	}

	return s.elasticsearchBookmarkRepository, nil
}

func (s *ServiceLocator) ElasticClient() (*elastic.Client, error) {
	if s.elasticClient == nil {
		ec, err := elastic.NewClient(
			elastic.SetURL(strings.Split(s.ElasticURLs(), ",")...),
			elastic.SetSniff(false),
			elastic.SetHealthcheck(false),
			elastic.SetHttpClient(
				&http.Client{
					Transport: &http.Transport{
						TLSClientConfig: &tls.Config{
							InsecureSkipVerify: true,
						},
					},
				},
			),
		)

		if err != nil {
			return nil, err
		}

		s.elasticClient = ec
	}

	return s.elasticClient, nil
}

func (s *ServiceLocator) ElasticURLs() string {
	return os.Getenv("ELASTICSEARCH_URL")
}

func (s *ServiceLocator) HTTPServerPort() (int, error) {
	portInt, err := strconv.Atoi(os.Getenv("HTTP_SERVER_PORT"))
	if err != nil {
		return 0, err
	}
	return portInt, nil
}




