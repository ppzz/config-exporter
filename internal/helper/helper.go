package helper

import (
	"errors"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ListFilenameByExt 根据扩展名取文件, 扩展名不区分大小写
func ListFilenameByExt(dir string, extNames ...string) []string {
	if !IsExist(dir) { // 如果路径不存在
		return []string{}
	}

	infos, err := os.ReadDir(dir)
	if err != nil {
		return []string{}
	}

	extNames = lo.Map(extNames, func(item string, index int) string { return strings.ToLower(item) })

	return lo.FilterMap(infos, func(item os.DirEntry, index int) (string, bool) {
		itemExt := strings.ToLower(filepath.Ext(item.Name()))
		return item.Name(), !item.IsDir() && lo.Contains(extNames, itemExt)
	})
}

// IsExist filepath 对应的文件是否存在
func IsExist(filepath string) bool {
	_, err := os.Stat(filepath)
	// 如果 err == nil ，则存在
	return err == nil
}

func MakeSureExist(dir string) {
	// 判断目录是否存在
	exist := IsExist(dir)
	if exist {
		return
	}
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		panic(err)
	}
}

// IsEmptyDir 检查指定的路径是否是一个空目录
func IsEmptyDir(path string) bool {

	// 使用Stat获取文件或目录的信息
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true
	}
	if err != nil {
		log.Fatal("stat failed, path: ", path, err.Error())
		return false // 其他错误，返回错误信息
	}

	// 检查是否为目录
	if !fileInfo.IsDir() {
		return false
	}

	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		return false // 打开目录时出错，返回错误信息
	}
	defer dir.Close()

	// 读取目录内容
	_, err = dir.Readdirnames(1) // 尝试读取至少一个条目
	if errors.Is(err, io.EOF) || errors.Is(err, os.ErrNotExist) {
		return true
	}
	if err != nil {
		log.Fatal("read dir failed, path: ", path, err.Error())
		return false // 读取目录出错，返回错误信息
	}

	// 如果能够读取到至少一个文件或目录，那么目录不为空
	return false
}

func FileBareName(fullPath string) string {
	basename := path.Base(fullPath)
	ext := path.Ext(basename)
	return basename[:len(basename)-len(ext)]
}
