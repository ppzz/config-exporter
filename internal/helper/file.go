package helper

import (
	"encoding/csv"
	"io"
	"log"
	"os"
)

func FileCsvRead(item string) [][]string {
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

func FileCsvWrite(csvFilePath string, grid [][]string) {
	file, err := os.Create(csvFilePath)
	if err != nil {
		log.Fatalf("Failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// 写入数据
	err = writer.WriteAll(grid) // 自动调用 Flush
	if err != nil {
		log.Fatalf("Failed to write data to CSV: %v", err)
	}

	// 检查是否有任何写入错误
	if err := writer.Error(); err != nil {
		log.Fatalf("Error writing CSV: %v", err)
	}
}
