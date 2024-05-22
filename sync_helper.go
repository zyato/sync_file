package main

import (
	"log"
	"os"
	"path/filepath"
	"time"
)

type syncHelper struct {
	sourceDir string
	targetDir string
	interval  time.Duration
	err       error
	l         *log.Logger
}

func newSyncHelper(sourceDir, targetDir string) *syncHelper {
	return &syncHelper{
		sourceDir: sourceDir,
		targetDir: targetDir,
		interval:  time.Hour,
	}
}

func (h *syncHelper) withLogger(l *log.Logger) *syncHelper {
	h.l = l
	return h
}

func (h *syncHelper) withInterval(intervalSeconds int) *syncHelper {
	h.interval = time.Duration(intervalSeconds) * time.Second
	return h
}

func (h *syncHelper) syncFile() *syncHelper {
	if h.err != nil {
		return h
	}
	h.err = filepath.Walk(h.sourceDir, h.processFile)
	return h
}

func (h *syncHelper) processFile(path string, fInfo os.FileInfo, err error) error {
	if err != nil {
		return wrap(err, "filepath.Walk fn")
	}
	relativePath, err := filepath.Rel(h.sourceDir, path)
	if err != nil {
		return wrap(err, "Rel")
	}
	targetPath := filepath.Join(h.targetDir, relativePath)
	// 创建target下对应的的目录路径
	if fInfo.IsDir() {
		return wrap(os.MkdirAll(targetPath, os.ModePerm), "Mkdir")
	}
	// 目标文件不存在，直接同步
	dstInfo, err := os.Stat(targetPath)
	if os.IsNotExist(err) {
		err = copyFile(path, targetPath)
		if err != nil {
			h.l.Printf("同步新文件错误 %s -> %s\n", path, targetPath)
			return wrap(err, "copyFile")
		}
		h.l.Printf("成功同步新文件 %s -> %s\n", path, targetPath)
		return nil
	}
	// 目标文件存在，判断是否修改
	srcInfo, err := os.Stat(path)
	if err != nil {
		return wrap(err, "Stat(path)")
	}
	if !srcInfo.ModTime().After(dstInfo.ModTime()) {
		return nil
	}
	// 目标文件被修改过，判断md5是否变动
	srcMD5, err := getFileMD5(path)
	if err != nil {
		return wrap(err, "getFileMD5(path)")
	}
	dstMD5, err := getFileMD5(targetPath)
	if err != nil {
		return wrap(err, "getFileMD5(targetPath)")
	}
	if srcMD5 == dstMD5 {
		return nil
	}
	err = copyFile(path, targetPath)
	if err != nil {
		h.l.Printf("同步改动文件错误 %s -> %s\n", path, targetPath)
		return wrap(err, "copyFile")
	}
	h.l.Printf("成功同步改动文件 %s -> %s\n", path, targetPath)
	return nil
}

func (h *syncHelper) run() error {
	if h.err != nil {
		return h.err
	}
	if h.l == nil {
		h.l = &log.Logger{}
	}
	for {
		h.l.Println("start")
		h.syncFile()
		if h.err != nil {
			h.l.Println("出现错误, err", h.err)
		}
		h.l.Println("end")
		time.Sleep(h.interval)
	}
}
