package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// 中文：Storage 定义当前包使用的数据结构或接口。
// English: Storage defines a data structure or interface used by this package.
// Storage 对象存储接口
type Storage interface {
	// 中文：Upload 声明该接口需要实现的行为。
	// English: Upload declares behavior required by this interface.
	// Upload 上传文件
	Upload(ctx context.Context, bucket, key string, data []byte, contentType string) error
	// 中文：UploadStream 声明该接口需要实现的行为。
	// English: UploadStream declares behavior required by this interface.
	// UploadStream 流式上传
	UploadStream(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) error
	// 中文：Download 声明该接口需要实现的行为。
	// English: Download declares behavior required by this interface.
	// Download 下载文件
	Download(ctx context.Context, bucket, key string) ([]byte, error)
	// 中文：Delete 声明该接口需要实现的行为。
	// English: Delete declares behavior required by this interface.
	// Delete 删除文件
	Delete(ctx context.Context, bucket, key string) error
	// 中文：GetURL 声明该接口需要实现的行为。
	// English: GetURL declares behavior required by this interface.
	// GetURL 获取文件访问URL
	GetURL(bucket, key string) string
	// 中文：PresignedURL 声明该接口需要实现的行为。
	// English: PresignedURL declares behavior required by this interface.
	// PresignedURL 获取预签名URL（限时访问）
	PresignedURL(ctx context.Context, bucket, key string, expirySeconds int) (string, error)
	// 中文：Exists 声明该接口需要实现的行为。
	// English: Exists declares behavior required by this interface.
	// Exists checks whether an object exists.
	Exists(ctx context.Context, bucket, key string) (bool, error)
	// 中文：List 声明该接口需要实现的行为。
	// English: List declares behavior required by this interface.
	// List returns object keys under a prefix. limit <= 0 means no explicit limit.
	List(ctx context.Context, bucket, prefix string, limit int) ([]string, error)
}

// 中文：MinIOStorage 定义当前包使用的数据结构或接口。
// English: MinIOStorage defines a data structure or interface used by this package.
// MinIOStorage MinIO对象存储实现
type MinIOStorage struct {
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *minio.Client
	// 中文：endpoint 保存当前结构中的配置或数据值。
	// English: endpoint stores a configuration or data value for this struct.
	endpoint string
	// 中文：useSSL 保存当前结构中的配置或数据值。
	// English: useSSL stores a configuration or data value for this struct.
	useSSL bool
}

// 中文：MinIOConfig 定义当前包使用的数据结构或接口。
// English: MinIOConfig defines a data structure or interface used by this package.
// MinIOConfig MinIO配置
type MinIOConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string `yaml:"endpoint"`
	// 中文：AccessKey 保存当前结构中的配置或数据值。
	// English: AccessKey stores a configuration or data value for this struct.
	AccessKey string `yaml:"access_key"`
	// 中文：SecretKey 保存当前结构中的配置或数据值。
	// English: SecretKey stores a configuration or data value for this struct.
	SecretKey string `yaml:"secret_key"`
	// 中文：UseSSL 保存当前结构中的配置或数据值。
	// English: UseSSL stores a configuration or data value for this struct.
	UseSSL bool `yaml:"use_ssl"`
}

// 中文：CephConfig 定义当前包使用的数据结构或接口。
// English: CephConfig defines a data structure or interface used by this package.
// CephConfig contains S3-compatible Ceph RGW settings.
type CephConfig struct {
	// 中文：Endpoint 保存当前结构中的配置或数据值。
	// English: Endpoint stores a configuration or data value for this struct.
	Endpoint string `yaml:"endpoint"`
	// 中文：AccessKey 保存当前结构中的配置或数据值。
	// English: AccessKey stores a configuration or data value for this struct.
	AccessKey string `yaml:"access_key"`
	// 中文：SecretKey 保存当前结构中的配置或数据值。
	// English: SecretKey stores a configuration or data value for this struct.
	SecretKey string `yaml:"secret_key"`
	// 中文：UseSSL 保存当前结构中的配置或数据值。
	// English: UseSSL stores a configuration or data value for this struct.
	UseSSL bool `yaml:"use_ssl"`
	// 中文：PublicURL 保存当前结构中的配置或数据值。
	// English: PublicURL stores a configuration or data value for this struct.
	PublicURL string `yaml:"public_url"`
}

// 中文：NewMinIOStorage 创建并返回对应组件实例。
// English: NewMinIOStorage creates and returns the corresponding component instance.
// NewMinIOStorage 创建MinIO存储
func NewMinIOStorage(cfg MinIOConfig) (*MinIOStorage, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("minio: create client: %w", err)
	}
	return &MinIOStorage{client: client, endpoint: cfg.Endpoint, useSSL: cfg.UseSSL}, nil
}

