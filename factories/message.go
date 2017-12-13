package factories

import (
	"presistentQueue/models"
	"encoding/json"
	"fmt"
	"time"
	"github.com/satori/go.uuid"
)

func Messages(rb string) *models.Message {
	if m := messageFromJSONString(rb); m != nil {
		return m
	}
	return nil
}

func messageFromJSONString(s string) *models.Message {
	var m models.PushedMessage
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println("Error parsing " + s)
		return nil
	}
	return &models.Message{Size:int64(m.Size),
		Data:m.Data,
		UUID: uuid.NewV4().String(),
		QueueId: 123123,
		StorageDriver: "inline",
		Timestamp: int(time.Now().UnixNano() / 1000000)}
}
