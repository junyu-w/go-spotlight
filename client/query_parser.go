package client

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"time"
)

func compileQuery(timeRange, fileExtension, hint string) *query.ConjunctionQuery {
	start, end := parseTimeRange(timeRange)
	q1 := query.NewDateRangeQuery(start, end)
	q2 := query.NewQueryStringQuery(fmt.Sprintf("Extension:%s", fileExtension))
	q3 := query.NewFuzzyQuery(hint)
	return bleve.NewConjunctionQuery(q1, q2, q3)
}

func parseTimeRange(t string) (time.Time, time.Time) {
	start, _ := time.Parse(time.RFC3339, "2017-01-02T15:04:05Z07:00")
	end, _ := time.Parse(time.RFC3339, "2017-03-11T15:04:05Z07:00")
	return start, end
}
