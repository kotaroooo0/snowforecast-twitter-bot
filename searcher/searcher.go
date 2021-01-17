package searcher

type Searchr interface {
	FindSimilarSnowResort(string) (*SnowResortDto, error)
}

type SnowResortDto struct {
	Name      string `json:"name"`
	SearchKey string `json:"search_key"`
}
