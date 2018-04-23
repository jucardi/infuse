package ioutils

import (
	"io"
	"io/ioutil"
	"os"
)

type fileWriter struct {
	filename string
}

func (w *fileWriter) Write(p []byte) (n int, err error) {
	if err := ioutil.WriteFile(w.filename, p, os.FileMode(644)); err != nil {
		return 0, err
	}
	return len(p), nil
}

// NewFileWriter creates a new instance of a file writer
func NewFileWriter(filename string) io.Writer {
	return &fileWriter{filename: filename}
}
