package services

import (
	"../models"
	"bytes"
	"sync"
)

type MessagingInterface interface {
	Push(m *models.Message)
	Pop(queueId string, n int) bytes.Buffer
}

func InitMessaging() *Messaging {
	ser := Messaging{fileQueues: make(map[string]*FileQueue), mutex: sync.Mutex{}}
	return &ser
}

type Messaging struct {
	mutex sync.Mutex
	fileQueues map[string]*FileQueue
}


func (s *Messaging) Push(m *models.Message) {
	//Create a single file descriptor for each queue_id
	if _, ok := s.fileQueues[m.QueueId]; !ok {
		s.mutex.Lock()
		if _, ok := s.fileQueues[m.QueueId]; !ok {
			s.fileQueues[m.QueueId] = NewFileQueue(m.QueueId, 1024*1024*500)
		}
		s.mutex.Unlock()
	}
	b := bytes.Buffer{}
	b.WriteString(m.ToSqlInsert())
	s.fileQueues[m.QueueId].Push(b)
}

func (s *Messaging) Pop(queueId string, n int) bytes.Buffer {
	b := bytes.Buffer{}
	if _, ok := s.fileQueues[queueId]; !ok {
		return b
	}
	return s.fileQueues[queueId].Pop(n)
}