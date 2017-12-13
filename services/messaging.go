package services

import (
	"presistentQueue/models"
	"database/sql"
	"bytes"
	"sync"
	"log"
)

type MessagingInterface interface {
	Push(m *models.Message)
	Pop(n int)
	Ack(m chan *models.Message)
}

func InitMessaging(db *sql.DB, buffer int64) *Messaging {
	ser := Messaging{
		dbConn: db,
		buffer: buffer,
	}
	go ser.flush()
	return &ser
}

type Messaging struct {
	mutex        sync.Mutex
	dbConn       *sql.DB
	buffer       int64
	currentSize  int64 //why not use a channel for sync: in the begging of the flow, msg's size will be the whole data so i wont be able to buffer a lot of data in memory conservative system until the data will be uploaded to external storage
	insertValues bytes.Buffer
	aggValues    bytes.Buffer
	queries      chan bytes.Buffer
}

func (s *Messaging) Push(m *models.Message) {
	//change storage driver
	//todo: write this part

	s.mutex.Lock()
	//insert to buffer
	s.insertValues.WriteString(m.ToSqlInsert())
	//check if need to flush
	if s.currentSize >= s.buffer {

		s.queries <- s.insertValues
		s.insertValues = bytes.Buffer{}
		s.currentSize = 0
	} else {
		s.currentSize++
		s.insertValues.WriteString(",\n")
	}
	s.mutex.Unlock()

}

func (s *Messaging) flush() {
	s.queries = make(chan bytes.Buffer, 1)
	for ins := range s.queries {
		//build insert query
		iq := bytes.Buffer{}
		iq.WriteString("INSERT INTO `messages` (`uuid`, `queue_id`, `size`, `data`, `storage_driver`) VALUES ")
		iq.WriteString(ins.String())
		log.Println(iq.String())
		go func() {
			_, err := s.dbConn.Query(iq.String())
			if err != nil {
				log.Println(err)
			}
		}()

	}

}

func (s *Messaging) Pop(n int) {}

func (s *Messaging) Ack(m chan *models.Message) {}
