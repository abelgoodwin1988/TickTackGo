package game

import (
	"math"
	"math/rand"

	"github.com/abelgoodwin1988/TickTackGo/internal/board"
	"github.com/abelgoodwin1988/TickTackGo/internal/client"
	"github.com/pkg/errors"
)

// Codes holds all the active game room codes
var Codes = map[string]*Game{}

// Game is a current running game
type Game struct {
	Code string
	board.Board
	Turn    int
	Clients []*client.Client
}

// NewGame creates a new game
func NewGame() *Game {
	g := &Game{
		Code:  NewCode(),
		Turn:  0,
		Board: board.CreateBoard(),
	}
	Codes[g.Code] = g
	return g
}

// GameOver cleans up a game
func (g *Game) GameOver() {
	for _, c := range g.Clients {
		c.Close("game over. closing server connection.")
	}
}

// NewCode returns a 4 character all uppercase code
// which is used by a second user to join a game
// that another user has already created
func NewCode() string {
	// 65-90 represents the hex value for all uppercase english alphabet characters
	c := []rune{}
	for i := 1; i < 5; i++ {
		c = append(c, rune(rand.Intn(90-65)+65))
	}
	return string(c)
}

// NextTurn changed the Turn value from 0 to 1, or 1 to 0
func (g *Game) NextTurn() {
	g.Turn = int(math.Abs(float64(g.Turn - 1)))
}

// AddClient attempts to add a client to a game. If game already has two clients,
// an error will be returned
func (g *Game) AddClient(c *client.Client) error {
	if len(g.Clients) >= 2 {
		return errors.New("game has max number of clients")
	}
	g.Clients = append(g.Clients, c)
	return nil
}

// Blast sends a message to all clients in a game
func (g *Game) Blast(msg string) error {
	for _, c := range g.Clients {
		if err := c.Msg(msg); err != nil {
			return errors.Wrapf(err, "failed to blast message to client %s", c)
		}
	}
	return nil
}
