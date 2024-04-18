package starter

import (
	"embed"
	"github.com/ppzz/config-exporter/internal/gen_code_go/gen"
	"github.com/ppzz/config-exporter/internal/gen_code_go/setting"
	"github.com/ppzz/config-exporter/internal/helper"
	"github.com/ppzz/config-exporter/internal/lib/csver"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"path"
	"sort"
	"text/template"
)

//go:embed template
var tmplFs embed.FS

func Start() {
	csvDir := setting.Get().InputFmtCsvDir
	codeDir := setting.Get().OutputCodeGoDir

	csvFilenameList := helper.ListFilenameByExt(csvDir, ".csv")
	csvFilePathList := lo.Map(csvFilenameList, func(item string, index int) string {
		return path.Join(csvDir, item)
	})

	// 读取csv文件, 返回一组 ConfigCsv 对象
	csvList := lo.Map(csvFilePathList, func(item string, index int) *csver.ConfigCsv {
		grid := helper.FileCsvRead(item)
		csv := csver.NewCsv(item, grid)
		return csver.CreateConfigCsv(csv)
	})

	// 分类,normal, global
	typ2List := lo.GroupBy(csvList, func(item *csver.ConfigCsv) helper.ConfigTyp {
		return helper.GetConfType(path.Base(item.Csv.FilePath))
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
	bareName := helper.NameBareName(csv.Csv.FilePath)
	goFilePath := path.Join(codeDir, "a_"+helper.NameToSnakeCase(bareName)+".autogen.go")

	tmpl, err := template.ParseFS(tmplFs, "template/global.tmpl")
	cobra.CheckErr(err)

	param := gen.CreateParamGoGlobal(csv)
	helper.RenderTemplate(goFilePath, tmpl, param)
}

// exportNormalGoFile 导出普通的go文件
func exportNormalGoFile(codeDir string, csv *csver.ConfigCsv) {
	bareName := helper.NameBareName(csv.Csv.FilePath)
	goFilePath := path.Join(codeDir, helper.NameToSnakeCase(bareName)+".autogen.go")

	tmpl, err := template.ParseFS(tmplFs, "template/normal.tmpl")
	cobra.CheckErr(err)

	param := gen.CrateParamGoNormal(csv)
	helper.RenderTemplate(goFilePath, tmpl, param)
}

func exportNormalIndexGoFile(codeDir string, list []*csver.ConfigCsv) {
	goFilePath := path.Join(codeDir, "a_index.autogen.go")

	tmpl, err := template.ParseFS(tmplFs, "template/normal_index.tmpl")
	cobra.CheckErr(err)

	classNameList := lo.Map(list, func(item *csver.ConfigCsv, index int) string {
		// 这里的算法需要跟 生成 normal config 文件中的类名保持一致
		camelCaseName := helper.NameToCamelCase(helper.NameBareName(item.Csv.FilePath))
		return helper.NameUpperFirstLetter(camelCaseName)
	})
	sort.Strings(classNameList)

	param := map[string]any{"List": classNameList}
	helper.RenderTemplate(goFilePath, tmpl, param)
}
