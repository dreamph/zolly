package utils

import (
	"os"
	"path/filepath"
)

func IsEmpty(value string) bool {
	return value == ""
}

func IsNotEmpty(value string) bool {
	return !IsEmpty(value)
}

func IsEmptyList[T any](value *[]T) bool {
	return value == nil || len(*value) == 0
}

func IsNotEmptyList[T any](value *[]T) bool {
	return !IsEmptyList(value)
}

func FileExists(name string) bool {
	_, err := os.Stat(name)
	return err == nil
}

func WriteFile(filePath string, data []byte) error {
	if _, err := os.Stat(filepath.Dir(filePath)); os.IsNotExist(err) {
		err := os.MkdirAll(filepath.Dir(filePath), 0700)
		if err != nil {
			return err
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
