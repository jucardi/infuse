package ioutils

import (
	"io"
	"io/ioutil"
	"os"
)

type fileWritter struct {
	filename string
}

func (w *fileWritter) Write(p []byte) (n int, err error) {
	if err := ioutil.WriteFile(w.filename, p, os.FileMode(644)); err != nil {
		return 0, err
	}
	return len(p), nil
}

// NewFileWritter creates a new instance of a file writter
func NewFileWritter(filename string) io.Writer {
	return &fileWritter{filename: filename}
}
