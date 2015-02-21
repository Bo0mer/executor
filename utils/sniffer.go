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

// Returns new Sniffer, which is sniffing on the specified writer
// All sniffed data is kept in memory.
func NewSniffer(writer io.Writer) Sniffer {
	return &sniffer{
		data:   new(bytes.Buffer),
		writer: writer,
	}
}

// Write to the sniffer's writer. The data will be sniffed
func (s *sniffer) Write(data []byte) (n int, err error) {
	n, err = s.writer.Write(data)
	s.data.Write(data[:n])
	return
}

// Returns the sniffed data, got by Write
func (s *sniffer) SniffedData() []byte {
	return s.data.Bytes()
}
