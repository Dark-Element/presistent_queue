package services

import (
	"presistentQueue/models"
	"database/sql"
	"bytes"
	"sync"
	"log"
	"strconv"
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
		aggValues: make(map[int64]int64),
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
	aggValues    map[int64]int64
	queries      chan [2]bytes.Buffer
}

func (s *Messaging) Push(m *models.Message) {
	//change storage driver
	//todo: write this part

	s.mutex.Lock()
	//insert to buffer
	s.insertValues.WriteString(m.ToSqlInsert())
	s.aggValues[m.QueueId] += m.Size
	//check if need to flush
	s.currentSize++
	if s.currentSize >= s.buffer {
		s.queries <- [2]bytes.Buffer{s.insertValues,s.aggToSQLInsert()}
		s.insertValues = bytes.Buffer{}
		s.aggValues =  make(map[int64]int64)
		s.currentSize = 0
	} else {
		s.insertValues.WriteString(",\n")
	}
	s.mutex.Unlock()

}

func (s *Messaging) flush() {
	s.queries = make(chan [2]bytes.Buffer, 1)
	for ins := range s.queries {
		//build insert query
		iq := buildInsertQuery(ins[0])
		aq := buildAggQuery(ins[1])
		go func() {
			_, err := s.dbConn.Exec(iq)
			if err != nil {
				log.Println(err)
				return
			}

			_, errA := s.dbConn.Exec(aq)
			if errA != nil {
				log.Println(errA)
				return
			}

		}()

	}

}

func (m *Messaging) aggToSQLInsert() bytes.Buffer{
	b := bytes.Buffer{}
	for idx, val := range m.aggValues{
		b.WriteString("(")
		b.WriteString(strconv.FormatInt(idx,10))
		b.WriteString(",")
		b.WriteString(strconv.FormatInt(val,10))
		b.WriteString(")\n")
	}
	return b
}

func buildInsertQuery(insertValues bytes.Buffer) string{
	iq := bytes.Buffer{}
	iq.WriteString("INSERT INTO `messages` (`uuid`, `queue_id`, `size`,  `storage_driver`, `data`, `timestamp`) VALUES ")
	iq.WriteString(insertValues.String())
	return iq.String()
}


func buildAggQuery(insertValues bytes.Buffer) string{
	iq := bytes.Buffer{}
	iq.WriteString("INSERT INTO `aggregation` (`queue_id`, `current_size`) VALUES ")
	iq.WriteString(insertValues.String())
	iq.WriteString("ON DUPLICATE KEY UPDATE `current_size` = `current_size` + values (`current_size`)")
	return iq.String()
}



func (s *Messaging) Pop(n int) {}

func (s *Messaging) Ack(m chan *models.Message) {}
