package main

import (
	"net/http"

	"github.com/gorilla/websocket"
)

// UserCommandMsg type representing a user's command to alter the board, as it came directly from the user
type UserCommandMsg struct {
	Coords Point                                               `json:"coords"`
	Cells  [UserCommandMaxCellsDim][UserCommandMaxCellsDim]int `json:"cells"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func commandListener(c *websocket.Conn, ch chan *UserCommand, userColor int) {
	for {
		cmd := UserCommandMsg{}
		err := c.ReadJSON(&cmd)
		if err != nil {
			break
		}

		cmdWithColor := UserCommand{
			Coords: cmd.Coords,
			Cells:  cmd.Cells,
			Color:  userColor,
		}

		ch <- &cmdWithColor
	}
}

func userUpdateLoop(c *websocket.Conn, ch chan *GameUpdate) {
	for data := range ch {
		err := c.WriteMessage(websocket.BinaryMessage, *data)
		if err != nil {
			return
		}
	}
}

func setupGameHandler(game *Game, getColor func() int) {
	http.HandleFunc("/connect", func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}

		defer c.Close()

		userColor := getColor()
		go commandListener(c, game.commandChannel, userColor)

		dataChannel := make(chan *GameUpdate)
		defer close(dataChannel)

		game.AddUserChannel(dataChannel)
		defer game.RemoveUserChannel(dataChannel)

		userUpdateLoop(c, dataChannel)
	})
}

// StartServer starts the game server and listens to user connections
func StartServer() error {
	getColor := initializeColorPicker()

	game := NewGame()
	go game.StartGameLoop()

	setupGameHandler(game, getColor)

	return http.ListenAndServe("localhost:8000", nil)
}
