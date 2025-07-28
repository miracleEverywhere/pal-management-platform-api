package utils

import (
	"fmt"
	"os"
	"strings"
)

// EnsureFileExists 检查文件是否存在，如果不存在则创建空文件
func EnsureFileExists(filePath string) error {
	// 检查文件是否存在
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		// 文件不存在，创建一个空文件
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		err = file.Close()
		if err != nil {
			return err
		}
	} else if err != nil {
		// 其他错误
		return err
	}

	return nil
}

// EnsureDirExists 检查目录是否存在，如果不存在则创建
func EnsureDirExists(dirPath string) error {
	if strings.HasPrefix(dirPath, "~") {
		dirPath = strings.Replace(dirPath, "~", HomeDir, 1)
	}
	// 检查目录是否存在
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		// 目录不存在，创建目录
		err := os.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			return fmt.Errorf("无法创建目录: %w", err)
		}
	} else if err != nil {
		// 其他错误
		return fmt.Errorf("检查目录时出错: %w", err)
	}

	return nil
}
