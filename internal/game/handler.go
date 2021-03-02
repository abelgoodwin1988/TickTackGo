package game

import (
	"bufio"
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

var (
	waitTime    = time.Second * 5
	waiting     = []byte(fmt.Sprintf("waiting for second player...\n%s before next check", waitTime))
	setYourMark = []byte("pick a mark between 0-8")
)

// Handle will handle the game between two clients
func (g *Game) Handle() {
	// Set up close/interruption cases from the clients
	// Wait for the game to have two players
	for len(g.Clients) < 2 {
		for _, c := range g.Clients {
			c.Conn.Write(waiting)
		}
		time.Sleep(waitTime)
	}
	if err := g.Play(); err != nil {
		// We need some kind of chan communication for error handling.. we need to set up a
		// a receiver at the creation point of the game type?
	}
}

// Play the game!
func (g *Game) Play() error {
	for g.WinCondition(g.Clients[g.Turn].Mark) {
		// Say who's turn it is
		if err := g.Blast(fmt.Sprintf("%s turn", g.Clients[g.Turn].Name)); err != nil {
			return errors.Wrap(err, "failed to blast turn msg to clients")
		}
		// Print the board
		if err := g.Blast(fmt.Sprint("%s", g.Board)); err != nil {
			return errors.Wrap(err, "failed to blast board msg to clients")
		}
		// Ask the player who's turn it is for input
		g.Clients[g.Turn].Msg(setYourMark)
		// Read client's input
		resp, err := bufio.NewReader(g.Clients[g.Turn].Conn).ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "failed to read client input for 'name'")
		}
		position, err := strconv.Atoi(resp)
		if err != nil {
			// iterate and tell the current client to give us a valid number. Ask until they go crazy
		}
		if position < 0 || position > 8 {
			// iterate and tell the current client to give us a valid number. Ask until they go crazy
		}
		if g.Board[position].Value != "" {
			// iterate and tell the current client to give us a valid position. Ask until they go crazy
		}
		g.Board[position].Value = g.Clients[g.Turn].Mark
		g.NextTurn()
	}
	// because we next-turn on end of loop, we need to set it back
	// to the winnder for easy reference
	g.NextTurn()
	// message about the winner, and then close the connections! GG Noobs
	g.Blast(fmt.Sprintf("%s", g.Board))
	g.Blast(fmt.Sprintf("Congratulations to the Winner %s!!!", g.Clients[g.Turn].Name))
	return nil
}
