package factories

import (
	"../models"
	"encoding/json"
	"fmt"
	"time"
)

func Messages(rb string, qi string) *models.Message {
	if m := messageFromJSONString(rb, qi); m != nil {
		return m
	}
	return nil
}

func messageFromJSONString(s string, qi string) *models.Message {
	var m models.PushedMessage
	err := json.Unmarshal([]byte(s), &m)
	if err != nil {
		fmt.Println("Error parsing " + s)
		return nil
	}
	return &models.Message{Size: int64(m.Size),
		Data:          m.Data,
		UUID:          "",
		QueueId:       qi,
		StorageDriver: "inline",
		Timestamp:     time.Now().UnixNano() / 1000000000}
}
