package minio

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Minio struct {
	s3Client *minio.Client
	bucket   string
}

func New(endpoint, key, secret, bucket, region string) (*Minio, error) {
	s3Client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(key, secret, ""),
		Region: region,
		Secure: true,
	})
	if err != nil {
		return nil, err
	}

	return &Minio{
		s3Client: s3Client,
		bucket:   bucket,
	}, nil
}

func (t *Minio) read(ctx context.Context, key string) ([]byte, error) {
	obj, err := t.s3Client.GetObject(
		ctx,
		t.bucket, key,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, fmt.Errorf("db: minio-s3, read: %v", err)
	}

	return io.ReadAll(obj)
}

func (t *Minio) Load(key []byte) ([]byte, error) {
	return t.read(context.TODO(), string(key))
}

func (t *Minio) Write(key []byte, value []byte) error {
	_, err := t.s3Client.PutObject(
		context.TODO(),
		t.bucket,
		string(key),
		bytes.NewBuffer(value),
		int64(len(value)),
		minio.PutObjectOptions{},
	)
	if err != nil {
		return fmt.Errorf("db: minio-s3, write: %v", err)
	}
	return nil
}

func (t *Minio) Scan(prefix []byte, openFn func(key []byte, value []byte) error) error {
	objsCn := t.s3Client.ListObjects(
		context.TODO(),
		t.bucket,
		minio.ListObjectsOptions{
			Prefix:    string(prefix),
			Recursive: true,
		},
	)

	for obj := range objsCn {
		value, err := t.read(context.TODO(), obj.Key)
		if err != nil {
			return err
		}

		err = openFn([]byte(obj.Key), value)
		if err != nil {
			return err
		}
	}

	return nil
}
