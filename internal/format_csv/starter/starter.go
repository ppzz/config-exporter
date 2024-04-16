package starter

import (
	"encoding/csv"
	csver "github.com/ppzz/golang-csv/internal/format_csv/csv"
	"github.com/ppzz/golang-csv/internal/format_csv/setting"
	"github.com/ppzz/golang-csv/internal/helper"
	"github.com/samber/lo"
	"io"
	"log"
	"os"
	"path"
	"regexp"
)

func Start() {
	csvDir := setting.Get().InputCsvDir
	fmtDir := setting.Get().OutputFmtDir
	csvFilenameSchema := setting.Get().CsvFileNameSchema

	csvFilenameList := helper.ListFilenameByExt(csvDir, ".csv")
	csvFilenameList = lo.Filter(csvFilenameList, func(item string, index int) bool {
		return regexp.MustCompile(csvFilenameSchema).MatchString(item)
	})
	csvFilenameList = lo.Map(csvFilenameList, func(item string, index int) string {
		return path.Join(csvDir, item)
	})

	csvList := lo.Map(csvFilenameList, func(item string, index int) *csver.Csv {
		grids := readCsv(item)
		return csver.NewCsv(item, grids)
	})

	configCsvList := lo.Map(csvList, func(item *csver.Csv, index int) *csver.ConfigCsv {
		return csver.CreateConfigCsv(item)
	})

	helper.MakeSureExist(fmtDir)

	emptyDir := helper.IsEmptyDir(fmtDir)
	if !emptyDir {
		log.Fatal("output dir is not empty: ", fmtDir)
	}
	lo.ForEach(configCsvList, func(item *csver.ConfigCsv, index int) {
		bareName := helper.FileBareName(item.Csv.FilePath)
		n := helper.CamelCaseToSnakeCase(helper.FilenameByType(bareName))

		newCsvFilePath := path.Join(fmtDir, n+".csv")
		saveCsv(item, newCsvFilePath)
	})

	//  0. csv 文件过滤
	// - 按照cs标记过滤列
	// - 按照 type 过滤列, 如果有不符合的列，报错
	// - type 映射到通用的 type
	// - name 补全,格式化

	// header line: name, type, desc, c-s-flag
}

func saveCsv(item *csver.ConfigCsv, csvPath string) {
	file, err := os.Create(csvPath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	headerGrid := [][]string{
		item.HeaderLineName,
		item.HeaderLineDesc,
		item.HeaderLineType,
		item.HeaderLineFlag,
	}

	// 写入数据
	err = writer.WriteAll(headerGrid) // 自动调用 Flush
	if err != nil {
		log.Fatalf("Failed to write data to CSV: %v", err)
	}

	// 写入数据
	err = writer.WriteAll(item.Grid) // 自动调用 Flush
	if err != nil {
		log.Fatalf("Failed to write data to CSV: %v", err)
	}

	// 检查是否有任何写入错误
	if err := writer.Error(); err != nil {
		log.Fatalf("Error writing CSV: %v", err)
	}
}

func readCsv(item string) [][]string {
	// 打开文件
	file, err := os.Open(item)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close() // 确保在函数结束时关闭文件

	// 创建 CSV 读取器
	reader := csv.NewReader(file)

	// 逐行读取
	var lines [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break // 文件结束，退出循环
		}
		if err != nil {
			log.Fatalf("Failed to read CSV: %v", err)
		}
		lines = append(lines, record)
	}
	return lines
}
