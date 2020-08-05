package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var ch = make(chan []byte)

func procMsg() {
	var messages [][]byte
	for {
		messages = append(messages, <-ch)

		for i, v := range messages {
			fmt.Printf("i: %v  v: %v\n", i, string(v))
		}
	}
}

func wsMessage(w http.ResponseWriter, r *http.Request) {
	var upgrader = websocket.Upgrader{} // use default options
	c, err := upgrader.Upgrade(w, r, nil)
	defer c.Close()
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	_, message, err := c.ReadMessage()
	if err != nil {
		log.Println("read:", err)
		return
	}

	log.Printf("recv: %s", message)
	ch <- message

	response := []byte("response")
	// mesage type = 1
	err = c.WriteMessage(1, response)
	if err != nil {
		log.Println("ws write err:", err)
		return
	}

	log.Println("ws sent response")
}

func main() {

	go procMsg()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	http.HandleFunc("/ws", wsMessage)

	log.Println("Running...")

	err := http.ListenAndServe(":80", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
