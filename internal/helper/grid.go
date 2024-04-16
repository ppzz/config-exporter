package helper

import (
	"github.com/samber/lo"
	"log"
	"strings"
)

// GridSplitHeader 分裂成两个grid
func GridSplitHeader(grid [][]string, headerRowCount int) (headerGrid, dataGrid [][]string) {
	if len(grid) <= headerRowCount {
		return grid, [][]string{}
	}

	return grid[:headerRowCount], grid[headerRowCount:]
}

// GridTrimCellSpace 去掉每个cell的空格
func GridTrimCellSpace(grid [][]string) [][]string {
	for i, row := range grid {
		for j, cell := range row {
			grid[i][j] = strings.TrimSpace(cell)
		}
	}
	return grid
}

// GridMaxColumnLen 返回有效grid的有效列数（不包含每行尾部的空值）
func GridMaxColumnLen(grids [][]string) int {
	return lo.Max(lo.Map(grids, func(row []string, index int) int { return len(row) }))
}

// GridFillColumn 填充列
func GridFillColumn(grid [][]string, colCount int) [][]string {
	rows := make([][]string, len(grid))
	for i, row := range grid {
		if len(row) > colCount {
			log.Fatal("Error: row len >= col")
		}

		for len(row) < colCount {
			row = append(row, "")
		}
		rows[i] = row
	}
	return rows
}

// GridOmitEmptyRow 删除空行
func GridOmitEmptyRow(grid [][]string) [][]string {
	return lo.Filter(grid, func(row []string, index int) bool {
		return !RowIsEmpty(row)
	})
}

// RowIsEmpty 判断行是否为空
func RowIsEmpty(row []string) bool {
	return lo.EveryBy(row, func(item string) bool {
		return len(strings.TrimSpace(item)) == 0
	})
}

// GridDoesColumnHasValue 判断grid的列是否有值
func GridDoesColumnHasValue(grid [][]string) []bool {
	validCol := make([]bool, len(grid[0]))
	for _, row := range grid {
		for i, cell := range row {
			if len(strings.TrimSpace(cell)) > 0 {
				validCol[i] = true
			}
		}
	}
	return validCol
}

// GridFilterColumn 过滤grid的列
func GridFilterColumn(grid [][]string, colHasVal []bool) [][]string {
	return lo.Map(grid, func(row []string, index int) []string {
		var newRow []string
		for i, cell := range row {
			if colHasVal[i] {
				newRow = append(newRow, cell)
			}
		}
		return newRow
	})
}
