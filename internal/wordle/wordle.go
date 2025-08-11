package wordle

import (
	"fmt"
	"slices"
)

var LETTERS map[rune]bool = make(map[rune]bool)

const (
	WORD_LEN = 5
	GUESSES  = 6
)

type Greens [WORD_LEN]rune
type Yellows map[rune][]int
type Greys map[rune]bool

type Wordle struct {
	board []Word
	green Greens
	yello Yellows
	greys Greys
}

type Colour int

const (
	EMPTY Colour = iota
	GREY
	YELLOW
	GREEN
)

type Word [WORD_LEN]Char

type Char struct {
	Val    rune
	Colour Colour
}

func NewWordle() Wordle {
	return Wordle{
		green: Greens{' ', ' ', ' ', ' ', ' '},
		yello: make(Yellows, len(LETTERS)),
		greys: make(Greys, len(LETTERS)),
	}
}

func (w *Wordle) AddWord(word Word) error {
	if len(w.board) >= GUESSES {
		return fmt.Errorf("Guess count reached limit of %d", GUESSES)
	}
	for i, c := range word {
		switch c.Colour {
		case EMPTY:
			return fmt.Errorf("Cannot add empty colour chars")
		case GREEN:
			if w.green[i] != ' ' && c.Val != w.green[i] {
				return fmt.Errorf("Invalid green at index %d", i)
			}
			if w.green[i] == ' ' {
				w.green[i] = c.Val
			}
		case GREY:
			w.greys[c.Val] = true
		case YELLOW:
			if !slices.Contains(w.yello[c.Val], i) {
				w.yello[c.Val] = append(w.yello[c.Val], i)
			}
		}
	}
	return nil
}

func init() {
	LETTERS = make(map[rune]bool, 'z'-'a'+1)
	for i := 'a'; i <= 'z'; i++ {
		LETTERS[rune(i)] = true
	}
}
