package client

import (
	"fmt"
	"github.com/blevesearch/bleve"
	"github.com/blevesearch/bleve/search/query"
	"regexp"
	"strconv"
	"time"
)

var timeRangeRegexp = regexp.MustCompile(`-(?P<dl>\d+)~-*(?P<dr>\d+)`)

func compileQuery(timeRange, fileExtension, hint string) *query.ConjunctionQuery {
	start, end := parseTimeRange(timeRange)
	qs := make([]query.Query, 0, 3)
	if timeRange != "" {
		qs = append(qs, query.NewDateRangeQuery(start, end))
	}
	if fileExtension != "" {
		qs = append(qs, query.NewQueryStringQuery(fmt.Sprintf("Extension:%s", fileExtension)))
	}
	if hint != "" {
		qs = append(qs, query.NewQueryStringQuery(hint))
	}
	return bleve.NewConjunctionQuery(qs...)
}

func parseTimeRange(t string) (time.Time, time.Time) {
	temp := timeRangeRegexp.FindStringSubmatch(t)
	dr, _ := strconv.Atoi(temp[0])
	dl, _ := strconv.Atoi(temp[1])
	start := time.Now().AddDate(0, 0, -1*dl)
	end := time.Now().AddDate(0, 0, -1*dr)
	return start, end
}
