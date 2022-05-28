package game

import (
	"testing"
)

func TestWhoWon(t *testing.T) {
	var tests = []struct {
		table  [tableLength][tableLength]string
		whoWon int
	}{
		// X won
		{
			[tableLength][tableLength]string{
				{empty, empty, X},
				{empty, X, empty},
				{X, empty, empty},
			},
			Won,
		},
		{
			[tableLength][tableLength]string{
				{empty, X, empty},
				{empty, X, empty},
				{empty, X, empty},
			},
			Won,
		},
		{
			[tableLength][tableLength]string{
				{X, X, X},
				{empty, empty, empty},
				{empty, empty, empty},
			},
			Won,
		},
		// O won
		{
			[tableLength][tableLength]string{
				{empty, X, empty},
				{X, X, empty},
				{O, O, O},
			},
			Lost,
		},
		// nobody won
		{
			[tableLength][tableLength]string{
				{empty, empty, empty},
				{empty, empty, empty},
				{empty, empty, empty},
			},
			Playing,
		},
		{
			[tableLength][tableLength]string{
				{empty, X, empty},
				{X, X, empty},
				{O, O, X},
			},
			Playing,
		},
	}

	for i, test := range tests {
		game := NewGame(X)
		game.table = test.table
		got := game.State()
		if test.whoWon != got {
			t.Errorf("%d : Expected: %d, Got: %d", i, test.whoWon, got)
		}
	}

}
