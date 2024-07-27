package sloger

import "io"

type MultiWriter struct {
	writers []io.Writer
}

func (m *MultiWriter) Write(p []byte) (n int, err error) {
	var maxN int
	for _, writer := range m.writers {
		n, err := writer.Write(p)
		if err != nil {
			return maxN, err
		}
		maxN = max(n, maxN)
	}
	return maxN, nil
}

func NewMultiWriter(mw ...io.Writer) io.Writer {
	return &MultiWriter{
		writers: mw,
	}
}
