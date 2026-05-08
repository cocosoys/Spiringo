package convert

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
)

// 中文：String 执行当前包中的对应流程。
// English: String executes the corresponding workflow in this package.
func String(value any) string {
	switch v := value.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprintf("%v", value)
	}
}

// 中文：Int 执行当前包中的对应流程。
// English: Int executes the corresponding workflow in this package.
func Int(value any) (int, error) {
	v, err := Int64(value)
	return int(v), err
}

// 中文：Int64 执行当前包中的对应流程。
// English: Int64 executes the corresponding workflow in this package.
func Int64(value any) (int64, error) {
	switch v := value.(type) {
	case int:
		return int64(v), nil
	case int8:
		return int64(v), nil
	case int16:
		return int64(v), nil
	case int32:
		return int64(v), nil
	case int64:
		return v, nil
	case uint:
		return int64(v), nil
	case uint8:
		return int64(v), nil
	case uint16:
		return int64(v), nil
	case uint32:
		return int64(v), nil
	case uint64:
		if v > math.MaxInt64 {
			return 0, fmt.Errorf("uint64 overflows int64: %d", v)
		}
		return int64(v), nil
	case float32:
		return int64(v), nil
	case float64:
		return int64(v), nil
	case string:
		return strconv.ParseInt(strings.TrimSpace(v), 10, 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to int64", value)
	}
}

// 中文：Float64 执行当前包中的对应流程。
// English: Float64 executes the corresponding workflow in this package.
func Float64(value any) (float64, error) {
	switch v := value.(type) {
	case float32:
		return float64(v), nil
	case float64:
		return v, nil
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, err := Int64(v)
		return float64(i), err
	case string:
		return strconv.ParseFloat(strings.TrimSpace(v), 64)
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

// 中文：Bool 执行当前包中的对应流程。
// English: Bool executes the corresponding workflow in this package.
func Bool(value any) (bool, error) {
	switch v := value.(type) {
	case bool:
		return v, nil
	case string:
		return strconv.ParseBool(strings.TrimSpace(v))
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		i, err := Int64(v)
		return i != 0, err
	default:
		return false, fmt.Errorf("cannot convert %T to bool", value)
	}
}

// 中文：Strings 执行当前包中的对应流程。
// English: Strings executes the corresponding workflow in this package.
func Strings(value any) ([]string, error) {
	if value == nil {
		return nil, nil
	}
	switch v := value.(type) {
	case []string:
		return append([]string(nil), v...), nil
	case string:
		if strings.TrimSpace(v) == "" {
			return nil, nil
		}
		parts := strings.Split(v, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			if item := strings.TrimSpace(part); item != "" {
				result = append(result, item)
			}
		}
		return result, nil
	}

	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return nil, fmt.Errorf("cannot convert %T to []string", value)
	}
	result := make([]string, 0, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		result = append(result, String(rv.Index(i).Interface()))
	}
	return result, nil
}

// 中文：JSONMap 执行当前包中的对应流程。
// English: JSONMap executes the corresponding workflow in this package.
func JSONMap(data []byte) (map[string]any, error) {
	var out map[string]any
	if len(data) == 0 {
		return map[string]any{}, nil
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return nil, err
	}
	if out == nil {
		out = map[string]any{}
	}
	return out, nil
}
