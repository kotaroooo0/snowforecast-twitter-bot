package searcher

import (
	"context"
	"log"

	"github.com/elastic/go-elasticsearch/v7"
	elasticsearch7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

type SnowResortSearcherEsImpl struct {
}

func NewSnowResortSearcherEsImple() {
	es, _ := elasticsearch7.NewDefaultClient()
	log.Println(elasticsearch.Version)
	log.Println(es.Info())
	log.Println(es.Cat.Indices())

	req := esapi.SearchRequest{
		SuggestText: "kagura",
	}
	ctx := context.Background()
	log.Println(req.Do(ctx, es))
}
