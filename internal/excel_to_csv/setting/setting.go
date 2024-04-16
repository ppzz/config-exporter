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
	InputExcelDir string `validate:"required"` // 输入的 Excel 文件夹
	OutputCsvDir  string `validate:"required"` // 输出的 Csv 文件夹
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
	s.InputExcelDir = v.GetString("excel")
	s.OutputCsvDir = v.GetString("csv")
}
