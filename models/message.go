package models

import (
	"bytes"
	"strconv"
)

type PushedMessage struct {
	Size float64 `json:"size"`
	Data string  `json:"data"`
}

type Message struct {
	Size          int64  `json:"size"`
	Data          string `json:"data"`
	UUID          string `json:"uuid,omitempty"`
	QueueId       string `json:"queue_id,omitempty"`
	StorageDriver string `json:"storage_drive,-"`
	Timestamp     int64  `json:"timestamp"`
}


func (m *Message) ToSqlInsert() string {
	var buffer bytes.Buffer
	buffer.WriteString("(")
	buffer.WriteString(parametrizeString(m.Data))
	buffer.WriteString(",")
	buffer.WriteString(parametrizeInt(m.Timestamp))

	buffer.WriteString(")")
	return buffer.String()
}

func parametrizeString(value string) string {
	var buffer bytes.Buffer
	buffer.WriteString("'")
	buffer.WriteString(value)
	buffer.WriteString("'")
	return buffer.String()
}

func parametrizeInt(value int64) string {
	var buffer bytes.Buffer
	buffer.WriteString(strconv.FormatInt(value, 10))
	return buffer.String()
}