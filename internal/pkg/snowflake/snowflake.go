package snowflake

import (
	"fmt"
	"sync"
	"time"
)

// 中文：Snowflake 定义当前包使用的数据结构或接口。
// English: Snowflake defines a data structure or interface used by this package.
// Snowflake ID生成器（Twitter Snowflake算法简化实现）
// 64位ID = 1位符号位 + 41位时间戳 + 10位机器ID + 12位序列号
type Snowflake struct {
	// 中文：mu 保存当前结构中的配置或数据值。
	// English: mu stores a configuration or data value for this struct.
	mu sync.Mutex
	// 中文：epoch 保存当前结构中的配置或数据值。
	// English: epoch stores a configuration or data value for this struct.
	epoch int64
	// 中文：machineID 保存当前结构中的配置或数据值。
	// English: machineID stores a configuration or data value for this struct.
	machineID int64
	// 中文：sequence 保存当前结构中的配置或数据值。
	// English: sequence stores a configuration or data value for this struct.
	sequence int64
	// 中文：lastTime 保存当前结构中的配置或数据值。
	// English: lastTime stores a configuration or data value for this struct.
	lastTime int64
}

// 中文：machineIDBits、sequenceBits、maxMachineID、... 声明当前包使用的常量。
// English: machineIDBits、sequenceBits、maxMachineID、... declares constants used by this package.
const (
	machineIDBits = 10
	sequenceBits  = 12

	maxMachineID = -1 ^ (-1 << machineIDBits) // 1023
	maxSequence  = -1 ^ (-1 << sequenceBits)  // 4095

	machineIDShift = sequenceBits
	timestampShift = sequenceBits + machineIDBits
)

// 中文：New 创建并返回对应组件实例。
// English: New creates and returns the corresponding component instance.
// New 创建Snowflake ID生成器
func New(machineID int64) (*Snowflake, error) {
	if machineID < 0 || machineID > maxMachineID {
		return nil, fmt.Errorf("machine ID must be between 0 and %d", maxMachineID)
	}

	// 使用2024-01-01作为纪元
	epoch := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).UnixMilli()

	return &Snowflake{
		epoch:     epoch,
		machineID: machineID,
	}, nil
}

// 中文：Generate 执行当前包中的对应流程。
// English: Generate executes the corresponding workflow in this package.
// Generate 生成唯一ID
func (s *Snowflake) Generate() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixMilli() - s.epoch

	if now == s.lastTime {
		s.sequence = (s.sequence + 1) & maxSequence
		if s.sequence == 0 {
			// 等待下一毫秒
			for now <= s.lastTime {
				now = time.Now().UnixMilli() - s.epoch
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastTime = now

	return (now << timestampShift) | (s.machineID << machineIDShift) | s.sequence
}

// 中文：GenerateString 执行当前包中的对应流程。
// English: GenerateString executes the corresponding workflow in this package.
// GenerateString 生成字符串ID
func (s *Snowflake) GenerateString() string {
	return fmt.Sprintf("%d", s.Generate())
}
