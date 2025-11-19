package internal

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func NewObjectStore(address string) *FileStore {
	store := FileStore{address: address, Context: context.Background(), Connected: false, buckets: make(map[string]bool)}
	return &store
}

func (store *FileStore) Connect(username string, password string) error {
	client, err := minio.New(store.address, &minio.Options{
		Creds:           credentials.NewStaticV4(username, password, ""),
		TrailingHeaders: true,
		Secure:          false,
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

// Uploads to FileStore by redirecting to a temporary file before uploading
// this increases CPU costs, but prevents large data amounts being loaded into
// memory
func (store *FileStore) UploadFS(context context.Context, bucket string, metadata *Metadata, reader io.Reader) error {
	store.createBucketIfNotExists(bucket)

	file, err := os.CreateTemp("", "tmpfile-")
	if err != nil {
		return err
	}

	defer file.Close()
	defer os.Remove(file.Name())

	writer := bufio.NewWriter(file)
	buf := make([]byte, 1024)
	for {
		n, err := reader.Read(buf)
		if err != nil && err != io.EOF {
			return err
		}

		if n == 0 { // very suspect
			break
		}

		if _, err := writer.Write(buf[:n]); err != nil {
			return err
		}
	}

	if err = writer.Flush(); err != nil {
		return err
	}

	// Now we do read :joy:
	uuid, err := store.Database.UploadImageMeta(metadata) // abort now if we can't post metadata
	if err != nil {
		return err
	}

	str := uuid.String()

	_, err = store.Client.FPutObject(context, bucket, str, file.Name(), minio.PutObjectOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (store *FileStore) createBucketIfNotExists(bucket string) error {
	_, ok := store.buckets[bucket]
	if ok {
		return nil
	}

	result, err := store.Client.BucketExists(context.Background(), bucket)
	if err != nil {
		return err
	}

	if result {
		store.buckets[bucket] = true
		return nil
	}

	err = store.Client.MakeBucket(context.Background(), bucket, minio.MakeBucketOptions{
		Region: "us-east-1",
	})

	if err != nil {
		return err
	}

	store.buckets[bucket] = true
	return nil
}

// Fun fact client#AppendObject does NOT WORK on the community edition. There is ABSOLUTELY 0 warning about this
// in ANY form of documentation. Whatever fuck head thought that was a real piece of shit
// func (store *FileStore) PartitionedUpload(bucket string, metadata Metadata, reader io.Reader) error {
// 	store.createBucketIfNotExists(bucket)
// 	client := store.Client

// 	// upload empty object
// 	_, err := client.PutObject(context.Background(), bucket, metadata.Name, nil, 0, minio.PutObjectOptions{})
// 	if err != nil {
// 		return err
// 	}

// 	partSize := int64(2 * 1024 * 1024) // 2mb parts
// 	buffer := make([]byte, partSize)
// 	totalUploaded := int64(0)

// 	for {
// 		n, err := reader.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		if err != nil {
// 			store.deletePartialUpload(bucket, metadata)
// 			return err
// 		}

// 		_, err = client.AppendObject(context.Background(), bucket, metadata.Name, reader, int64(n), minio.AppendObjectOptions{}) // you need premium :/
// 		if err != nil {
// 			store.deletePartialUpload(bucket, metadata)
// 			println(metadata.Name)
// 			log.Printf("[append error] %v", err)
// 			return err
// 		}

// 		totalUploaded += int64(n)
// 	}

// 	log.Println("Uploaded ", metadata.Name, " to the bucket ", bucket, " successfully. Wrote ", totalUploaded, " bytes")
// 	return nil
// }

// func (store *FileStore) deletePartialUpload(bucket string, metadata Metadata) error {
// 	// info, err := store.Client.StatObject(context.Background(), bucket, metadata.Name, minio.StatObjectOptions{})
// 	// if err != nil {
// 	// 	log.Println(err)
// 	// 	return err
// 	// }

// 	// println(info.Owner.DisplayName)
// 	// store.Client.RemoveObject(context.Background(), bucket, metadata.Name, minio.RemoveObjectOptions{})
// 	return nil
// }
