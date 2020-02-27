package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/context"
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

	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()

	go procMsg()

	http.Handle("/", http.FileServer(http.Dir("./static")))

	// http.HandleFunc("/home.html/", home)

	// http.HandleFunc("/echo", wsEcho)
	http.HandleFunc("/ws", wsMessage)

	// http.HandleFunc("/ws", wSocket)

	// http.Handle("/download/", http.StripPrefix("/download/",
	// 	http.FileServer(http.Dir("download"))))
	// http.Handle("/img/", http.StripPrefix("/img/",
	// 	http.FileServer(http.Dir("img"))))
	// http.Handle("/", http.StripPrefix("/",
	// 	http.FileServer(http.Dir("root"))))

	log.Println("Running...")

	// Gorilla mux
	// go func() {
	// 	err := http.ListenAndServeTLS(":443", "pki/server.crt", "pki/server.key",
	// 		context.ClearHandler(http.DefaultServeMux))
	// 	// err := http.ListenAndServe(":80", nil)
	// 	if err != nil {
	// 		panic("ListenAndServeTLS: " + err.Error())
	// 	}
	// }()

	err := http.ListenAndServe(":80", context.ClearHandler(http.DefaultServeMux))
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
