package services

import (
	"persistentQueue/models"
	"io"
	"persistentQueue/adapters"
	"sync"
)

type TopicManagerInterface interface {
	Push(data []byte, atomic bool)
	Pop(limit int64, LimitSize int64) io.Reader
	Peek() models.TopicContent
	Prefix() string
	Close()
}

func NewTopicManager(queue_id string, max_file_size int64) TopicManagerInterface {
	arr := []adapters.QueueInterface{adapters.InitMemoryQueue(1024), adapters.NewFileQueue(queue_id, max_file_size)}
	return &TopicManager{queues: arr, prefix: queue_id}
}

//I think applying all the underlying adapters as circuit breakers will work better than statically defining them
type TopicManager struct {
	queues []adapters.QueueInterface
	prefix string

	count int64
	size  int64
	sync.Mutex
}

func (t *TopicManager) Push(data []byte, atomic bool) {
	for _, a := range t.queues {
		if a.CanPush(int64(len(data)), atomic) {
			a.Push(data)
			break
		}
	}
}
func (t *TopicManager) Pop(targetCount int64, targetSize int64) io.Reader {
	r, w := io.Pipe()
	go func() {
		t.Lock()
		defer t.Unlock()
		defer w.Close()
		if !t.CanPop(targetCount, targetSize){
			return
		}
		for _, a := range t.queues {
			if targetCount == 0 || targetSize == 0{
				break
			}
			out := make(chan []byte, targetCount)
			go a.Pop(out, targetCount, targetSize)
			for {
				msg, ok := <-out
				if !ok {
					break
				}
				targetCount -= 1
				targetSize -= int64(len(msg))
				w.Write(msg)
			}
		}

	}()

	return r
}

func (t *TopicManager) Peek() models.TopicContent {
	totalCount := int64(0)
	totalSize := int64(0)
	for _, a := range t.queues{
		count, size := a.Peek()
		totalCount += count
		totalSize += size
	}
	return models.TopicContent{
		Count: totalCount,
		Size: totalSize,
	}
}

func (t *TopicManager) CanPop(targetCount int64, targetSize int64) bool {
	current := t.Peek()
	return (targetCount > 0 && current.Count >= targetCount) || (targetSize > 0 && current.Size >= targetSize )
}

func (t *TopicManager) Close() {}

func (t *TopicManager) Prefix() string {
	return t.prefix
}
