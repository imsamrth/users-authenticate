package database

import (
	"context"
	"log"
	"os"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

func StorageInstance() (sc *storage.BucketHandle) {
	ctx := context.Background()
	bucket := os.Getenv("BUCKET")

	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile("keys.json"))
	if err != nil {
		log.Fatal(err)
	}

	return storageClient.Bucket(bucket)

}

var BucketHandle *storage.BucketHandle = StorageInstance()
