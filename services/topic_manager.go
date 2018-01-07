package services

import (
	"persistentQueue/models"
	"io"
	"persistentQueue/adapters"
)

type TopicManagerInterface interface {
	Push(data []byte, atomic bool)
	Pop(limit int64) io.Reader
	Peek() models.TopicContent
	Close()
}



func NewTopicManager(queue_id string, max_file_size int64) TopicManagerInterface{
	arr := []adapters.QueueInterface{&adapters.MemoryQueue{}, &adapters.FileQueue{}}
	return  &TopicManager{queues: arr}
}
//I think applying all the underlying adapters as circuit breakers will work better than statically defining them
type TopicManager struct{
	queues []adapters.QueueInterface

	count int64
	size int64
}

func (t *TopicManager) Push(data []byte, atomic bool){

}
func (t *TopicManager) Pop(limit int64) io.Reader{
	r, _ := io.Pipe()
	return r
}
func (t *TopicManager) Peek() models.TopicContent{
	return models.TopicContent{}
}
func (t *TopicManager) Close(){}
