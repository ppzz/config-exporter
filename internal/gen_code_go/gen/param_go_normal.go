package gen

import (
	"github.com/ppzz/config-exporter/internal/helper"
	"github.com/ppzz/config-exporter/internal/lib/csver"
	"path"
)

type ParamFieldGoNormal struct {
	ColIndex    int    // 列序号
	ColName     string // 列名字 原csv列名
	ColDesc     string // 列描述 原csv列描述
	ColTyp      string // 列数据类型 原csv列类型
	GoFieldName string // 属性名字 (go 文件中的属性名字)
}

type ParamGoNormal struct {
	CsvFilename        string
	CsvFilePath        string
	CsvDataRowCount    int
	CsvDataColumnCount int

	CsvMetaVarName   string // go file 中 用户暂存 csv 信息的变量名字（首字母小写，xxxCsvMeta结尾）
	ManagerClassName string // go file 类名
	ClassName        string // 根据CSV名字转化而来的go类名

	Fields []*ParamFieldGoNormal
}

func CrateParamGoNormal(csv *csver.ConfigCsv) *ParamGoNormal {
	csvFilePath := csv.Csv.FilePath
	camelCaseName := helper.NameToCamelCase(helper.NameBareName(csvFilePath))
	className := helper.NameUpperFirstLetter(camelCaseName) // 这里的 需要跟外部生成 index config 文件的类名保持一致

	param := &ParamGoNormal{
		CsvFilename:        path.Base(csvFilePath),
		CsvFilePath:        csvFilePath,
		CsvDataRowCount:    len(csv.Grid),
		CsvDataColumnCount: len(csv.Grid[0]),
		CsvMetaVarName:     camelCaseName + "CsvMeta",
		ManagerClassName:   className + "Storage",
		ClassName:          className,
		Fields:             make([]*ParamFieldGoNormal, 0),
	}

	for i, cell := range csv.HeaderLineName {
		param.Fields = append(param.Fields, &ParamFieldGoNormal{
			ColIndex:    i,
			ColName:     cell,
			ColDesc:     csv.HeaderLineDesc[i],
			ColTyp:      csv.HeaderLineType[i],
			GoFieldName: helper.NameUpperFirstLetter(cell),
		})
	}

	return param
}

// func GetNormalClassName(bareName string) string {
// 	lowerName := strings.ToLower(bareName)
// 	// 要求表名必须以字母开头
// 	if !regexp.MustCompile(`^[a-z].*`).MatchString(lowerName) {
// 		log.Fatal("name not valid. must start with letter: ", bareName)
// 	}
//
// 	if regexp.MustCompile(`^i18n_`).MatchString(lowerName) {
// 		return strings.ReplaceAll(tempName, "_", "")
// 	}
//
// 	// 正则表达式，用于匹配中文字符
// 	englishName := regexp.MustCompile(`\p{Han}`).ReplaceAllString(tempName, "") // 移除中文字符
//
// 	// 移除第一个下划线后的内容
// 	r = regexp.MustCompile(`_.*`)
// 	return r.ReplaceAllString(englishName, "")
// }
