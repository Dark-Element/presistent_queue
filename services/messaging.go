package services

import (
	"persistentQueue/models"
	"sync"
	"fmt"
	"io"
)

type MessagingInterface interface {
	Push(m *models.Message, flush bool)
	Pop(queueId string, n int) io.Reader
	Close()
}

func InitMessaging() *Messaging {
	ser := Messaging{topicManager: make(map[string]TopicManagerInterface), mutex: sync.RWMutex{}}
	return &ser
}

type Messaging struct {
	mutex        sync.RWMutex
	topicManager map[string]TopicManagerInterface
}


func (s *Messaging) Push(m *models.Message) {
	//Create a single file descriptor for each queue_id
	s.mutex.RLock()
	val, ok := s.topicManager[m.QueueId]
	s.mutex.RUnlock()

	if !ok{
		s.mutex.Lock()
		if val, ok = s.topicManager[m.QueueId]; !ok {
			val = NewTopicManager(m.QueueId, 1024*1024*500)
			s.topicManager[m.QueueId] = val
		}
		s.mutex.Unlock()
	}

	val.Push(m.Data, false) //todo: add to message struct
}

func (s *Messaging) Pop(queueId string, n int64) io.Reader {
	if _, ok := s.topicManager[queueId]; !ok {
		return nil
	}
	return s.topicManager[queueId].Pop(n)


}

func (s *Messaging) Close(){
	s.mutex.Lock()
	fmt.Println("Closing messaging service")
	defer fmt.Println("Closed messaging service")
	if len(s.topicManager) == 0{
		return
	}
	for _, q := range s.topicManager {
		fmt.Println("Closing " + q.prefix)
		q.Close()
		fmt.Println("Closed " + q.prefix)
	}

}