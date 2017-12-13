package models

import (
	"github.com/satori/go.uuid"
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
	QueueId       int64  `json:"queue_id,omitempty"`
	StorageDriver string `json:"storage_drive,-"`
	Timestamp     int    `json:"timestamp"`
}

func (m *Message) setUUID() {
	m.UUID = uuid.NewV4().String()
}

func (m *Message) ToSqlInsert() string {
	if m.UUID == "" {
		m.setUUID()
	}
	var buffer bytes.Buffer
	buffer.WriteString("(")
	buffer.WriteString(parametrizeString(m.UUID))
	buffer.WriteString(",")
	buffer.WriteString(parametrizeInt(m.QueueId))
	buffer.WriteString(",")
	buffer.WriteString(parametrizeInt(m.Size))
	buffer.WriteString(",")
	buffer.WriteString(parametrizeString(m.Data))
	buffer.WriteString(",")
	buffer.WriteString(parametrizeString(m.StorageDriver))

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

func (m *Message) SetStorageDrive(s string) {
	m.StorageDriver = s
}

func (m *Message) SetData(d string) {
	m.Data = d
}

func (m *Message) SetTimestamp(d string) {

}
