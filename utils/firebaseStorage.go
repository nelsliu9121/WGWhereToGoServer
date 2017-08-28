package utils

import (
	"context"
	"fmt"
	"io"

	"cloud.google.com/go/storage"
	log "github.com/sirupsen/logrus"
)

// Storage Reference
type Storage struct {
	Bucket string
}

func (s *Storage) resource(path string) string {
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o?name=%s", s.Bucket, path)
}

func (s *Storage) client(ctx context.Context) (*storage.Client, error) {
	storage, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}

	return storage, err
}

// Put will store a file in Firebase Storage
func (s *Storage) Put(ctx context.Context, data io.Reader, path string) (*storage.ObjectHandle, *storage.ObjectAttrs, error) {
	client, err := s.client(ctx)
	if err != nil {
		log.WithError(err).Panic()
	}
	defer client.Close()

	bkt := client.Bucket(s.Bucket)
	obj := bkt.Object(path)
	writer := obj.NewWriter(ctx)
	if _, err := io.Copy(writer, data); err != nil {
		log.WithError(err).Panic()
	}
	if err := writer.Close(); err != nil {
		log.WithError(err).Panic()
	}
	if err := obj.ACL().Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.WithError(err).Panic()
	}
	attrs, err := obj.Attrs(ctx)
	return obj, attrs, err
}
