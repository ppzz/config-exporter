package csver

type Csv struct {
	FilePath string
	Grid     [][]string
}

func NewCsv(filePath string, grid [][]string) *Csv {
	return &Csv{FilePath: filePath, Grid: grid}
}
