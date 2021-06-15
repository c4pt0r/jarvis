package main

import (
	"io"
	"os"

	"github.com/google/uuid"
)

func GenerateUUID() string {
	return uuid.NewString()
}

type StdinReaderCloser struct{}

func NewStdinReaderCloser() io.ReadCloser {
	return &StdinReaderCloser{}
}

func (s *StdinReaderCloser) Read(p []byte) (n int, err error) {
	return os.Stdin.Read(p)
}

func (s *StdinReaderCloser) Close() error {
	return nil
}

type StdoutWriterCloser struct{}

func NewStdoutWriterCloser() io.WriteCloser {
	return &StdoutWriterCloser{}
}

func (s *StdoutWriterCloser) Write(p []byte) (n int, err error) {
	return os.Stdout.Write(p)
}

func (s *StdoutWriterCloser) Close() error {
	return nil
}
