package adapters

import (
	"io"
)

type QueueInterface interface {
	Push(data []byte)
	Pop(n int64, s int64) io.Reader
	Peek() (int64, int64)
	CanPush(s int, atomic bool) bool
	Close()
}