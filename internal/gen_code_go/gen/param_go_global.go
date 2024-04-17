package gen

import (
	"github.com/ppzz/config-exporter/internal/helper"
	"github.com/ppzz/config-exporter/internal/lib/csver"
	"github.com/samber/lo"
)

type ParamFieldGoGlobal struct {
	ColIndex int // 列序号

	OriginalName string // 原名字
	OriginalType string // 原类型
	OriginalVal  string // 原值
	OriginalDesc string // 原描述

	VariableName  string // go 变量名
	VariableType  string // go 变量类型
	VariableValue string // go 变量的值
}

type ParamGoGlobal struct {
	CsvFileBareName    string
	CsvFileFullPath    string
	CsvDataRowCount    int
	CsvDataColumnCount int

	CsvMetaVarName string // go file 中 用户暂存 csv 信息的变量名字（首字母小写，xxxCsvMeta结尾）

	Variables []*ParamFieldGoGlobal
}

func CreateParamGoGlobal(csv *csver.ConfigCsv) *ParamGoGlobal {
	csvFilePath := csv.Csv.FilePath
	bareName := helper.FileBareName(csvFilePath)
	snakeCaseName := helper.CamelCaseToSnakeCase(bareName)
	camelCaseName := helper.SnakeToCamel(snakeCaseName)

	param := &ParamGoGlobal{
		CsvFileBareName:    bareName,
		CsvFileFullPath:    csvFilePath,
		CsvDataRowCount:    len(csv.Grid),
		CsvDataColumnCount: len(csv.Grid[0]),
		CsvMetaVarName:     camelCaseName + "CsvMeta",
		Variables:          make([]*ParamFieldGoGlobal, 0),
	}

	const ColIdxName = 0 // 名字 列
	const ColIdxDesc = 1 // 描述 列
	const ColIdxType = 2 // 类型 列
	const ColIdxVal = 3  // 值 列

	param.Variables = lo.Map(csv.Grid, func(item []string, index int) *ParamFieldGoGlobal {
		name, desc, typ, val := item[ColIdxName], item[ColIdxDesc], item[ColIdxType], item[ColIdxVal]
		variableType := helper.FormatTypeName(typ)
		return &ParamFieldGoGlobal{
			ColIndex:      index,
			OriginalName:  name,
			OriginalType:  typ,
			OriginalVal:   val,
			OriginalDesc:  desc,
			VariableName:  "Global" + helper.UpperFirstLetter(name),
			VariableType:  variableType,
			VariableValue: helper.GetVariableLiteralValue(variableType, val),
		}
	})
	return param
}
