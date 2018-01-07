package adapters

import (
	"io"
)

type MemoryQueue struct{}

func (mq *MemoryQueue) Push(data []byte){}
func (mq *MemoryQueue) Pop(n int64, s int64) io.Reader{
	a, _ := io.Pipe()
	return a

}
func (mq *MemoryQueue) Peek() (int64, int64) {
	return 999999,999999
}
func (mq *MemoryQueue) CanPush(s int, atomic bool) bool {return true}
func (mq *MemoryQueue) Close(){}