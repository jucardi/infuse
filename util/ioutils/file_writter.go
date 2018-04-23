package ioutils

import (
	"io"
	"os"
	"github.com/jucardi/go-infuse/util/log"
)

type fileWriter struct {
	filename string
	file     *os.File
}

func (w *fileWriter) Write(p []byte) (n int, err error) {
	log.Debug("<--- ioutils.Write")
	return w.file.Write(p)
}

func (w *fileWriter) Close() error {
	return w.file.Close()
}

// NewFileWriter creates a new instance of a file writer
func NewFileWriter(filename string) (io.WriteCloser, error) {
	f, err := os.Create(filename)
	if err != nil {
		return nil, err
	}
	return &fileWriter{filename: filename, file: f}, nil
}
