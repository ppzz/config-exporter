package starter

import (
	"github.com/ppzz/config-exporter/internal/excel_to_csv/setting"
	"github.com/ppzz/config-exporter/internal/helper"
	"github.com/ppzz/config-exporter/internal/lib/excel"
	"github.com/samber/lo"
	"log"
	"path"
	"regexp"
)

func Start() {
	excelDir := setting.Get().InputExcelDir
	csvDir := setting.Get().OutputCsvDir
	excelFilenameSchema := setting.Get().ExcelFileNameSchema

	helper.DirMustEmpty(csvDir)

	excelFilename := helper.ListFilenameByExt(excelDir, []string{".xlsx", ".xls"}...)

	excelFilename = lo.Filter(excelFilename, func(item string, index int) bool {
		return regexp.MustCompile(excelFilenameSchema).MatchString(item)
	})

	excelFilePath := lo.Map(excelFilename, func(item string, index int) string {
		return path.Join(excelDir, item)
	})

	h := excel.NewExcelHandler()

	excelList := lo.Map(excelFilePath, func(item string, index int) *excel.Excel {
		ex, err := h.ReadOne(item)
		if err != nil {
			log.Fatal("read excel file failed: ", item, " ", err.Error())
		}
		return ex
	})

	helper.MakeSureExist(csvDir)

	// 导出csv
	for _, item := range excelList {
		h.ExportToCsv(csvDir, item)
	}
}
