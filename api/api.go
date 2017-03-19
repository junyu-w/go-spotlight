package api

import (
	"encoding/json"
	"github.com/DrakeW/go-spotlight/models"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"github.com/syndtr/goleveldb/leveldb/errors"
	"strings"
)

var (
	OpenIndexError      = errors.New("ERROR: couldn't open index")
	JsonConversionError = errors.New("ERROR: couldn't convert result to JSON")
)

func compileQuery(q string) query.Query {
	items := strings.Split(q, " ")
	queries := make([]query.Query, len(items), len(items))
	for i, item := range items {
		queries[i] = query.NewQueryStringQuery(item)
	}
	conjQuery := bleve.NewConjunctionQuery(queries...)
	return conjQuery
}

func RunQuery(dir, queryStr string) (string, error) {
	q := compileQuery(queryStr)
	fr_index, err := models.GetFrIndex(dir)
	defer fr_index.Close()
	if err != nil {
		return "", OpenIndexError
	}
	searchRequest := bleve.NewSearchRequest(q)
	searchRequest.Highlight = bleve.NewHighlightWithStyle("html")
	searchResult, err := fr_index.Search(searchRequest)
	res, err := json.Marshal(searchResult)
	if err != nil {
		return "", JsonConversionError
	}
	return string(res), nil
}
