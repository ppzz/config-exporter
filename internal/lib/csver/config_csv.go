package csver

import (
    "strconv"
    "strings"

    "github.com/samber/lo"

    "github.com/ppzz/config-exporter/internal/helper"
)

type ConfigCsv struct {
    Csv            *Csv
    HeaderLineName []string   // csv header line name
    HeaderLineDesc []string   // csv header line desc
    HeaderLineType []string   // csv header line type
    HeaderLineFlag []string   // csv header line flag
    Grid           [][]string // csv data grid
}

func (c *ConfigCsv) ToGrid() [][]string {
    headerGrid := [][]string{
        c.HeaderLineName,
        c.HeaderLineDesc,
        c.HeaderLineType,
        c.HeaderLineFlag,
    }
    return append(headerGrid, c.Grid...)
}

func CreateConfigCsv(item *Csv) *ConfigCsv {
    // - csv 文件名过滤
    // - 按照 cs 标记过滤列
    // - 按照 type 检查, 如果有不符合的列，报错
    // - type 映射到通用的 type
    // - name 补全, 格式化

    const HeaderLineCount = 4        // header 行数
    const HeaderLineNameRowIndex = 0 // header 行索引 name
    const HeaderLineDescRowIndex = 1 // header 行索引 desc
    const HeaderLineTypeRowIndex = 2 // header 行索引 type
    const HeaderLineFlagRowIndex = 3 // header 行索引 flag

    grid := helper.GridTrimCellSpace(item.Grid) // 去掉每个cell的空格
    grid = helper.GridOmitEmptyRow(grid)        // 删除空行

    headerGrid, dataGrid := helper.GridSplitHeader(grid, HeaderLineCount) // 拆分 header 和 data
    if len(headerGrid) != HeaderLineCount || len(dataGrid) == 0 {
        return nil
    }

    // 去掉空的列
    headColHasVal := helper.GridDoesHeaderColumnHasValue(headerGrid[HeaderLineNameRowIndex]) // 判断列是否有效(至少要包含一个值)
    //dataColHasVal := helper.GridDoesColumnHasValue(dataGrid)                                 // 判断列是否有效(至少要包含一个值)
    dataGrid = helper.GridFilterColumn(dataGrid, headColHasVal)     // 删去空列
    headerGrid = helper.GridFilterColumn(headerGrid, headColHasVal) // 删去空列

    // 去掉没有s标记的列
    flagColValid := colHasSvrFlag(headerGrid[HeaderLineFlagRowIndex]) // 判断列是否有效(至少要包含一个值)
    dataGrid = helper.GridFilterColumn(dataGrid, flagColValid)        // 删去非 s 列
    headerGrid = helper.GridFilterColumn(headerGrid, flagColValid)    // 删去非 s 列

    dataGrid = helper.GridOmitEmptyRow(dataGrid) // 删除空行
    if len(dataGrid) == 0 {
        return nil
    }

    // type 映射到通用的 type , 找不到匹配的项会报错
    headerGrid[HeaderLineTypeRowIndex] = fmtTypeName(headerGrid[HeaderLineTypeRowIndex])

    // name 补全, 格式化
    headerGrid[HeaderLineNameRowIndex] = fillName(headerGrid[HeaderLineNameRowIndex])

    return &ConfigCsv{
        Csv:            item,
        HeaderLineName: headerGrid[HeaderLineNameRowIndex],
        HeaderLineDesc: headerGrid[HeaderLineDescRowIndex],
        HeaderLineType: headerGrid[HeaderLineTypeRowIndex],
        HeaderLineFlag: headerGrid[HeaderLineFlagRowIndex],
        Grid:           dataGrid,
    }
}

func fillName(nameRow []string) []string {
    return lo.Map(nameRow, func(item string, index int) string {
        if item == "" {
            return "col" + strconv.Itoa(index+1)
        }
        return item
    })
}

// fmtTypeName 格式化 type 名称
func fmtTypeName(typeRow []string) []string {
    return lo.Map(typeRow, func(item string, index int) string {
        return helper.FormatTypeName(item)
    })
}

// isValidType 判断列是否有有效的 type
func isValidType(typeRow []string) []bool {
    allValidTypes := helper.AllValidTypeNames()
    return lo.Map(typeRow, func(item string, index int) bool {
        return lo.Contains(allValidTypes, item)
    })
}

// colHasSvrFlag 判断列是否有 s 字符
func colHasSvrFlag(flagRow []string) []bool {
    return lo.Map(flagRow, func(item string, index int) bool {
        return strings.Index(item, "s") != -1
    })
}
