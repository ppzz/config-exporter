package starter

import (
	"github.com/ppzz/golang-csv/internal/excel_to_csv/excel"
	"github.com/ppzz/golang-csv/internal/excel_to_csv/setting"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func Start() {
	ExcelDir := setting.Get().InputExcelDir
	CsvDir := setting.Get().OutputCsvDir

	emptyDir, err := IsEmptyDir(CsvDir)
	if err != nil {
		log.Fatal("check csv dir failed: ", err.Error())
	}
	if !emptyDir {
		log.Fatal("output dir is not empty")
	}

	filenames := getExcelNames(ExcelDir)

	h := excel.NewExcelHandler()

	excelList := lo.Map(filenames, func(item string, index int) *excel.Excel {
		ex, err := h.ReadOne(item)
		if err != nil {
			log.Fatal("read excel file failed: ", item, err.Error())
		}
		return ex
	})

	// 导出csv
	for _, item := range excelList {
		h.ExportToCsv(CsvDir, item)
	}
}

// IsEmptyDir 检查指定的路径是否是一个空目录
func IsEmptyDir(path string) (bool, error) {
	// 使用Stat获取文件或目录的信息
	fileInfo, err := os.Stat(path)
	if err != nil {
		// 如果错误类型是os.ErrNotExist，说明文件或目录不存在
		if os.IsNotExist(err) {
			return true, nil // 路径不存在，认为是空目录
		}
		return false, err // 其他错误，返回错误信息
	}

	// 检查是否为目录
	if !fileInfo.IsDir() {
		return false, nil // 如果不是目录（即是文件），返回false
	}

	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		return false, err // 打开目录时出错，返回错误信息
	}
	defer dir.Close()

	// 读取目录内容
	_, err = dir.Readdirnames(1) // 尝试读取至少一个条目
	if err == io.EOF {
		return true, nil
	}
	if err != nil {
		// 如果读到文件末尾，说明目录为空
		if err == os.ErrNotExist {
			return true, nil
		}
		return false, err // 读取目录出错，返回错误信息
	}

	// 如果能够读取到至少一个文件或目录，那么目录不为空
	return false, nil
}

// getExcelNames 返回 dir 下的 excel 文件名列表
func getExcelNames(dir string) []string {
	list := ListFilenameByExt(dir, setting.ConstExcelExtNames...)
	return lo.Map(list, func(item string, index int) string {
		return path.Join(dir, item)
	})
}

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
