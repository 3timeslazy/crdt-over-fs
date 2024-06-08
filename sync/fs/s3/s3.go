package s3

import (
	"bytes"
	"io"

	"github.com/3timeslazy/crdt-over-fs/sync/fs"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
)

type FS struct {
	client *s3.S3
	bucket *string
}

func NewFS(client *s3.S3, bucket string) *FS {
	return &FS{
		client: client,
		bucket: &bucket,
	}
}

type DirEntry struct {
	object *s3.Object
}

func (dr *DirEntry) Name() string {
	return *dr.object.Key
}

func (dr *DirEntry) IsDir() bool {
	// TODO: sub-buckets?
	return false
}

func (s3fs *FS) ReadDir(name string) ([]fs.DirEntry, error) {
	out, err := s3fs.client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: s3fs.bucket,
	})
	if err != nil {
		return nil, err
	}

	entries := []fs.DirEntry{}
	for _, obj := range out.Contents {
		entries = append(entries, &DirEntry{obj})
	}

	return entries, nil
}

func (s3fs *FS) MakeDir(name string) error {
	panic("not implemented")
}

func (s3fs *FS) WriteFile(name string, data []byte) error {
	_, err := s3fs.client.PutObject(&s3.PutObjectInput{
		Key:    aws.String(name),
		Bucket: s3fs.bucket,
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}

func (s3fs *FS) ReadFile(name string) ([]byte, error) {
	out, err := s3fs.client.GetObject(&s3.GetObjectInput{
		Bucket: s3fs.bucket,
		Key:    aws.String(name),
	})
	if err != nil {
		if e, ok := err.(awserr.Error); ok {
			switch e.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, fs.ErrNotExist
			}
		}
		return nil, err
	}
	defer out.Body.Close()

	data, err := io.ReadAll(out.Body)
	return data, err
}
