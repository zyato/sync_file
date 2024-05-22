package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

func wrap(err error, s string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", s, err)
}

// 计算文件MD5
func getFileMD5(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", wrap(err, "Open")
	}
	defer file.Close()

	hash := md5.New()
	if _, err = io.Copy(hash, file); err != nil {
		return "", wrap(err, "Copy")
	}

	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

// 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return wrap(err, "Open")
	}
	defer sourceFile.Close()
	destFile, err := os.Create(dst)
	if err != nil {
		return wrap(err, "Create")
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return wrap(err, "Copy")
	}

	// 确保数据写入到磁盘
	err = destFile.Sync()
	if err != nil {
		return wrap(err, "Sync")
	}

	return nil
}
