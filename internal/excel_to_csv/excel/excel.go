package excel

import (
	"github.com/samber/lo"
	"github.com/tealeg/xlsx"
	"strings"
)

type Excel struct {
	FilePath string   // 文件名，表明读取的excel，带有路径和 ext后缀
	Size     int      // 文件大小 byte
	Sheets   []*Sheet // 表单
}

// Sheet 对 sheet 的映射，用来存放数据
type Sheet struct {
	Index          int        // sheet序号
	Name           string     // sheet名字
	OriginalMaxRow int        // 原始sheet穿过来的max行数
	OriginalMaxCol int        // 原始sheet穿过来的max列数
	RowCount       int        // 数据的有效行数(不算metadata，排除掉文件的空行)
	ColCount       int        // 列数
	Grid           [][]string // 有效数据表格, grid[0] 表示第一个数据行，grid[0][1]： 表示第一个数据行的第二个格子的值
}

func NewSheet(sheetIndex int, sheet *xlsx.Sheet) *Sheet {
	tempGrid := lo.Map(sheet.Rows, func(row *xlsx.Row, _ int) []string {
		return lo.Map(row.Cells, func(cell *xlsx.Cell, _ int) string {
			return strings.TrimSpace(cell.String())
		})
	})

	_, dataGrid := SplitGrid(tempGrid, MetaDataRowCount) // tempGrid 拆成两部分
	dataGrid = RemoveEmptyRow(tempGrid)                  // 删除空行
	columnCount := GetColumnCount(tempGrid)              // 计算有效列
	tempGrid = Tailor(columnCount, tempGrid)             // 裁剪多余的列，并补足不够的列

	return &Sheet{
		Index:          sheetIndex,
		Name:           sheet.Name,
		OriginalMaxRow: sheet.MaxRow,
		OriginalMaxCol: sheet.MaxCol,
		RowCount:       len(tempGrid),
		ColCount:       columnCount,
		Grid:           dataGrid,
	}
}

// ColumnMeta column 的描述信息
type ColumnMeta struct {
	Index int    // 列序号
	Name  string // 列名字 (对应sheet的第一行)
	Desc  string // 列描述（对应sheet的第二行）
	Typ   string // 列数据类型（对应sheet的第三行）
	Flag  string // 列 c/s 标记（对应sheet的第四行）
}

func NewColumnMeta(index int, name, desc, typ, flag string) ColumnMeta {
	t := ColumnMeta{
		Index: index,
		Name:  name,
		Desc:  desc,
		Typ:   typ,
		Flag:  flag,
	}
	return t
}
