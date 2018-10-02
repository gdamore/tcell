package log

import "log"

// DefaultWriter is simple wrapper for log from stdlib
type DefaultWriter struct{}

func (dw DefaultWriter) Write(p []byte) (int, error) {
	log.Println(string(p))
	return len(p), nil
}
