package log

import (
	"log"
)

// AsyncWriter is asynchronous writer, writes in separate goroutine
type AsyncWriter struct {
	lines chan []byte
}

// NewAsyncWriter creates writer and starts writing routine
func NewAsyncWriter() *AsyncWriter {
	aw := &AsyncWriter{
		lines: make(chan []byte, 100),
	}
	go aw.writerLoop()
	return aw
}

func (aw *AsyncWriter) Write(v []byte) (int, error) {
	aw.lines <- v
	return 0, nil
}

func (aw *AsyncWriter) writerLoop() {
	for {
		line := <-aw.lines
		log.Println(string(line))
	}
}
