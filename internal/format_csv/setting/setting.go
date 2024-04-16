package setting

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var ConstExcelExtNames = []string{".xlsx", ".xls"} // excel 文件后缀

var instance = NewSetting()

func Get() *Setting {
	return instance
}

// ---------------------------------------------------------------------------------------------------------------------

type Setting struct {
	InputCsvDir  string `validate:"required"` // 输出的 Csv 文件夹
	OutputFmtDir string `validate:"required"` // 输入的 Excel 文件夹

	CsvFileNameSchema string `validate:"required"` // csv 文件名的正则表达式
}

func NewSetting() *Setting {
	return &Setting{}
}

func (s *Setting) Init() {
	v := viper.GetViper()
	s.SetAttribute(v)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(s)
	cobra.CheckErr(err)
}

func (s *Setting) SetAttribute(v *viper.Viper) {
	v.SetDefault("csv_filename_schema", `^\w+\.\d+\.\w+\.csv$`)

	s.InputCsvDir = v.GetString("csv")
	s.OutputFmtDir = v.GetString("fmt")
	s.CsvFileNameSchema = v.GetString("csv_filename_schema")
}
