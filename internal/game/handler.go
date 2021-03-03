package game

import (
	"bufio"
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

var (
	waitTime    = time.Second * 5
	waiting     = fmt.Sprintf("waiting for second player...")
	setYourMark = "pick a mark between 0-8"
)

// Handle will handle the game between two clients
func (g *Game) Handle() {
	// Set up close/interruption cases from the clients
	// Wait for the game to have two players
	for len(g.Clients) < 2 {
		for _, c := range g.Clients {
			c.Msg(waiting)
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
	for i := 0; !g.WinCondition(g.Clients[int(math.Abs(float64(g.Turn-1)))].Mark); i++ {
		// Say who's turn it is
		log.Info().Int("play", i).
			Str("player", g.Clients[g.Turn].Name).
			Msg("turn")
		fmt.Printf("%s\n", g.Board)
		if err := g.Blast(fmt.Sprintf("%s's turn", g.Clients[g.Turn].Name)); err != nil {
			return errors.Wrap(err, "failed to blast turn msg to clients")
		}
		// Print the board
		if err := g.Blast(fmt.Sprintf("%s", g.Board)); err != nil {
			return errors.Wrap(err, "failed to blast board msg to clients")
		}
		// Ask the player who's turn it is for input
		g.Clients[g.Turn].Msg(setYourMark)
		// Read client's input
		resp, err := bufio.NewReader(g.Clients[g.Turn].Conn).ReadString('\n')
		if err != nil {
			return errors.Wrap(err, "failed to read client input for 'name'")
		}
		position, err := strconv.Atoi(resp[:len(resp)-1])
		if err != nil {
			// iterate and tell the current client to give us a valid number. Ask until they go crazy
		}
		if position < 0 || position > 8 {
			// iterate and tell the current client to give us a valid number. Ask until they go crazy
		}
		if g.Board[position].Value != "" {
			// iterate and tell the current client to give us a valid position. Ask until they go crazy
		}
		log.Info().Int("position", position).Msg("position selected by player")
		g.Board[position].Value = g.Clients[g.Turn].Mark
		log.Info().Str("value", g.Board[position].Value).Str("playermark", g.Clients[g.Turn].Mark)
		g.NextTurn()
	}
	// because we next-turn on end of loop, we need to set it back
	// to the winnder for easy reference
	g.NextTurn()
	// message about the winner, and then close the connections! GG Noobs
	g.Blast(fmt.Sprintf("%s", g.Board))
	g.Blast(fmt.Sprintf("Congratulations to the Winner %s!!!", g.Clients[g.Turn].Name))
	g.GameOver()
	return nil
}
