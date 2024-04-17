package starter

import (
	"github.com/ppzz/config-exporter/internal/gen_code_go/gen"
	"github.com/ppzz/config-exporter/internal/gen_code_go/setting"
	"github.com/ppzz/config-exporter/internal/helper"
	"github.com/ppzz/config-exporter/internal/lib/csver"
	"github.com/samber/lo"
	"path"
)

func Start() {
	csvDir := setting.Get().InputFmtCsvDir
	codeDir := setting.Get().OutputCodeGoDir

	csvFilenameList := helper.ListFilenameByExt(csvDir, ".csv")
	csvFilenameList = lo.Map(csvFilenameList, func(item string, index int) string {
		return path.Join(csvDir, item)
	})

	// 读取csv文件, 返回一组 ConfigCsv 对象
	csvList := lo.Map(csvFilenameList, func(item string, index int) *csver.ConfigCsv {
		grid := helper.FileCsvRead(item)
		csv := csver.NewCsv(item, grid)
		return csver.CreateConfigCsv(csv)
	})

	// 分类,normal, global
	typ2List := lo.GroupBy(csvList, func(item *csver.ConfigCsv) helper.ConfigTyp {
		return helper.GetConfType(helper.FileBareName(item.Csv.FilePath))
	})

	normalCsvList := typ2List[helper.ConfigTypNormal] // 普通配置
	globalCsvList := typ2List[helper.ConfigTypGlobal] // 全局配置

	helper.DirMustEmpty(codeDir)
	helper.MakeSureExist(codeDir)

	// go normal config
	lo.ForEach(normalCsvList, func(item *csver.ConfigCsv, index int) {
		exportNormalGoFile(codeDir, item)
	})

	// go global config
	lo.ForEach(globalCsvList, func(item *csver.ConfigCsv, index int) {
		exportGlobalGoFile(codeDir, item)
	})

	// go normal index config
	exportNormalIndexGoFile(codeDir, normalCsvList)
}

func exportGlobalGoFile(codeDir string, csv *csver.ConfigCsv) {
	bareName := helper.FileBareName(csv.Csv.FilePath)
	goFilePath := path.Join(codeDir, "a_"+helper.CamelCaseToSnakeCase(bareName)+".autogen.go")

	templateDir := "internal/gen_code_go/template"
	TemplatePath := path.Join(templateDir, "global.tmpl")

	param := gen.CreateParamGoGlobal(csv)
	helper.RenderTemplate(goFilePath, TemplatePath, param)

}

// exportNormalGoFile 导出普通的go文件
func exportNormalGoFile(codeDir string, csv *csver.ConfigCsv) {
	bareName := helper.FileBareName(csv.Csv.FilePath)
	goFilePath := path.Join(codeDir, helper.CamelCaseToSnakeCase(bareName)+".autogen.go")

	templateDir := "internal/gen_code_go/template"
	TemplatePath := path.Join(templateDir, "normal.tmpl")

	param := gen.CrateParamGoNormal(csv)
	helper.RenderTemplate(goFilePath, TemplatePath, param)
}

func exportNormalIndexGoFile(codeDir string, list []*csver.ConfigCsv) {
	goFilePath := path.Join(codeDir, "a_index.autogen.go")

	templateDir := "internal/gen_code_go/template"
	TemplatePath := path.Join(templateDir, "normal_index.tmpl")

	classNameList := lo.Map(list, func(item *csver.ConfigCsv, index int) string {
		csvFilePath := item.Csv.FilePath
		bareName := helper.FileBareName(csvFilePath)
		snakeCaseName := helper.CamelCaseToSnakeCase(bareName)
		camelCaseName := helper.SnakeToCamel(snakeCaseName)
		return helper.UpperFirstLetter(camelCaseName)
	})
	param := map[string]any{"List": classNameList}
	helper.RenderTemplate(goFilePath, TemplatePath, param)
}
