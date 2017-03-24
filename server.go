package main

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

const maxUserReqBurst = 100

func proxyUserCommands(game *Game, user *User) {
	limiter := make(chan bool, maxUserReqBurst)

	for i := 0; i < maxUserReqBurst; i++ {
		limiter <- true
	}

	go func() {
		for range time.Tick(time.Second) {
			if !user.IsConnected() {
				break
			}

			limiter <- true
		}

		close(limiter)
	}()

	for {
		cmd, err := user.ReadCommand()
		if err != nil {
			break
		}

		select {
		case <-limiter:
			game.commandChannel <- cmd
		default:
			continue
		}
	}
}

func setupGameHandler(game *Game, getColor func() int) {
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		c.SetReadLimit(1000)

		user := NewUser(c, getColor())

		defer user.Disconnect()

		game.AddUser(user)
		defer game.RemoveUser(user)

		proxyUserCommands(game, user)
	})
}

// StartServer starts the game server and listens to user connections
func StartServer() error {
	getColor := initializeColorPicker()

	game := NewGame()
	go game.StartGameLoop()

	setupGameHandler(game, getColor)

	return http.ListenAndServe("0.0.0.0:8000", nil)
}
