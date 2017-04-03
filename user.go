package main

import (
	"sync"
	"sync/atomic"

	"time"

	"github.com/gorilla/websocket"
)

var userIDCounter int64 = 1

// User represents a connected user
type User struct {
	id            int64
	conn          *websocket.Conn
	connected     bool
	updateChannel chan *GameUpdate
	cmdChannel    chan *UserCommand
	color         int
	lock          sync.Mutex
}

// UserCommand type representing a user's command to alter the board
type UserCommand struct {
	user    *User
	payload *UserCommandPayload
}

// UserCommandMaxCellsDim is the maximum size of the grid sent in a user command
const UserCommandMaxCellsDim = 8

// UserCommandPayload represents the exact command the user sent
type UserCommandPayload struct {
	Coords Point                                               `json:"coords"`
	Cells  [UserCommandMaxCellsDim][UserCommandMaxCellsDim]int `json:"cells"`
}

// NewUser creates a new user
func NewUser(c *websocket.Conn, color int) *User {
	return &User{
		nextUserID(),
		c,
		true,
		make(chan *GameUpdate, 1),
		make(chan *UserCommand, 1),
		color,
		sync.Mutex{},
	}
}

func (u *User) IsConnected() bool {
	u.lock.Lock()
	defer u.lock.Unlock()

	return u.connected
}

func (u *User) Disconnect() error {
	u.lock.Lock()
	defer u.lock.Unlock()

	err := u.conn.Close()
	u.connected = false
	close(u.cmdChannel)
	close(u.updateChannel)

	return err
}

func (u *User) SendGameUpdate(data *GameUpdate) error {
	u.lock.Lock()
	defer u.lock.Unlock()

	u.conn.SetWriteDeadline(time.Now().Add(time.Second / 3))
	err := u.conn.WriteMessage(websocket.BinaryMessage, *data)
	if err != nil {
		return err
	}

	return nil
}

func (u *User) ReadCommand() (*UserCommand, error) {
	payload := UserCommandPayload{}
	err := u.conn.ReadJSON(&payload)

	if err != nil {
		return nil, err
	}

	cmd := UserCommand{
		u,
		&payload,
	}

	return &cmd, nil
}

func nextUserID() int64 {
	return atomic.AddInt64(&userIDCounter, 1)
}
