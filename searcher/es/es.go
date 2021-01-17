package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kotaroooo0/snowforecast-twitter-bot/searcher"
	"github.com/olivere/elastic/v7"
)

var (
	targetIndex = "snow_resorts_alias"
	targetField = []string{"name", "search_key"}
)

type SearcherEsImpl struct {
	Client *elastic.Client
}

func NewSnowResortSearcherEsImpl() (*SearcherEsImpl, error) {
	// ref: https://github.com/olivere/elastic/wiki/Docker
	client, err := elastic.NewClient(elastic.SetSniff(false))
	if err != nil {
		return nil, err
	}
	return &SearcherEsImpl{
		Client: client,
	}, nil
}

func (s *SearcherEsImpl) FindSimilarSnowResort(source string) (*searcher.SnowResortDto, error) {
	multiMatchQuery := elastic.NewMultiMatchQuery(source, targetField...).Type("most_fields")
	res, err := s.Client.Search().Index(targetIndex).Query(multiMatchQuery).Size(1).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if res.Hits.TotalHits.Value == 0 {
		return nil, fmt.Errorf("error: document not found")
	}

	var sr searcher.SnowResortDto
	if err = json.Unmarshal(res.Hits.Hits[0].Source, &sr); err != nil {
		log.Println(err)
	}
	return &sr, nil
}
