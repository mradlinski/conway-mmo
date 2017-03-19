package main

import (
	"math/rand"
	"sync"
	"time"
)

// Game is a structure holding a single instance of Conway MMO
type Game struct {
	board             GameBoard
	userNotifChannels []chan *GameUpdate
	commandChannel    chan *UserCommand
	lock              sync.Mutex
}

// GameSize is the size of the game board
const GameSize = 500

// EmptyCell value representing an empty (dead) cell
const EmptyCell int = -1

// GameBoard type for the game board
type GameBoard [GameSize][GameSize]int

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
	Coords Point
	Cells  [UserCommandMaxCellsDim][UserCommandMaxCellsDim]int
	Color  int
}

// NewGame creates a new Game instance
func NewGame() *Game {
	return &Game{
		GameBoard{},
		make([]chan *GameUpdate, 0, 100),
		make(chan *UserCommand, 100),
		sync.Mutex{},
	}
}

// AddUserChannel adds a channel to the list of channels notified on game update
func (g *Game) AddUserChannel(ch chan *GameUpdate) {
	g.lock.Lock()
	g.userNotifChannels = append(g.userNotifChannels, ch)
	g.lock.Unlock()
}

// RemoveUserChannel removes a channel from the list of channels notified on game update
func (g *Game) RemoveUserChannel(ch chan *GameUpdate) bool {
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

func notifyUser(ch chan *GameUpdate, update *GameUpdate) {
	ch <- update
}

func (g *Game) notifyUsers(update GameUpdate) {
	g.lock.Lock()

	for _, ch := range g.userNotifChannels {
		go notifyUser(ch, &update)
	}

	g.lock.Unlock()
}

func countNeighbours(board *GameBoard, x, y int) (int, int) {
	num := 0
	colors := make(map[int]int)
	majorityColor := board[x][y]
	colors[majorityColor] = 1

	for i := -1; i <= 1; i++ {
		for j := -1; j <= 1; j++ {
			nX := (x + i) % GameSize
			nY := (y + j) % GameSize

			if nX < 0 {
				nX = GameSize + nX
			}

			if nY < 0 {
				nY = GameSize + nY
			}

			if !(i == 0 && j == 0) && board[nX][nY] != EmptyCell {
				num++

				color := board[nX][nY]
				colors[color]++

				if colors[color] > colors[majorityColor] {
					majorityColor = color
				}
			}
		}
	}

	return num, majorityColor
}

func (g *Game) calcGameUpdate() GameUpdate {
	update := make([]int, 0)

	oldBoard := g.board
	g.board = GameBoard{}

	for x := range oldBoard {
		for y := range oldBoard[x] {
			neighbours, majorityColor := countNeighbours(&oldBoard, x, y)

			if oldBoard[x][y] != EmptyCell && (neighbours < 2 || neighbours > 3) {
				g.board[x][y] = EmptyCell
			} else if oldBoard[x][y] != EmptyCell && (neighbours == 2 || neighbours == 3) {
				g.board[x][y] = majorityColor
			} else if oldBoard[x][y] == EmptyCell && (neighbours == 3 || neighbours == 6) {
				g.board[x][y] = majorityColor
			} else {
				g.board[x][y] = EmptyCell
			}

			if g.board[x][y] != EmptyCell {
				update = append(update, x, y, g.board[x][y])
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
			if cmd.Cells[i][j] == 1 {
				g.board[x+i][y+j] = cmd.Color
			}
		}
	}
}

// StartGameLoop starts updating the game board and notifying users
func (g *Game) StartGameLoop() {
	for x := range g.board {
		for y := range g.board[x] {
			var r = rand.Intn(100)
			if r > 50 {
				g.board[x][y] = 0
			} else {
				g.board[x][y] = EmptyCell
			}
		}
	}

	ticker := time.NewTicker(time.Millisecond * 500)

	// lastUpdate := time.Now().UnixNano()

	for {
		select {
		case <-ticker.C:
			update := g.calcGameUpdate()
			g.notifyUsers(update)
			// timeNow := time.Now().UnixNano()
			// fmt.Printf("Update after %dms\n", (timeNow-lastUpdate)/1000000)
			// lastUpdate = timeNow
		case cmd := <-g.commandChannel:
			g.ApplyCommand(cmd)
		}
	}
}
