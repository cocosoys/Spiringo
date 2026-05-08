package utils

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

// 中文：FileExists 执行当前包中的对应流程。
// English: FileExists executes the corresponding workflow in this package.
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 中文：DirExists 执行当前包中的对应流程。
// English: DirExists executes the corresponding workflow in this package.
func DirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// 中文：EnsureDir 执行当前包中的对应流程。
// English: EnsureDir executes the corresponding workflow in this package.
func EnsureDir(path string) error {
	if path == "" {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}

// 中文：ReadJSON 执行当前包中的对应流程。
// English: ReadJSON executes the corresponding workflow in this package.
func ReadJSON(path string, target any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, target)
}

// 中文：WriteJSON 执行当前包中的对应流程。
// English: WriteJSON executes the corresponding workflow in this package.
func WriteJSON(path string, value any) error {
	if err := EnsureDir(filepath.Dir(path)); err != nil {
		return err
	}
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

// 中文：RemoveIfExists 执行当前包中的对应流程。
// English: RemoveIfExists executes the corresponding workflow in this package.
func RemoveIfExists(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}
