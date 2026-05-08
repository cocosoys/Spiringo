package storage

import (
	"context"
	"io"
	"strconv"
	"strings"
	"testing"
	"time"
)

// 中文：_、_ 声明当前包使用的变量。
// English: _、_ declares variables used by this package.
var (
	_ Storage = (*MinIOStorage)(nil)
	_ Storage = (*CephStorage)(nil)
)

// 中文：TestCephStorageGetURLUsesPublicURL 验证相关行为符合预期。
// English: TestCephStorageGetURLUsesPublicURL verifies the related behavior.
func TestCephStorageGetURLUsesPublicURL(t *testing.T) {
	storage, err := NewCephStorage(CephConfig{
		Endpoint:  "127.0.0.1:7480",
		AccessKey: "access",
		SecretKey: "secret",
		PublicURL: "https://cdn.example.com/base/",
	})
	if err != nil {
		t.Fatalf("new ceph storage: %v", err)
	}
	if got := storage.GetURL("bucket", "path/file.txt"); got != "https://cdn.example.com/base/bucket/path/file.txt" {
		t.Fatalf("url = %s", got)
	}
}

// 中文：TestCephStorageRequiresEndpoint 验证相关行为符合预期。
// English: TestCephStorageRequiresEndpoint verifies the related behavior.
func TestCephStorageRequiresEndpoint(t *testing.T) {
	if _, err := NewCephStorage(CephConfig{}); err == nil {
		t.Fatalf("expected endpoint error")
	}
}

// 中文：TestBucketStorageAdaptsBlueprintOSSContract 验证相关行为符合预期。
// English: TestBucketStorageAdaptsBlueprintOSSContract verifies the related behavior.
func TestBucketStorageAdaptsBlueprintOSSContract(t *testing.T) {
	backend := newMemoryStorage()
	oss, err := NewBucketStorage(backend, "assets")
	if err != nil {
		t.Fatalf("NewBucketStorage returned error: %v", err)
	}

	if err := oss.Put(context.Background(), "docs/readme.txt", strings.NewReader("hello"), int64(len("hello"))); err != nil {
		t.Fatalf("Put returned error: %v", err)
	}
	if backend.contentTypes["assets/docs/readme.txt"] != "text/plain; charset=utf-8" {
		t.Fatalf("content type = %q", backend.contentTypes["assets/docs/readme.txt"])
	}
	exists, err := oss.Exists(context.Background(), "docs/readme.txt")
	if err != nil || !exists {
		t.Fatalf("Exists = %v, %v; want true, nil", exists, err)
	}
	reader, err := oss.Get(context.Background(), "docs/readme.txt")
	if err != nil {
		t.Fatalf("Get returned error: %v", err)
	}
	defer reader.Close()
	data, err := io.ReadAll(reader)
	if err != nil {
		t.Fatalf("ReadAll returned error: %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("data = %q", data)
	}
	url, err := oss.PresignedURL(context.Background(), "docs/readme.txt", 1500*time.Millisecond)
	if err != nil {
		t.Fatalf("PresignedURL returned error: %v", err)
	}
	if url != "signed://assets/docs/readme.txt?exp=2" {
		t.Fatalf("url = %q", url)
	}
	keys, err := oss.List(context.Background(), "docs/", 10)
	if err != nil {
		t.Fatalf("List returned error: %v", err)
	}
	if len(keys) != 1 || keys[0] != "docs/readme.txt" {
		t.Fatalf("keys = %#v", keys)
	}
	if err := oss.Delete(context.Background(), "docs/readme.txt"); err != nil {
		t.Fatalf("Delete returned error: %v", err)
	}
}

// 中文：memoryStorage 定义当前包使用的数据结构或接口。
// English: memoryStorage defines a data structure or interface used by this package.
type memoryStorage struct {
	// 中文：objects 保存当前结构中的配置或数据值。
	// English: objects stores a configuration or data value for this struct.
	objects map[string][]byte
	// 中文：contentTypes 保存当前结构中的配置或数据值。
	// English: contentTypes stores a configuration or data value for this struct.
	contentTypes map[string]string
}

// 中文：newMemoryStorage 执行当前包中的对应流程。
// English: newMemoryStorage executes the corresponding workflow in this package.
func newMemoryStorage() *memoryStorage {
	return &memoryStorage{
		objects:      make(map[string][]byte),
		contentTypes: make(map[string]string),
	}
}

// 中文：Upload 执行当前包中的对应流程。
// English: Upload executes the corresponding workflow in this package.
func (s *memoryStorage) Upload(_ context.Context, bucket, key string, data []byte, contentType string) error {
	id := bucket + "/" + key
	s.objects[id] = append([]byte(nil), data...)
	s.contentTypes[id] = contentType
	return nil
}

// 中文：UploadStream 执行当前包中的对应流程。
// English: UploadStream executes the corresponding workflow in this package.
func (s *memoryStorage) UploadStream(ctx context.Context, bucket, key string, reader io.Reader, _ int64, contentType string) error {
	data, err := io.ReadAll(reader)
	if err != nil {
		return err
	}
	return s.Upload(ctx, bucket, key, data, contentType)
}

// 中文：Download 执行当前包中的对应流程。
// English: Download executes the corresponding workflow in this package.
func (s *memoryStorage) Download(_ context.Context, bucket, key string) ([]byte, error) {
	return append([]byte(nil), s.objects[bucket+"/"+key]...), nil
}

// 中文：Delete 执行当前包中的对应流程。
// English: Delete executes the corresponding workflow in this package.
func (s *memoryStorage) Delete(_ context.Context, bucket, key string) error {
	delete(s.objects, bucket+"/"+key)
	return nil
}

// 中文：GetURL 执行当前包中的对应流程。
// English: GetURL executes the corresponding workflow in this package.
func (s *memoryStorage) GetURL(bucket, key string) string {
	return "mem://" + bucket + "/" + key
}

// 中文：PresignedURL 执行当前包中的对应流程。
// English: PresignedURL executes the corresponding workflow in this package.
func (s *memoryStorage) PresignedURL(_ context.Context, bucket, key string, expirySeconds int) (string, error) {
	return "signed://" + bucket + "/" + key + "?exp=" + strconv.Itoa(expirySeconds), nil
}

// 中文：Exists 执行当前包中的对应流程。
// English: Exists executes the corresponding workflow in this package.
func (s *memoryStorage) Exists(_ context.Context, bucket, key string) (bool, error) {
	_, ok := s.objects[bucket+"/"+key]
	return ok, nil
}

// 中文：List 执行当前包中的对应流程。
// English: List executes the corresponding workflow in this package.
func (s *memoryStorage) List(_ context.Context, bucket, prefix string, limit int) ([]string, error) {
	keys := make([]string, 0)
	for id := range s.objects {
		key := strings.TrimPrefix(id, bucket+"/")
		if key == id || !strings.HasPrefix(key, prefix) {
			continue
		}
		keys = append(keys, key)
		if limit > 0 && len(keys) >= limit {
			break
		}
	}
	return keys, nil
}
