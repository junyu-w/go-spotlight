package models

import (
	"time"
)

type FileRecord struct {
	Path        string
	CreatedTime time.Time
	UpdatedTime time.Time
}
