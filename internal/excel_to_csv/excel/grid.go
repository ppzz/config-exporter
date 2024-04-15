package excel

import (
	"github.com/samber/lo"
	"strings"
)

// 与 [][]string 相关的操作

// SplitGrid 分裂成两个grid
func SplitGrid(grid [][]string, rowCount int) (metaGrid, dataGrid [][]string) {
	count := len(grid)
	if count <= rowCount {
		metaGrid = grid
		return
	}
	metaGrid = grid[:rowCount]
	dataGrid = grid[rowCount:count]

	return
}

// Tailor 裁剪or补足grid的列
func Tailor(colCount int, grids [][]string) [][]string {
	for i, row := range grids {
		if len(row) >= colCount {
			grids[i] = row[:colCount]
			continue
		}

		for len(row) < colCount {
			row = append(row, "")
		}
		grids[i] = row
	}
	return grids
}

// GetColumnCount 返回有效grid的有效列数（不包含每行尾部的空值）
func GetColumnCount(grids [][]string) int {
	count := 0
	for _, row := range grids {
		countOfRow := len(row)
		if countOfRow <= count {
			continue
		}
		for i := count; i < countOfRow; i++ {
			if len(strings.TrimSpace(row[i])) > 0 {
				count = i + 1
			}
		}
	}
	return count
}

// RemoveEmptyRow 删除空行
func RemoveEmptyRow(grid [][]string) [][]string {
	return lo.Filter(grid, func(row []string, index int) bool {
		return !IsAllEmpty(row)
	})
}

func IsAllEmpty(row []string) bool {
	return lo.EveryBy(row, func(item string) bool {
		return len(strings.TrimSpace(item)) == 0
	})
}
