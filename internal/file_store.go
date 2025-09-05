package internal

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileStore struct {
	address   string
	Context   context.Context
	Client    *minio.Client
	Connected bool
}

func NewObjectStore(address string) *FileStore {
	store := FileStore{address: address, Context: context.Background(), Connected: false}
	return &store
}

func (store *FileStore) Connect(username string, password string) error {
	client, err := minio.New(store.address, &minio.Options{
		Creds:  credentials.NewStaticV4(username, password, ""),
		Secure: false,
	})

	if err != nil {
		return err
	}

	store.Client = client
	store.Connected = true
	return nil
}

func (store *FileStore) Close() error {
	store.Client = nil
	store.Connected = false
	return nil
}
