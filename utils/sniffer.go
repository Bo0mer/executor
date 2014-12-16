package utils

import (
	"bytes"
	"io"
)

type Sniffer interface {
	Write([]byte) (int, error)
	SniffedData() []byte
}

type sniffer struct {
	data   *bytes.Buffer
	writer io.Writer
}

func NewSniffer(writer io.Writer) Sniffer {
	return &sniffer{
		data:   new(bytes.Buffer),
		writer: writer,
	}
}

func (s *sniffer) Write(data []byte) (n int, err error) {
	n, err = s.writer.Write(data)
	s.data.Write(data[0:n])
	return
}

func (s *sniffer) SniffedData() []byte {
	return s.data.Bytes()
}
