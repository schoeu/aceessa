package utils

import (
	"os"
	"log"
	"time"
	"strings"
	"fmt"
)

// 获取程序cwd
func GetCwd() string {
	dir, err := os.Getwd()
	ErrHadle(err)
	return dir
}

// 错误信息梳理
func ErrHadle(err interface{}) {
	if err != nil {
		log.Println(err)
	}
}

// 清除临时文件&文件夹
func CleanTmp(p string) {
	if p == "" {
		return
	}
	err := os.RemoveAll(p)
	ErrHadle(err)
}

// 获取当前时间字符串
func GetDateString(date time.Time) string {
	t := date.String()
	return strings.Replace(strings.Split(t, " ")[0], "-", "", -1)
}

// 创建临时文件夹存放中间文件
func EnsureDir(path string) {
	_, err := os.Stat(path)
	if err != nil {
		fmt.Println("path not exists ", path)
		err := os.MkdirAll(path, 0711)

		if err != nil {
			log.Fatal(err)
		}
	}
}
