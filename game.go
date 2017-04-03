package main

import (
	"math/rand"
	"sync"
	"time"
)

// Game is a structure holding a single instance of Conway MMO
type Game struct {
	board          GameBoard
	users          []*User
	commandChannel chan *UserCommand
	lock           sync.Mutex
}

// GameSize is the size of the game board
const GameSize = 500

// PredefinedCell value representing a server-filled cell (without player color)
const PredefinedCell int = 0

// EmptyCell value representing an empty (dead) cell
const EmptyCell int = -1

// GameBoard type for the game board
type GameBoard [GameSize][GameSize]int

// GameUpdate type for updates sent to users
type GameUpdate []byte

// NewGame creates a new Game instance
func NewGame() *Game {
	return &Game{
		GameBoard{},
		make([]*User, 0, 100),
		make(chan *UserCommand, 100),
		sync.Mutex{},
	}
}

// AddUser adds a user to the game
func (g *Game) AddUser(u *User) {
	g.lock.Lock()
	defer g.lock.Unlock()

	g.users = append(g.users, u)
}

// RemoveUser removes a user from the game
func (g *Game) RemoveUser(u *User) bool {
	g.lock.Lock()
	defer g.lock.Unlock()

	foundIdx := -1
	for idx, u2 := range g.users {
		if u.id == u2.id {
			foundIdx = idx
			break
		}
	}

	if foundIdx == -1 {
		return false
	}

	g.users = append(
		g.users[:foundIdx],
		g.users[foundIdx+1:]...,
	)

	return true
}

func (g *Game) notifyUsers(update GameUpdate) {
	g.lock.Lock()
	defer g.lock.Unlock()

	for _, u := range g.users {
		go u.SendGameUpdate(&update)
	}
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
	payload := cmd.payload
	x := payload.Coords.X
	y := payload.Coords.Y

	if x < 0 || x+UserCommandMaxCellsDim >= GameSize ||
		y < 0 || y+UserCommandMaxCellsDim >= GameSize {
		return
	}

	for i := range payload.Cells {
		for j := range payload.Cells[i] {
			if payload.Cells[i][j] == 1 {
				g.board[x+i][y+j] = cmd.user.color
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
				g.board[x][y] = PredefinedCell
			} else {
				g.board[x][y] = EmptyCell
			}
		}
	}

	ticker := time.NewTicker(time.Millisecond * 500)

	for {
		select {
		case <-ticker.C:
			update := g.calcGameUpdate()
			g.notifyUsers(update)
		case cmd := <-g.commandChannel:
			g.ApplyCommand(cmd)
		}
	}
}
