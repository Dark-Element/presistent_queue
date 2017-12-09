package controllers

import (
	"net/http"
	"presistentQueue/initializers"
	"io"
	"log"
	"fmt"
)

func  IsAlive(w http.ResponseWriter, r *http.Request, registry *initializers.Registry){
	rows, err := registry.Db.Query("SELECT uuid,data, queue_id, size  FROM messages")
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var uuid string
		var data []byte
		var queue_id int
		var size int
		err = rows.Scan(&uuid, &data, &queue_id, &size)
		fmt.Println(uuid)
		fmt.Println(data)
		fmt.Println(queue_id)
		fmt.Println(size)
	}
	io.WriteString(w, "OK")
}