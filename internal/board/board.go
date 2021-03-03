package board

import (
	"fmt"

	"github.com/abelgoodwin1988/TickTackGo/internal/validate"
	"github.com/rs/zerolog/log"
)

// CreateBoard creates an empty tick tack toe board
func CreateBoard() Board {
	b := make(Board, 9)
	return b
}

// FillSquare fills a square on a board
func (b Board) FillSquare(elem int, s Square) error {
	if elem > len(b) {
		return fmt.Errorf("cannot fill position %d. board only has positions 0-8", elem)
	}
	if b[elem].Value == "" {
		return fmt.Errorf("cannot fill position %d, it is already filled with value %s", elem, b[elem].Value)
	}
	if err := validate.V(s); err != nil {
		return err
	}
	b[elem] = s
	return nil
}

var winConditions = [8][3]int{
	{0, 1, 2}, {3, 4, 5}, {6, 7, 8}, // left to right
	{0, 3, 6}, {1, 4, 7}, {2, 5, 8}, // top to bottom
	{0, 4, 8}, {2, 4, 6}, // cross
}

// WinCondition returns true if the input mark has a win condition
func (b Board) WinCondition(m string) bool {
	for _, o := range winConditions {
		if m == b[o[0]].Value &&
			m == b[o[1]].Value &&
			m == b[o[2]].Value {
			log.Info().Msg("win condition!")
			return true
		}
	}
	log.Info().Msg("no win condition yet")
	return false
}
