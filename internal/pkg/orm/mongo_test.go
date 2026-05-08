package orm

import (
	"context"
	"testing"
	"time"
)

// 中文：TestNewMongoStoreValidatesConfig 验证相关行为符合预期。
// English: TestNewMongoStoreValidatesConfig verifies the related behavior.
func TestNewMongoStoreValidatesConfig(t *testing.T) {
	if _, err := NewMongoStore(context.Background(), MongoConfig{}); err == nil {
		t.Fatalf("expected uri error")
	}
	if _, err := NewMongoStore(context.Background(), MongoConfig{URI: "mongodb://127.0.0.1:27017"}); err == nil {
		t.Fatalf("expected database error")
	}
}

// 中文：TestMongoStoreCollectionValidation 验证相关行为符合预期。
// English: TestMongoStoreCollectionValidation verifies the related behavior.
func TestMongoStoreCollectionValidation(t *testing.T) {
	store, err := NewMongoStore(context.Background(), MongoConfig{
		URI:      "mongodb://127.0.0.1:27017",
		Database: "spiringo",
		Timeout:  10 * time.Millisecond,
	})
	if err != nil {
		t.Fatalf("new store: %v", err)
	}
	defer store.Close(context.Background())

	if _, err := store.collection(""); err == nil {
		t.Fatalf("expected collection error")
	}
	if _, err := store.collection("audit_logs"); err != nil {
		t.Fatalf("collection: %v", err)
	}
}
