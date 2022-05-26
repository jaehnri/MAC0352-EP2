package game

import "fmt"

type cellPosition struct {
	i int
	j int
}

const (
	X         = "x"
	O         = "o"
	emptyCell = ""
)

const tableLength = 3

type Game struct {
	table [tableLength][tableLength]string
}

func NewGame() Game {
	game := Game{}
	for i := 0; i < tableLength; i++ {
		for j := 0; j < tableLength; j++ {
			game.table[i][j] = emptyCell
		}
	}
	return game
}

const amountLines = 8

var lines [amountLines][tableLength]cellPosition = [amountLines][tableLength]cellPosition{
	// Columns
	{{0, 0}, {0, 1}, {0, 2}},
	{{1, 0}, {1, 1}, {1, 2}},
	{{2, 0}, {2, 1}, {2, 2}},
	// Rows
	{{0, 0}, {1, 0}, {2, 0}},
	{{0, 1}, {1, 1}, {2, 1}},
	{{0, 2}, {1, 2}, {2, 2}},
	// Diagonals
	{{0, 0}, {1, 1}, {2, 2}},
	{{2, 0}, {1, 1}, {0, 2}},
}

func (g Game) WhoWon() string {
	for i := 0; i < amountLines; i++ {
		hasWon := g.lineHasWon(lines[i])
		if hasWon != emptyCell {
			return hasWon
		}
	}
	return emptyCell
}
func (g Game) lineHasWon(line [tableLength]cellPosition) string {
	firstCell := g.table[line[0].i][line[0].j]
	for _, cell := range line {
		if g.table[cell.i][cell.j] != firstCell {
			return emptyCell
		}
	}
	return firstCell
}

func (g Game) UserPlayed(i int, j int, user string) error {
	if g.table[i][j] != emptyCell {
		return fmt.Errorf("game: cell (%d,%d) is not empty", i, j)
	}
	g.table[i][j] = user
	return nil
}
