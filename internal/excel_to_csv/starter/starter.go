package starter

import (
	"github.com/ppzz/config-exporter/internal/excel_to_csv/setting"
	"github.com/ppzz/config-exporter/internal/helper"
	"github.com/ppzz/config-exporter/internal/lib/excel"
	"github.com/samber/lo"
	"log"
	"path"
)

func Start() {
	excelDir := setting.Get().InputExcelDir
	csvDir := setting.Get().OutputCsvDir

	emptyDir := helper.IsEmptyDir(csvDir)
	if !emptyDir {
		log.Fatal("output dir is not empty")
	}

	filenames := getExcelNames(excelDir)

	h := excel.NewExcelHandler()

	excelList := lo.Map(filenames, func(item string, index int) *excel.Excel {
		ex, err := h.ReadOne(item)
		if err != nil {
			log.Fatal("read excel file failed: ", item, err.Error())
		}
		return ex
	})

	helper.DirMustEmpty(csvDir)
	helper.MakeSureExist(csvDir)

	// 导出csv
	for _, item := range excelList {
		h.ExportToCsv(csvDir, item)
	}
}

// getExcelNames 返回 dir 下的 excel 文件名列表
func getExcelNames(dir string) []string {
	list := helper.ListFilenameByExt(dir, setting.ConstExcelExtNames...)
	return lo.Map(list, func(item string, index int) string {
		return path.Join(dir, item)
	})
}
