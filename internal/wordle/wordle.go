package wordle

import (
	"fmt"
	"slices"
)

var LETTERS map[rune]bool = make(map[rune]bool)

const (
	WORD_LEN = 5
	GUESSES  = 5
)

type Greens [WORD_LEN]rune
type Yellows map[rune][]int
type Greys map[rune]bool

type Wordle struct {
	Board []Word
	Green Greens
	Yello Yellows
	Greys Greys
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

func NewWord() Word {
	w := Word{}	
	for i := range w {
		w[i].Colour = EMPTY
		w[i].Val = ' '
	}
	return w
}

func NewWordle() Wordle {
	return Wordle{
		Green: Greens{' ', ' ', ' ', ' ', ' '},
		Yello: make(Yellows, len(LETTERS)),
		Greys: make(Greys, len(LETTERS)),
	}
}

func (w *Wordle) AddWord(word Word) error {
	if len(w.Board) >= GUESSES {
		return fmt.Errorf("Guess count reached limit of %d", GUESSES)
	}
	for i, c := range word {
		switch c.Colour {
		case EMPTY:
			return fmt.Errorf("Cannot add empty colour chars")
		case GREEN:
			if w.Green[i] != ' ' && c.Val != w.Green[i] {
				return fmt.Errorf("Invalid green at index %d", i)
			}
			if w.Green[i] == ' ' {
				w.Green[i] = c.Val
			}
		case GREY:
			w.Greys[c.Val] = true
		case YELLOW:
			if !slices.Contains(w.Yello[c.Val], i) {
				w.Yello[c.Val] = append(w.Yello[c.Val], i)
			}
		}
	}
	w.Board = append(w.Board, word)
	return nil
}

func init() {
	LETTERS = make(map[rune]bool, 'z'-'a'+1)
	for i := 'a'; i <= 'z'; i++ {
		LETTERS[rune(i)] = true
	}
}
