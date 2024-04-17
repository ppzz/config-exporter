package starter

import (
	"github.com/ppzz/golang-csv/internal/format_csv/setting"
	"github.com/ppzz/golang-csv/internal/helper"
	"github.com/ppzz/golang-csv/internal/lib/csver"
	"github.com/samber/lo"
	"path"
	"regexp"
)

func Start() {
	csvDir := setting.Get().InputCsvDir
	fmtDir := setting.Get().OutputFmtDir
	csvFilenameSchema := setting.Get().CsvFileNameSchema

	csvFilenameList := helper.ListFilenameByExt(csvDir, ".csv")
	csvFilenameList = lo.Filter(csvFilenameList, func(item string, index int) bool {
		return regexp.MustCompile(csvFilenameSchema).MatchString(item)
	})
	csvFilenameList = lo.Map(csvFilenameList, func(item string, index int) string {
		return path.Join(csvDir, item)
	})

	csvList := lo.Map(csvFilenameList, func(item string, index int) *csver.Csv {
		grids := helper.FileCsvRead(item)
		return csver.NewCsv(item, grids)
	})

	configCsvList := lo.Map(csvList, func(item *csver.Csv, index int) *csver.ConfigCsv {
		return csver.CreateConfigCsv(item)
	})

	configCsvList = lo.Without(configCsvList, nil)
	if len(configCsvList) == 0 {
		return
	}

	helper.DirMustEmpty(fmtDir)
	helper.MakeSureExist(fmtDir)
	lo.ForEach(configCsvList, func(item *csver.ConfigCsv, index int) {
		bareName := helper.FileBareName(item.Csv.FilePath)
		n := helper.CamelCaseToSnakeCase(helper.FilenameByType(bareName))

		newCsvFilePath := path.Join(fmtDir, n+".csv")
		helper.FileCsvWrite(item.ToGrid(), newCsvFilePath)
	})

	//  0. csv 文件过滤
	// - 按照cs标记过滤列
	// - 按照 type 过滤列, 如果有不符合的列，报错
	// - type 映射到通用的 type
	// - name 补全,格式化

	// header line: name, type, desc, c-s-flag
}
