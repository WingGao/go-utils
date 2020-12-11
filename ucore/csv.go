package ucore

import "github.com/ungerik/go-dry"

type CSVWriter struct {
	FilePath string
	StringBuilder
}

func NewCSVWriter(fp string) (*CSVWriter, error) {
	w := &CSVWriter{
		FilePath: fp,
	}
	w.Write("\uFEFF") //BOM,excel直接打开不乱码
	return w, nil
}

func (m *CSVWriter) Save() error {
	return dry.FileSetBytes(m.FilePath, m.Bytes())
}
