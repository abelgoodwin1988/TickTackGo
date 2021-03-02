package board

import (
	"fmt"

	"github.com/abelgoodwin1988/TickTackGo/internal/validate"
)

// CreateBoard creates an empty tick tack toe board
func CreateBoard() Board {
	b := make(Board, 8)
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
