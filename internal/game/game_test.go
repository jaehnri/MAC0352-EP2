package game

import (
	"testing"
)

func TestWhoWon(t *testing.T) {
	var tests = []struct {
		table  [tableLength][tableLength]string
		whoWon string
	}{
		// X won
		{
			[tableLength][tableLength]string{
				{emptyCell, emptyCell, X},
				{emptyCell, X, emptyCell},
				{X, emptyCell, emptyCell},
			},
			X,
		},
		{
			[tableLength][tableLength]string{
				{emptyCell, X, emptyCell},
				{emptyCell, X, emptyCell},
				{emptyCell, X, emptyCell},
			},
			X,
		},
		{
			[tableLength][tableLength]string{
				{X, X, X},
				{emptyCell, emptyCell, emptyCell},
				{emptyCell, emptyCell, emptyCell},
			},
			X,
		},
		// O won
		{
			[tableLength][tableLength]string{
				{emptyCell, X, emptyCell},
				{X, X, emptyCell},
				{O, O, O},
			},
			O,
		},
		// nobody won
		{
			[tableLength][tableLength]string{
				{emptyCell, emptyCell, emptyCell},
				{emptyCell, emptyCell, emptyCell},
				{emptyCell, emptyCell, emptyCell},
			},
			emptyCell,
		},
		{
			[tableLength][tableLength]string{
				{emptyCell, X, emptyCell},
				{X, X, emptyCell},
				{O, O, X},
			},
			emptyCell,
		},
	}

	for i, test := range tests {
		game := Game{}
		game.table = test.table
		got := game.WhoWon()
		if test.whoWon != got {
			t.Errorf("%d : Expected: %s, Got: %s", i, test.whoWon, got)
		}
	}

}
