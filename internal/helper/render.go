package helper

import (
	"bufio"
	"log"
	"os"
	"text/template"
)

// RenderTemplate 渲染模板
func RenderTemplate(outputFilePath string, templateFilePath string, param any) {
	// 渲染模板
	f, err := os.OpenFile(outputFilePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		log.Fatal("open file error: ", err)
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	t, err := template.ParseFiles(templateFilePath)
	if err != nil {
		log.Fatal("parse template error: ", err)
	}
	err = t.Execute(writer, param)
	if err != nil {
		log.Fatal("execute template error: ", err)
	}
}
