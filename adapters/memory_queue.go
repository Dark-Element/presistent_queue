package adapters

import (
	"sync"
)

type MemoryQueue struct {
	data      chan []byte

	maxSize int64

	sizeBytes int64
	sync.Mutex
}

func InitMemoryQueue(buffer int64) *MemoryQueue {
	c := make(chan []byte, buffer)
	return &MemoryQueue{data: c, maxSize: buffer}
}

func (mq *MemoryQueue) Push(data []byte) {
	mq.data <- append(data, "\n"...)
	mq.sizeIncr(int64(len(data)))
}
func (mq *MemoryQueue) Pop(out chan []byte, targetCount int64, targetSize int64) {
	for len(mq.data) > 0 && (targetCount > 0 || targetSize > 0) {
		d := <-mq.data
		out <- d
		mq.sizeDecr(1)
		targetCount--
		targetSize -= int64(len(d))
	}
	close(out)
}
func (mq *MemoryQueue) Peek() (int64, int64) {
	mq.Lock()
	defer mq.Unlock()
	return int64(len(mq.data)), mq.sizeBytes
}
func (mq *MemoryQueue) CanPush(s int64, atomic bool) bool {
	mq.Lock()
	defer mq.Unlock()
	return !atomic && mq.maxSize > s + mq.sizeBytes
}
func (mq *MemoryQueue) Close() {}

func (mq *MemoryQueue) sizeIncr(incr int64) {
	mq.Lock()
	defer mq.Unlock()
	mq.sizeBytes += incr
}

func (mq *MemoryQueue) sizeDecr(decr int64) {
	mq.Lock()
	defer mq.Unlock()
	mq.sizeBytes -= decr
}
