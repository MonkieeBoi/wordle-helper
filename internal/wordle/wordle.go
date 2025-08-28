package wordle

import (
	"fmt"
	"maps"
	"slices"
)

var LETTERS map[rune]bool = make(map[rune]bool)

const (
	WORD_LEN = 5
	GUESSES  = 5
)

type Board []Word
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
	w.board = append(w.board, word)
	return nil
}

func (w Wordle) Board() Board {
	b := make(Board, len(w.board))
	copy(b, w.board)
	return b
}

func (w Wordle) Greens() Greens {
	return w.green
}

func (w Wordle) Yellows() Yellows {
	y := make(Yellows, len(w.yello))
	for k, v := range w.yello {
		s := make([]int, len(v))
		copy(s, v)
		y[k] = s
	}
	return y
}

func (w Wordle) Greys() Greys {
	g := make(Greys, len(w.greys))
	maps.Copy(g, w.greys)
	return g
}

func init() {
	LETTERS = make(map[rune]bool, 'z'-'a'+1)
	for i := 'a'; i <= 'z'; i++ {
		LETTERS[rune(i)] = true
	}
}
