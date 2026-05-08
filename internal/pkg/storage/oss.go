package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"
)

// 中文：defaultContentType 声明当前包使用的常量。
// English: defaultContentType declares constants used by this package.
const defaultContentType = "application/octet-stream"

// 中文：OSS 定义当前包使用的数据结构或接口。
// English: OSS defines a data structure or interface used by this package.
// OSS is the blueprint-style object storage contract with a configured bucket.
type OSS interface {
	// 中文：Put 声明该接口需要实现的行为。
	// English: Put declares behavior required by this interface.
	Put(ctx context.Context, key string, reader io.Reader, size int64) error
	// 中文：Get 声明该接口需要实现的行为。
	// English: Get declares behavior required by this interface.
	Get(ctx context.Context, key string) (io.ReadCloser, error)
	// 中文：Delete 声明该接口需要实现的行为。
	// English: Delete declares behavior required by this interface.
	Delete(ctx context.Context, key string) error
	// 中文：PresignedURL 声明该接口需要实现的行为。
	// English: PresignedURL declares behavior required by this interface.
	PresignedURL(ctx context.Context, key string, expire time.Duration) (string, error)
	// 中文：Exists 声明该接口需要实现的行为。
	// English: Exists declares behavior required by this interface.
	Exists(ctx context.Context, key string) (bool, error)
	// 中文：List 声明该接口需要实现的行为。
	// English: List declares behavior required by this interface.
	List(ctx context.Context, prefix string, limit int) ([]string, error)
}

// 中文：BucketStorage 定义当前包使用的数据结构或接口。
// English: BucketStorage defines a data structure or interface used by this package.
// BucketStorage adapts the lower-level Storage interface to OSS by binding one
// bucket at construction time.
type BucketStorage struct {
	// 中文：backend 保存当前结构中的配置或数据值。
	// English: backend stores a configuration or data value for this struct.
	backend Storage
	// 中文：bucket 保存当前结构中的配置或数据值。
	// English: bucket stores a configuration or data value for this struct.
	bucket string
}

// 中文：_ 声明当前包使用的变量。
// English: _ declares variables used by this package.
var _ OSS = (*BucketStorage)(nil)

// 中文：NewBucketStorage 创建并返回对应组件实例。
// English: NewBucketStorage creates and returns the corresponding component instance.
func NewBucketStorage(backend Storage, bucket string) (*BucketStorage, error) {
	if backend == nil {
		return nil, fmt.Errorf("storage backend is required")
	}
	bucket = strings.TrimSpace(bucket)
	if bucket == "" {
		return nil, fmt.Errorf("storage bucket is required")
	}
	return &BucketStorage{backend: backend, bucket: bucket}, nil
}

// 中文：NewOSS 创建并返回对应组件实例。
// English: NewOSS creates and returns the corresponding component instance.
func NewOSS(backend Storage, bucket string) (OSS, error) {
	return NewBucketStorage(backend, bucket)
}

// 中文：Bucket 执行当前包中的对应流程。
// English: Bucket executes the corresponding workflow in this package.
func (s *BucketStorage) Bucket() string {
	if s == nil {
		return ""
	}
	return s.bucket
}

// 中文：Put 执行当前包中的对应流程。
// English: Put executes the corresponding workflow in this package.
func (s *BucketStorage) Put(ctx context.Context, key string, reader io.Reader, size int64) error {
	if err := s.validateKey(key); err != nil {
		return err
	}
	if reader == nil {
		return fmt.Errorf("storage reader is required")
	}
	contentType := contentTypeForKey(key)
	if size >= 0 {
		return s.backend.UploadStream(ctx, s.bucket, key, reader, size, contentType)
	}
	data, err := io.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("read object data: %w", err)
	}
	return s.backend.Upload(ctx, s.bucket, key, data, contentType)
}

// 中文：Get 执行当前包中的对应流程。
// English: Get executes the corresponding workflow in this package.
func (s *BucketStorage) Get(ctx context.Context, key string) (io.ReadCloser, error) {
	if err := s.validateKey(key); err != nil {
		return nil, err
	}
	data, err := s.backend.Download(ctx, s.bucket, key)
	if err != nil {
		return nil, err
	}
	return io.NopCloser(bytes.NewReader(data)), nil
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (s *BucketStorage) Delete(ctx context.Context, key string) error {
	if err := s.validateKey(key); err != nil {
		return err
	}
	return s.backend.Delete(ctx, s.bucket, key)
}

// 中文：PresignedURL 执行当前包中的对应流程。
// English: PresignedURL executes the corresponding workflow in this package.
func (s *BucketStorage) PresignedURL(ctx context.Context, key string, expire time.Duration) (string, error) {
	if err := s.validateKey(key); err != nil {
		return "", err
	}
	if expire <= 0 {
		return "", fmt.Errorf("presigned url expiry must be positive")
	}
	seconds := int(expire / time.Second)
	if expire%time.Second != 0 {
		seconds++
	}
	if seconds <= 0 {
		seconds = 1
	}
	return s.backend.PresignedURL(ctx, s.bucket, key, seconds)
}

// 中文：Exists 执行当前包中的对应流程。
// English: Exists executes the corresponding workflow in this package.
func (s *BucketStorage) Exists(ctx context.Context, key string) (bool, error) {
	if err := s.validateKey(key); err != nil {
		return false, err
	}
	return s.backend.Exists(ctx, s.bucket, key)
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (s *BucketStorage) List(ctx context.Context, prefix string, limit int) ([]string, error) {
	if s == nil || s.backend == nil {
		return nil, fmt.Errorf("storage backend is required")
	}
	return s.backend.List(ctx, s.bucket, prefix, limit)
}

// 中文：validateKey 执行当前包中的对应流程。
// English: validateKey executes the corresponding workflow in this package.
func (s *BucketStorage) validateKey(key string) error {
	if s == nil || s.backend == nil {
		return fmt.Errorf("storage backend is required")
	}
	if s.bucket == "" {
		return fmt.Errorf("storage bucket is required")
	}
	if strings.TrimSpace(key) == "" {
		return fmt.Errorf("storage key is required")
	}
	return nil
}

// 中文：contentTypeForKey 执行当前包中的对应流程。
// English: contentTypeForKey executes the corresponding workflow in this package.
func contentTypeForKey(key string) string {
	contentType := mime.TypeByExtension(filepath.Ext(key))
	if contentType == "" {
		return defaultContentType
	}
	return contentType
}
