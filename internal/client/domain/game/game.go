package game

import "fmt"

// //////////////////////////////////////////////////////////////////////
// DATA
// //////////////////////////////////////////////////////////////////////

type cellPosition struct {
	i int
	j int
}

const (
	X     = "x"
	O     = "o"
	empty = ""
)

func opositePlayer(player string) string {
	if player == X {
		return O
	}
	return X
}

const (
	Won = iota
	Lost
	Draw
	Playing
)

const playsFirst = X
const tableLength = 3

type Game struct {
	table      [tableLength][tableLength]string
	emptyCells int
	User       string
	Oponent    string
	turn       string
}

// //////////////////////////////////////////////////////////////////////
// METHODS
// //////////////////////////////////////////////////////////////////////

func NewGame(user string) Game {
	game := Game{
		User:       user,
		Oponent:    opositePlayer(user),
		turn:       playsFirst,
		emptyCells: tableLength * tableLength,
	}
	for i := 0; i < tableLength; i++ {
		for j := 0; j < tableLength; j++ {
			game.table[i][j] = empty
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

func (g Game) Play(i int, j int) error {
	return g.play(i, j, g.User)
}

func (g Game) OponentPlayed(i int, j int) error {
	return g.play(i, j, g.Oponent)
}

func (g *Game) play(i int, j int, player string) error {
	if player != g.turn {
		return fmt.Errorf("it's not the turn of %s", player)
	}
	if g.table[i][j] != empty {
		return fmt.Errorf("cell (%d,%d) is not empty", i, j)
	}
	g.emptyCells--
	g.table[i][j] = player
	g.turn = opositePlayer(g.turn)
	return nil
}

func (g Game) State() int {
	userWon := g.findCompletedLine()
	if userWon == g.User {
		return Won
	}
	if userWon == g.Oponent {
		return Lost
	}
	if g.emptyCells == 0 {
		return Draw
	}
	return Playing
}
func (g Game) findCompletedLine() string {
	for i := 0; i < amountLines; i++ {
		hasWon := g.lineHasWon(lines[i])
		if hasWon != empty {
			return hasWon
		}
	}
	return empty
}
func (g Game) lineHasWon(line [tableLength]cellPosition) string {
	firstCell := g.table[line[0].i][line[0].j]
	for _, cell := range line {
		if g.table[cell.i][cell.j] != firstCell {
			return empty
		}
	}
	return firstCell
}

func (g Game) PrintTable() {
	for i, line := range g.table {
		fmt.Printf(" %s | %s | %s \n", line[0], line[1], line[2])
		if i != tableLength-1 {
			fmt.Println("-------------")
		}
	}
}
