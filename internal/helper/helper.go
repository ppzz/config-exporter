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

func DirMustEmpty(dir string) {
	// 使用Stat获取文件或目录的信息
	fileInfo, err := os.Stat(dir)
	if os.IsNotExist(err) {
		return
	}
	if err != nil {
		log.Fatal("stat failed, dir: ", dir, err.Error())
	}

	// 检查是否为目录
	if !fileInfo.IsDir() {
		log.Fatal("not a dir: ", dir)
	}

	// 打开目录
	opened, err := os.Open(dir)
	if err != nil {
		log.Fatal("open dir failed, dir: ", dir, err.Error())
	}
	defer opened.Close()

	// 读取目录内容
	_, err = opened.Readdirnames(1) // 尝试读取至少一个条目
	if errors.Is(err, io.EOF) || errors.Is(err, os.ErrNotExist) {
		return
	}
	if err != nil {
		log.Fatal("read dir failed, dir: ", dir, err.Error())
	}

	log.Fatal("dir not empty: ", dir)
}

func NameBareName(fullPath string) string {
	basename := path.Base(fullPath)
	ext := path.Ext(basename)
	return basename[:len(basename)-len(ext)]
}

func BareNameOfCsvFile(fullPath string) string {
	bareName := NameBareName(fullPath)
	bareName = strings.Split(bareName, ".0.")[0]
	bareName = NameRemoveHan(bareName)
	idx := strings.LastIndex(bareName, "_")
	if idx == -1 {
		log.Fatal("invalid csv file name, not found '_' :", fullPath)
	}
	return bareName[:idx]
}

func FmtCsvFileNameOfCsvFile(fullPath string) string {
	return BareNameOfCsvFile(fullPath) + ".csv"
}
