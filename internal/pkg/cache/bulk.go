package cache

import (
	"fmt"
	"reflect"
)

// 中文：assignCacheMap 执行当前包中的对应流程。
// English: assignCacheMap executes the corresponding workflow in this package.
func assignCacheMap(dest any, values map[string]any) error {
	if dest == nil {
		return nil
	}
	target := reflect.ValueOf(dest)
	if target.Kind() != reflect.Ptr || target.IsNil() {
		return fmt.Errorf("cache destination must be a non-nil pointer")
	}
	elem := target.Elem()
	if elem.Kind() != reflect.Map {
		return assign(dest, values)
	}
	if elem.Type().Key().Kind() != reflect.String {
		return fmt.Errorf("cache destination map key must be string")
	}

	result := reflect.MakeMapWithSize(elem.Type(), len(values))
	valueType := elem.Type().Elem()
	for key, value := range values {
		mapValue, err := cacheMapValue(valueType, value)
		if err != nil {
			return fmt.Errorf("cache value for key %q: %w", key, err)
		}
		result.SetMapIndex(reflect.ValueOf(key).Convert(elem.Type().Key()), mapValue)
	}
	elem.Set(result)
	return nil
}

// 中文：cacheMapValue 执行当前包中的对应流程。
// English: cacheMapValue executes the corresponding workflow in this package.
func cacheMapValue(target reflect.Type, value any) (reflect.Value, error) {
	if value == nil {
		return reflect.Zero(target), nil
	}
	source := reflect.ValueOf(value)
	if source.Type().AssignableTo(target) {
		return source, nil
	}
	if source.Type().ConvertibleTo(target) {
		return source.Convert(target), nil
	}
	if target.Kind() == reflect.Interface && source.Type().Implements(target) {
		return source, nil
	}
	return reflect.Value{}, fmt.Errorf("type %T cannot be assigned to %s", value, target)
}
