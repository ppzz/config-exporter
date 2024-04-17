package starter

import (
	"github.com/ppzz/golang-csv/internal/gen_code_go/gen"
	"github.com/ppzz/golang-csv/internal/gen_code_go/setting"
	"github.com/ppzz/golang-csv/internal/helper"
	"github.com/ppzz/golang-csv/internal/lib/csver"
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

	// fmt.Println("")
	// for i, item := range csvFilenameList {
	// 	fmt.Println(i, item)
	// }
	// fmt.Println("")

	csvList := lo.Map(csvFilenameList, func(item string, index int) *csver.ConfigCsv {
		grid := helper.FileCsvRead(item)
		return csver.CreateConfigCsv(csver.NewCsv(item, grid))
	})

	// fmt.Println("")
	// for i, item := range csvList {
	// 	fmt.Println(i, item == nil)
	// }
	// fmt.Println("")

	// 分类,normal, global
	typ2List := lo.GroupBy(csvList, func(item *csver.ConfigCsv) helper.ConfigTyp {
		return helper.GetConfType(helper.FileBareName(item.Csv.FilePath))
	})

	normalCsvList := typ2List[helper.ConfigTypNormal]
	globalCsvList := typ2List[helper.ConfigTypGlobal]

	helper.DirMustEmpty(codeDir)
	helper.MakeSureExist(codeDir)

	// go normal
	lo.ForEach(normalCsvList, func(item *csver.ConfigCsv, index int) {
		exportNormalGoFile(codeDir, item)
	})

	// go config
	lo.ForEach(globalCsvList, func(item *csver.ConfigCsv, index int) {
		exportGlobalGoFile(codeDir, item)
	})

	// go normal index
	exportNormalIndexGoFile(codeDir, normalCsvList)

	// config index

	// globalCsvList := typ2List[helper.ConfigTypGlobal]

	//  0. csv 文件过滤
	// - 按照cs标记过滤列
	// - 按照 type 过滤列, 如果有不符合的列，报错
	// - type 映射到通用的 type
	// - name 补全,格式化
	// header line: name, type, desc, c-s-flag
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
