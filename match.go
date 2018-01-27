package gintersect

import (
	"github.com/pkg/errors"
)

var (
	ErrBadImplementation = errors.New("this logical path is invalid")
)

func Match(t1 Token, t2 Token) (bool, error) {
	var temp Token
	if t1.Type() > t2.Type() {
		temp = t1
		t1 = t2
		t2 = temp
	}

	switch t1.Type() {
	case TTCharacter:
		ch := t1.(character)

		switch t2.Type() {
		case TTCharacter:
			return matchCharacters(ch, t2.(character)), nil
		case TTDot:
			return matchCharacterDot(ch, t2.(dot)), nil
		case TTSet:
			return matchCharacterSet(ch, t2.(set)), nil
		default:
			return false, ErrBadImplementation
		}

	case TTDot:
		d := t1.(dot)

		switch t2.Type() {
		case TTDot:
			return matchDots(d, t2.(dot)), nil
		case TTSet:
			return matchDotSet(d, t2.(set)), nil
		default:
			return false, ErrBadImplementation
		}

	case TTSet:
		switch t2.Type() {
		case TTSet:
			return matchSets(t1.(set), t2.(set)), nil
		default:
			return false, ErrBadImplementation
		}

	default:
		return false, ErrBadImplementation

	}
}

func matchCharacters(a character, b character) bool {
	return a.Rune() == b.Rune()
}

func matchCharacterDot(a character, b dot) bool {
	return true
}

func matchCharacterSet(a character, b set) bool {
	_, ok := b.Runes()[a.Rune()]
	return ok
}

func matchDots(a dot, b dot) bool {
	return true
}

func matchDotSet(a dot, b set) bool {
	return true
}

func matchSets(a set, b set) bool {
	for k, _ := range a.Runes() {
		if _, ok := b.Runes()[k]; ok {
			return true
		}
	}

	return false
}
