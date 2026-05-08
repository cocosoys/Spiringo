package oss

import (
	"github.com/spiringo/spiringo/internal/pkg/storage"
)

// 中文：OSS 定义当前包使用的数据结构或接口。
// English: OSS defines a data structure or interface used by this package.
type OSS = storage.OSS

// 中文：Storage 定义当前包使用的数据结构或接口。
// English: Storage defines a data structure or interface used by this package.
type Storage = storage.Storage

// 中文：BucketStorage 定义当前包使用的数据结构或接口。
// English: BucketStorage defines a data structure or interface used by this package.
type BucketStorage = storage.BucketStorage

// 中文：MinIOStorage 定义当前包使用的数据结构或接口。
// English: MinIOStorage defines a data structure or interface used by this package.
type MinIOStorage = storage.MinIOStorage

// 中文：CephStorage 定义当前包使用的数据结构或接口。
// English: CephStorage defines a data structure or interface used by this package.
type CephStorage = storage.CephStorage

// 中文：MinIOConfig 定义当前包使用的数据结构或接口。
// English: MinIOConfig defines a data structure or interface used by this package.
type MinIOConfig = storage.MinIOConfig

// 中文：CephConfig 定义当前包使用的数据结构或接口。
// English: CephConfig defines a data structure or interface used by this package.
type CephConfig = storage.CephConfig

// 中文：NewBucketStorage 声明当前包使用的变量。
// English: NewBucketStorage declares variables used by this package.
var NewBucketStorage = storage.NewBucketStorage

// 中文：NewOSS 声明当前包使用的变量。
// English: NewOSS declares variables used by this package.
var NewOSS = storage.NewOSS

// 中文：NewMinIOStorage 声明当前包使用的变量。
// English: NewMinIOStorage declares variables used by this package.
var NewMinIOStorage = storage.NewMinIOStorage

// 中文：NewCephStorage 声明当前包使用的变量。
// English: NewCephStorage declares variables used by this package.
var NewCephStorage = storage.NewCephStorage