// 中文：Upload 执行当前包中的对应流程。
// English: Upload executes the corresponding workflow in this package.
func (s *MinIOStorage) Upload(ctx context.Context, bucket, key string, data []byte, contentType string) error {
	// 确保bucket存在
	if err := s.ensureBucket(ctx, bucket); err != nil {
		return err
	}

	_, err := s.client.PutObject(ctx, bucket, key, bytes.NewReader(data), int64(len(data)), minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// 中文：UploadStream 执行当前包中的对应流程。
// English: UploadStream executes the corresponding workflow in this package.
func (s *MinIOStorage) UploadStream(ctx context.Context, bucket, key string, reader io.Reader, size int64, contentType string) error {
	if err := s.ensureBucket(ctx, bucket); err != nil {
		return err
	}

	_, err := s.client.PutObject(ctx, bucket, key, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	return err
}

// 中文：Download 执行当前包中的对应流程。
// English: Download executes the corresponding workflow in this package.
func (s *MinIOStorage) Download(ctx context.Context, bucket, key string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, bucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, err
	}
	defer obj.Close()
	return io.ReadAll(obj)
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (s *MinIOStorage) Delete(ctx context.Context, bucket, key string) error {
	return s.client.RemoveObject(ctx, bucket, key, minio.RemoveObjectOptions{})
}

// 中文：GetURL 执行当前包中的对应流程。
// English: GetURL executes the corresponding workflow in this package.
func (s *MinIOStorage) GetURL(bucket, key string) string {
	scheme := "http"
	if s.useSSL {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s/%s/%s", scheme, s.endpoint, bucket, key)
}

// 中文：PresignedURL 执行当前包中的对应流程。
// English: PresignedURL executes the corresponding workflow in this package.
func (s *MinIOStorage) PresignedURL(ctx context.Context, bucket, key string, expirySeconds int) (string, error) {
	reqParams := make(url.Values)
	u, err := s.client.PresignedGetObject(ctx, bucket, key, time.Duration(expirySeconds)*time.Second, reqParams)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

// 中文：Exists 执行当前包中的对应流程。
// English: Exists executes the corresponding workflow in this package.
func (s *MinIOStorage) Exists(ctx context.Context, bucket, key string) (bool, error) {
	if bucket == "" || key == "" {
		return false, nil
	}
	if _, err := s.client.StatObject(ctx, bucket, key, minio.StatObjectOptions{}); err != nil {
		resp := minio.ToErrorResponse(err)
		if resp.StatusCode == http.StatusNotFound || resp.Code == "NoSuchKey" || resp.Code == "NoSuchBucket" {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (s *MinIOStorage) List(ctx context.Context, bucket, prefix string, limit int) ([]string, error) {
	if bucket == "" {
		return nil, fmt.Errorf("bucket is required")
	}
	keys := make([]string, 0)
	for obj := range s.client.ListObjects(ctx, bucket, minio.ListObjectsOptions{
		Prefix:    prefix,
		Recursive: true,
	}) {
		if obj.Err != nil {
			return nil, obj.Err
		}
		keys = append(keys, obj.Key)
		if limit > 0 && len(keys) >= limit {
			break
		}
	}
	return keys, nil
}

// 中文：ensureBucket 执行当前包中的对应流程。
// English: ensureBucket executes the corresponding workflow in this package.
func (s *MinIOStorage) ensureBucket(ctx context.Context, bucket string) error {
	exists, err := s.client.BucketExists(ctx, bucket)
	if err != nil {
		return fmt.Errorf("minio: check bucket: %w", err)
	}
	if !exists {
		if err := s.client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return fmt.Errorf("minio: create bucket: %w", err)
		}
	}
	return nil
}

// 中文：CephStorage 定义当前包使用的数据结构或接口。
// English: CephStorage defines a data structure or interface used by this package.
type CephStorage struct {
	// 中文：*MinIOStorage 嵌入复用该类型提供的能力。
	// English: *MinIOStorage embeds reusable behavior from that type.
	*MinIOStorage
	// 中文：publicURL 保存当前结构中的配置或数据值。
	// English: publicURL stores a configuration or data value for this struct.
	publicURL string
}

// 中文：NewCephStorage 创建并返回对应组件实例。
// English: NewCephStorage creates and returns the corresponding component instance.
func NewCephStorage(cfg CephConfig) (*CephStorage, error) {
	if strings.TrimSpace(cfg.Endpoint) == "" {
		return nil, fmt.Errorf("ceph: endpoint is required")
	}
	backend, err := NewMinIOStorage(MinIOConfig{
		Endpoint:  cfg.Endpoint,
		AccessKey: cfg.AccessKey,
		SecretKey: cfg.SecretKey,
		UseSSL:    cfg.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("ceph: create s3 compatible client: %w", err)
	}
	return &CephStorage{
		MinIOStorage: backend,
		publicURL:    strings.TrimRight(strings.TrimSpace(cfg.PublicURL), "/"),
	}, nil
}

// 中文：GetURL 执行当前包中的对应流程。
// English: GetURL executes the corresponding workflow in this package.
func (s *CephStorage) GetURL(bucket, key string) string {
	if s.publicURL == "" {
		return s.MinIOStorage.GetURL(bucket, key)
	}
	return fmt.Sprintf("%s/%s/%s", s.publicURL, bucket, key)
}
