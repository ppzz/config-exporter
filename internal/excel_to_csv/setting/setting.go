package setting

import (
	"github.com/go-playground/validator/v10"
	"github.com/ppzz/golang-csv/internal/helper"
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
	// 日志相关
	LogFilename string   `validate:"required"`
	LogEncoder  string   `validate:"required,oneof=console json"`
	LogOutput   []string `validate:"required,min=1"`
	LogLevel    string   `validate:"required,oneof=debug info warn error dpanic panic fatal"`

	InputExcelDir string `validate:"required,dirpath"` // 输入的 Excel 文件夹
	OutputCsvDir  string `validate:"required,dirpath"` // 输出的 Csv 文件夹
	// OutputCsvMetaDir string `validate:"required,dir"` // 输出的 Csv Meta 文件夹
}

func NewSetting() *Setting {
	return &Setting{}
}

func (s *Setting) Name() string {
	return helper.TypeName(s)
}

func (s *Setting) Init() {
	v := viper.GetViper()
	s.SetAttribute(v)

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.Struct(s)
	cobra.CheckErr(err)

	// if err != nil {
	// 	t, _ := json.Marshal(s)
	// 	fmt.Println(string(t))
	// 	cobra.CheckErr(err)
	// 	helper.AssertNoErrWithFmt("inits failed: "+s.Name(), err)
	// }
}

func (s *Setting) SetAttribute(v *viper.Viper) {
	v.SetDefault("log.filename", "app")
	v.SetDefault("log.level", "debug")
	v.SetDefault("log.encoder", "console")
	v.SetDefault("log.output", []string{"stdout"})

	// 日志相关
	s.LogFilename = v.GetString("log.filename")
	s.LogLevel = v.GetString("log.level")
	s.LogEncoder = v.GetString("log.encoder")
	s.LogOutput = v.GetStringSlice("log.output")

	s.InputExcelDir = v.GetString("excel")
	s.OutputCsvDir = v.GetString("output")
}

// ---------------------------------------------------------------------------------------------------------------------

func getNameWithEnv(filename string, e string) string {
	if e == "" {
		return filename
	}
	return filename + "." + e
}
