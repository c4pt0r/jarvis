package main

import (
	"encoding/json"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	jaegerlog "github.com/uber/jaeger-client-go/log"
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

func CreateTracer(servieName string) (opentracing.Tracer, io.Closer, error) {
	var cfg = jaegercfg.Configuration{
		ServiceName: servieName,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  jaeger.SamplerTypeConst,
			Param: 1,
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans: true,
		},
	}

	jLogger := jaegerlog.StdLogger
	tracer, closer, err := cfg.NewTracer(
		jaegercfg.Logger(jLogger),
	)
	return tracer, closer, err
}

func EncodeMap(m map[string]string) string {
	b, _ := json.Marshal(m)
	return string(b)
}
