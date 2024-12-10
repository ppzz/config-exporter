package starter

import (
    "path"
    "regexp"

    "github.com/samber/lo"

    "github.com/ppzz/config-exporter/internal/format_csv/setting"
    "github.com/ppzz/config-exporter/internal/helper"
    "github.com/ppzz/config-exporter/internal/lib/csver"
)

func Start() {
    csvDir := setting.Get().InputCsvDir
    fmtDir := setting.Get().OutputFmtDir
    csvFilenameSchema := setting.Get().CsvFileNameSchema

    csvFilenameList := helper.ListFilenameByExt(csvDir, ".csv")
    csvFilenameList = lo.Filter(csvFilenameList, func(item string, index int) bool {
        return regexp.MustCompile(csvFilenameSchema).MatchString(item)
    })
    csvFilePathList := lo.Map(csvFilenameList, func(item string, index int) string {
        return path.Join(csvDir, item)
    })

    configCsvList := lo.Map(csvFilePathList, func(item string, index int) *csver.ConfigCsv {
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
        fmtCsvFilePath := path.Join(fmtDir, helper.FmtCsvFileNameOfCsvFile(item.Csv.FilePath))
        helper.FileCsvWrite(fmtCsvFilePath, item.ToGrid())
    })
}
