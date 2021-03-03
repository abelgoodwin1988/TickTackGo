package board

import (
	"fmt"
	"strings"
)

// Board is a tick tack toe board
type Board []Square

// Square represents a single square on a tick tack toe board
type Square struct {
	Value string `validate:"oneof:X O  "`
}

var boardRule = strings.Repeat("-", 17)

// String returns the string representation of a tick tack toe board
func (b Board) String() string {
	return fmt.Sprintf(`  %s  |  %s  |  %s  
%s
  %s  |  %s  |  %s  
%s
  %s  |  %s  |  %s  `,
		b[0], b[1], b[2],
		boardRule,
		b[3], b[4], b[5],
		boardRule,
		b[6], b[7], b[8])
}

// String returns a string representation of a sqaure.
// Only thing notable here is that an empty square is string represented
// by a white space
func (s Square) String() string {
	if s.Value == "" {
		return " "
	}
	return s.Value
}
