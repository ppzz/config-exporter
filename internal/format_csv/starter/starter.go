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

	configCsvList := lo.Map(csvFilenameList, func(item string, index int) *csver.ConfigCsv {
		grids := helper.FileCsvRead(item)
		csv := csver.NewCsv(item, grids)
		return csver.CreateConfigCsv(csv)
	})

	configCsvList = lo.Without(configCsvList, nil)
	if len(configCsvList) == 0 {
		return
	}

	helper.DirMustEmpty(fmtDir)
	helper.MakeSureExist(fmtDir)
	lo.ForEach(configCsvList, func(item *csver.ConfigCsv, index int) {
		bareName := helper.FileBareName(item.Csv.FilePath)
		bareName = helper.CamelCaseToSnakeCase(helper.FilenameByType(bareName))
		fmtCsvFilePath := path.Join(fmtDir, bareName+".csv")
		helper.FileCsvWrite(fmtCsvFilePath, item.ToGrid())
	})
}
