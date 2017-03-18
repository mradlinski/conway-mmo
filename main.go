package main

import (
	"flag"
	"log"
	"net/http"

	"fmt"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:8000", "http service address")

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func initGame() {
	game := NewGame()

	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer c.Close()

		dataChannel := make(chan GameUpdate)
		defer close(dataChannel)

		game.AddUserChannel(dataChannel)
		defer game.RemoveUserChannel(dataChannel)

		cmdChannel := make(chan *UserCommand)
		go func(ch chan *UserCommand) {
			for {
				cmd := UserCommand{}
				err := c.ReadJSON(&cmd)
				if err != nil {
					break
				}

				ch <- &cmd
			}
		}(cmdChannel)
		defer close(cmdChannel)

		for {
			select {
			case data := <-dataChannel:
				fmt.Println("Sending update...")
				err = c.WriteMessage(websocket.BinaryMessage, data)
				if err != nil {
					return
				}
			case cmd := <-cmdChannel:
				game.commandChannel <- cmd
			}
		}
	})

	go game.StartGameLoop()
}

func main() {
	flag.Parse()
	log.SetFlags(0)
	initGame()
	log.Println("Starting...")
	log.Fatal(http.ListenAndServe(*addr, nil))
}
