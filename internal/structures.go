package internal

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
)

type ImageMeta struct {
	Id          uuid.UUID
	ImageName   string
	ImageType   string
	Description string
	Tags        []string
}

type Database struct {
	conn *sql.DB
}

type Metadata struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	ImageType   string
}

type IdResponse struct {
	Ids     uuid.UUIDs `json:"ids"`
	Entries int        `json:"entries"`
}

type FileStore struct {
	address   string
	Context   context.Context
	Client    *minio.Client
	Database  *Database
	Connected bool
	buckets   map[string]bool
}
