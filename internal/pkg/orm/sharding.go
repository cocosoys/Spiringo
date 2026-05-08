package orm

import (
	"fmt"
	"hash/fnv"
	"strconv"
	"strings"
)

// 中文：ErrShardKeyRequired、ErrShardNotFound 声明当前包使用的变量。
// English: ErrShardKeyRequired、ErrShardNotFound declares variables used by this package.
var (
	ErrShardKeyRequired = fmt.Errorf("shard key is required")
	ErrShardNotFound    = fmt.Errorf("shard database not found")
)

// 中文：ShardTarget 定义当前包使用的数据结构或接口。
// English: ShardTarget defines a data structure or interface used by this package.
type ShardTarget struct {
	// 中文：Database 保存当前结构中的配置或数据值。
	// English: Database stores a configuration or data value for this struct.
	Database string
	// 中文：Table 保存当前结构中的配置或数据值。
	// English: Table stores a configuration or data value for this struct.
	Table string
	// 中文：Index 保存当前结构中的配置或数据值。
	// English: Index stores a configuration or data value for this struct.
	Index uint64
}

// 中文：ShardStrategy 定义当前包使用的数据结构或接口。
// English: ShardStrategy defines a data structure or interface used by this package.
type ShardStrategy interface {
	// 中文：Target 声明该接口需要实现的行为。
	// English: Target declares behavior required by this interface.
	Target(key any) (ShardTarget, error)
}

// 中文：HashShardStrategy 定义当前包使用的数据结构或接口。
// English: HashShardStrategy defines a data structure or interface used by this package.
type HashShardStrategy struct {
	// 中文：Databases 保存当前结构中的配置或数据值。
	// English: Databases stores a configuration or data value for this struct.
	Databases []string
	// 中文：TablePrefix 保存当前结构中的配置或数据值。
	// English: TablePrefix stores a configuration or data value for this struct.
	TablePrefix string
	// 中文：TableCount 保存当前结构中的配置或数据值。
	// English: TableCount stores a configuration or data value for this struct.
	TableCount uint64
}

// 中文：NewHashShardStrategy 创建并返回对应组件实例。
// English: NewHashShardStrategy creates and returns the corresponding component instance.
func NewHashShardStrategy(databases []string, tablePrefix string, tableCount uint64) (*HashShardStrategy, error) {
	cleaned := make([]string, 0, len(databases))
	for _, db := range databases {
		db = strings.TrimSpace(db)
		if db != "" {
			cleaned = append(cleaned, db)
		}
	}
	if len(cleaned) == 0 {
		return nil, fmt.Errorf("at least one shard database is required")
	}
	if tableCount == 0 {
		tableCount = 1
	}
	return &HashShardStrategy{
		Databases:   cleaned,
		TablePrefix: strings.TrimSpace(tablePrefix),
		TableCount:  tableCount,
	}, nil
}

// 中文：Target 执行当前包中的对应流程。
// English: Target executes the corresponding workflow in this package.
func (s *HashShardStrategy) Target(key any) (ShardTarget, error) {
	if s == nil || len(s.Databases) == 0 {
		return ShardTarget{}, fmt.Errorf("shard strategy is not configured")
	}
	hash, err := hashShardKey(key)
	if err != nil {
		return ShardTarget{}, err
	}

	dbIndex := hash % uint64(len(s.Databases))
	tableIndex := hash % s.TableCount
	target := ShardTarget{
		Database: s.Databases[dbIndex],
		Index:    tableIndex,
	}
	if s.TablePrefix != "" {
		target.Table = fmt.Sprintf("%s_%02d", s.TablePrefix, tableIndex)
	}
	return target, nil
}

// 中文：ShardedDB 定义当前包使用的数据结构或接口。
// English: ShardedDB defines a data structure or interface used by this package.
type ShardedDB struct {
	// 中文：defaultDB 保存当前结构中的配置或数据值。
	// English: defaultDB stores a configuration or data value for this struct.
	defaultDB *DB
	// 中文：shards 保存当前结构中的配置或数据值。
	// English: shards stores a configuration or data value for this struct.
	shards map[string]*DB
	// 中文：strategy 保存当前结构中的配置或数据值。
	// English: strategy stores a configuration or data value for this struct.
	strategy ShardStrategy
}

// 中文：NewShardedDB 创建并返回对应组件实例。
// English: NewShardedDB creates and returns the corresponding component instance.
func NewShardedDB(defaultDB *DB, shards map[string]*DB, strategy ShardStrategy) *ShardedDB {
	copied := make(map[string]*DB, len(shards))
	for name, db := range shards {
		name = strings.TrimSpace(name)
		if name != "" && db != nil {
			copied[name] = db
		}
	}
	return &ShardedDB{
		defaultDB: defaultDB,
		shards:    copied,
		strategy:  strategy,
	}
}

// 中文：ForKey 执行当前包中的对应流程。
// English: ForKey executes the corresponding workflow in this package.
func (s *ShardedDB) ForKey(key any) (*DB, ShardTarget, error) {
	if s == nil || s.strategy == nil {
		return nil, ShardTarget{}, fmt.Errorf("sharded db is not configured")
	}
	target, err := s.strategy.Target(key)
	if err != nil {
		return nil, ShardTarget{}, err
	}

	db := s.defaultDB
	if target.Database != "" {
		if shard, ok := s.shards[target.Database]; ok {
			db = shard
		} else if s.defaultDB == nil {
			return nil, target, fmt.Errorf("%w: %s", ErrShardNotFound, target.Database)
		}
	}
	if db == nil {
		return nil, target, ErrShardNotFound
	}
	if target.Table == "" {
		return db, target, nil
	}
	return db.clone(db.db.Table(target.Table), db.reader().Table(target.Table)), target, nil
}

// 中文：hashShardKey 执行当前包中的对应流程。
// English: hashShardKey executes the corresponding workflow in this package.
func hashShardKey(key any) (uint64, error) {
	switch v := key.(type) {
	case nil:
		return 0, ErrShardKeyRequired
	case string:
		if strings.TrimSpace(v) == "" {
			return 0, ErrShardKeyRequired
		}
		return hashString(v), nil
	case []byte:
		if len(v) == 0 {
			return 0, ErrShardKeyRequired
		}
		return hashBytes(v), nil
	case int:
		return uint64(v), nil
	case int8:
		return uint64(v), nil
	case int16:
		return uint64(v), nil
	case int32:
		return uint64(v), nil
	case int64:
		return uint64(v), nil
	case uint:
		return uint64(v), nil
	case uint8:
		return uint64(v), nil
	case uint16:
		return uint64(v), nil
	case uint32:
		return uint64(v), nil
	case uint64:
		return v, nil
	case fmt.Stringer:
		return hashShardKey(v.String())
	default:
		return hashString(fmt.Sprint(v)), nil
	}
}

// 中文：hashString 执行当前包中的对应流程。
// English: hashString executes the corresponding workflow in this package.
func hashString(value string) uint64 {
	return hashBytes([]byte(value))
}

// 中文：hashBytes 执行当前包中的对应流程。
// English: hashBytes executes the corresponding workflow in this package.
func hashBytes(value []byte) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(value)
	return h.Sum64()
}

// 中文：ShardTable 执行当前包中的对应流程。
// English: ShardTable executes the corresponding workflow in this package.
func ShardTable(prefix string, index uint64) string {
	prefix = strings.TrimSpace(prefix)
	if prefix == "" {
		return ""
	}
	return prefix + "_" + strconv.FormatUint(index, 10)
}
