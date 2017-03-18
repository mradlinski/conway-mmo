package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Game is a structure holding a single instance of Conway MMO
type Game struct {
	board             GameBoard
	userNotifChannels []chan GameUpdate
	commandChannel    chan *UserCommand
	lock              sync.Mutex
}

// GameSize is the size of the game board
const GameSize = 1000

// GameBoard type for the game board
type GameBoard [GameSize][GameSize]byte

// GameUpdate type for updates sent to users
type GameUpdate []byte

// Point represents a point on a grid
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// UserCommandMaxCellsDim is the maximum size of the grid sent in a user command
const UserCommandMaxCellsDim = 8

// UserCommand type representing a user's command to alter the board
type UserCommand struct {
	Coords Point                                               `json:"coords"`
	Cells  [UserCommandMaxCellsDim][UserCommandMaxCellsDim]int `json:"cells"`
}

// NewGame creates a new Game instance
func NewGame() *Game {
	return &Game{
		GameBoard{},
		make([]chan GameUpdate, 0, 100),
		make(chan *UserCommand, 100),
		sync.Mutex{},
	}
}

// AddUserChannel adds a channel to the list of channels notified on game update
func (g *Game) AddUserChannel(ch chan GameUpdate) {
	g.lock.Lock()
	g.userNotifChannels = append(g.userNotifChannels, ch)
	g.lock.Unlock()
}

// RemoveUserChannel removes a channel from the list of channels notified on game update
func (g *Game) RemoveUserChannel(ch chan GameUpdate) bool {
	g.lock.Lock()

	foundIdx := -1
	for idx, otherChannel := range g.userNotifChannels {
		if otherChannel == ch {
			foundIdx = idx
			break
		}
	}

	if foundIdx == -1 {
		return false
	}

	g.userNotifChannels = append(
		g.userNotifChannels[:foundIdx],
		g.userNotifChannels[foundIdx+1:]...,
	)

	g.lock.Unlock()

	return true
}

func notifyUser(ch chan GameUpdate, update GameUpdate) {
	ch <- update
}

func (g *Game) notifyUsers(update GameUpdate) {
	g.lock.Lock()

	for _, ch := range g.userNotifChannels {
		go notifyUser(ch, update)
	}

	g.lock.Unlock()
}

func countNeighbours(board *GameBoard, x, y int) int {
	num := 0

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			nX := x + i
			nY := y + j

			if nX >= 0 && nX < GameSize && nY >= 0 && nY < GameSize && !(i == 0 && j == 0) && board[nX][nY] == 1 {
				num++
			}
		}
	}

	return num
}

func (g *Game) calcGameUpdate() GameUpdate {
	update := make([]int, 100)

	oldBoard := g.board
	g.board = GameBoard{}

	for x := range oldBoard {
		for y := range oldBoard[x] {
			neighbours := countNeighbours(&oldBoard, x, y)

			if oldBoard[x][y] == 1 && (neighbours < 2 || neighbours > 3) {
				g.board[x][y] = 0
			} else if oldBoard[x][y] == 1 && (neighbours == 2 || neighbours == 3) {
				g.board[x][y] = 1
			} else if oldBoard[x][y] == 0 && neighbours == 3 {
				g.board[x][y] = 1
			} else {
				g.board[x][y] = 0
			}

			if g.board[x][y] == 1 {
				update = append(update, x, y)
			}
		}
	}

	return convertIntsToBytes(update)
}

// ApplyCommand applies the user's command to the board
func (g *Game) ApplyCommand(cmd *UserCommand) {
	x := cmd.Coords.X
	y := cmd.Coords.Y

	if x < 0 || x+UserCommandMaxCellsDim >= GameSize ||
		y < 0 || y+UserCommandMaxCellsDim >= GameSize {
		return
	}

	for i := range cmd.Cells {
		for j := range cmd.Cells[i] {
			g.board[x+i][y+j] = byte(cmd.Cells[i][j])
		}
	}
}

// StartGameLoop starts updating the game board and notifying users
func (g *Game) StartGameLoop() {
	for x := range g.board {
		for y := range g.board[x] {
			var r = rand.Intn(100)
			if r > 90 {
				g.board[x][y] = 1
			}
		}
	}

	ticker := time.NewTicker(time.Millisecond * 500)

	lastUpdate := time.Now().UnixNano()

	for {
		select {
		case <-ticker.C:
			update := g.calcGameUpdate()
			g.notifyUsers(update)
			timeNow := time.Now().UnixNano()
			fmt.Printf("Update after %dms\n", (timeNow-lastUpdate)/1000000)
			lastUpdate = timeNow
		case cmd := <-g.commandChannel:
			g.ApplyCommand(cmd)
		}
	}
}
