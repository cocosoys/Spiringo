package orm

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// 中文：MongoConfig 定义当前包使用的数据结构或接口。
// English: MongoConfig defines a data structure or interface used by this package.
type MongoConfig struct {
	// 中文：URI 保存当前结构中的配置或数据值。
	// English: URI stores a configuration or data value for this struct.
	URI string
	// 中文：Database 保存当前结构中的配置或数据值。
	// English: Database stores a configuration or data value for this struct.
	Database string
	// 中文：Timeout 保存当前结构中的配置或数据值。
	// English: Timeout stores a configuration or data value for this struct.
	Timeout time.Duration
}

// 中文：DocumentStore 定义当前包使用的数据结构或接口。
// English: DocumentStore defines a data structure or interface used by this package.
type DocumentStore interface {
	// 中文：InsertOne 声明该接口需要实现的行为。
	// English: InsertOne declares behavior required by this interface.
	InsertOne(ctx context.Context, collection string, document any) (any, error)
	// 中文：FindOne 声明该接口需要实现的行为。
	// English: FindOne declares behavior required by this interface.
	FindOne(ctx context.Context, collection string, filter any, out any) error
	// 中文：UpdateOne 声明该接口需要实现的行为。
	// English: UpdateOne declares behavior required by this interface.
	UpdateOne(ctx context.Context, collection string, filter any, update any) (int64, error)
	// 中文：DeleteOne 声明该接口需要实现的行为。
	// English: DeleteOne declares behavior required by this interface.
	DeleteOne(ctx context.Context, collection string, filter any) (int64, error)
	// 中文：Ping 声明该接口需要实现的行为。
	// English: Ping declares behavior required by this interface.
	Ping(ctx context.Context) error
	// 中文：Close 声明该接口需要实现的行为。
	// English: Close declares behavior required by this interface.
	Close(ctx context.Context) error
}

// 中文：MongoStore 定义当前包使用的数据结构或接口。
// English: MongoStore defines a data structure or interface used by this package.
type MongoStore struct {
	// 中文：client 保存当前结构中的配置或数据值。
	// English: client stores a configuration or data value for this struct.
	client *mongo.Client
	// 中文：db 保存当前结构中的配置或数据值。
	// English: db stores a configuration or data value for this struct.
	db *mongo.Database
}

// 中文：NewMongoStore 创建并返回对应组件实例。
// English: NewMongoStore creates and returns the corresponding component instance.
func NewMongoStore(ctx context.Context, cfg MongoConfig) (*MongoStore, error) {
	if strings.TrimSpace(cfg.URI) == "" {
		return nil, fmt.Errorf("mongodb uri is required")
	}
	if strings.TrimSpace(cfg.Database) == "" {
		return nil, fmt.Errorf("mongodb database is required")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	client, err := mongo.Connect(options.Client().ApplyURI(cfg.URI))
	if err != nil {
		return nil, fmt.Errorf("mongodb connect: %w", err)
	}
	return &MongoStore{
		client: client,
		db:     client.Database(cfg.Database),
	}, nil
}

// 中文：InsertOne 执行当前包中的对应流程。
// English: InsertOne executes the corresponding workflow in this package.
func (s *MongoStore) InsertOne(ctx context.Context, collection string, document any) (any, error) {
	coll, err := s.collection(collection)
	if err != nil {
		return nil, err
	}
	result, err := coll.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}
	return result.InsertedID, nil
}

// 中文：FindOne 执行当前包中的对应流程。
// English: FindOne executes the corresponding workflow in this package.
func (s *MongoStore) FindOne(ctx context.Context, collection string, filter any, out any) error {
	coll, err := s.collection(collection)
	if err != nil {
		return err
	}
	if filter == nil {
		filter = bson.M{}
	}
	return coll.FindOne(ctx, filter).Decode(out)
}

// 中文：UpdateOne 执行当前包中的对应流程。
// English: UpdateOne executes the corresponding workflow in this package.
func (s *MongoStore) UpdateOne(ctx context.Context, collection string, filter any, update any) (int64, error) {
	coll, err := s.collection(collection)
	if err != nil {
		return 0, err
	}
	if filter == nil {
		filter = bson.M{}
	}
	result, err := coll.UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return 0, err
	}
	return result.ModifiedCount, nil
}

// 中文：DeleteOne 执行当前包中的对应流程。
// English: DeleteOne executes the corresponding workflow in this package.
func (s *MongoStore) DeleteOne(ctx context.Context, collection string, filter any) (int64, error) {
	coll, err := s.collection(collection)
	if err != nil {
		return 0, err
	}
	if filter == nil {
		filter = bson.M{}
	}
	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return 0, err
	}
	return result.DeletedCount, nil
}

// 中文：Ping 执行当前包中的对应流程。
// English: Ping executes the corresponding workflow in this package.
func (s *MongoStore) Ping(ctx context.Context) error {
	if s == nil || s.client == nil {
		return fmt.Errorf("mongodb client is not initialized")
	}
	return s.client.Ping(ctx, nil)
}

// 中文：Close 执行当前包中的对应流程。
// English: Close executes the corresponding workflow in this package.
func (s *MongoStore) Close(ctx context.Context) error {
	if s == nil || s.client == nil {
		return nil
	}
	return s.client.Disconnect(ctx)
}

// 中文：collection 执行当前包中的对应流程。
// English: collection executes the corresponding workflow in this package.
func (s *MongoStore) collection(name string) (*mongo.Collection, error) {
	if s == nil || s.db == nil {
		return nil, fmt.Errorf("mongodb database is not initialized")
	}
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("mongodb collection is required")
	}
	return s.db.Collection(name), nil
}
