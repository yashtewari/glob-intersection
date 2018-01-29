package gintersect

import (
	"fmt"

	"github.com/pkg/errors"
)

// Modifier is a special character that affects lexical analysis.
type Modifier uint

const (
	ModifierBackslash Modifier = iota
)

var (
	// Special runes.
	tokenTypeRunes = map[rune]TokenType{
		'.': TTDot,
		'[': TTSet,
		']': TTSet,
	}
	flagRunes = map[rune]Flag{
		'+': FlagPlus,
		'*': FlagStar,
	}
	modifierRunes = map[rune]Modifier{
		'\\': ModifierBackslash,
	}

	// Errors.
	ErrInvalidInput = errors.New("the input provided is invalid")
	errEndOfInput   = errors.New("reached end of input")
)

func Tokenize(input []rune) ([]Token, error) {
	tokens := []Token{}
	for i, t, err := nextToken(0, input); err != errEndOfInput; i, t, err = nextToken(i, input) {
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, t)
	}

	return tokens, nil
}

func nextToken(index int, input []rune) (newIndex int, token Token, err error) {
	var r rune
	var escaped bool

	newIndex, r, escaped, err = nextRune(index, input)
	if err != nil {
		return
	}

	if !escaped {
		if ttype, ok := tokenTypeRunes[r]; ok {
			switch ttype {
			case TTDot:
				token = NewDot()

			case TTSet:
				if r == ']' {
					err = errors.Wrap(ErrInvalidInput, invalidInputMessage(newIndex, "set-close ']' with no preceding '['"))
					return
				}

				newIndex, token, err = nextTokenSet(newIndex, input)

			default:
				err = errors.Wrapf(ErrBadImplementation, "encountered unhandled token type: %v", ttype)
				return
			}

		} else if _, ok := flagRunes[r]; ok {
			err = errors.Wrap(ErrInvalidInput, invalidInputMessage(newIndex, "flag '%s' must be preceded by a non-flag", string(r)))
			return

		} else if m, ok := modifierRunes[r]; ok {
			err = errors.Wrapf(ErrBadImplementation, "encountered unhandled modifier: %v", m)
		} else {
			// Nothing special to do.
			token = NewCharacter(r)
		}
	} else {
		// Nothing special to do.
		token = NewCharacter(r)
	}

	var f Flag
	newIndex, f, err = nextFlag(newIndex, input)
	if err == errEndOfInput {
		// Let this err be passed in the next cycle, after the current token is consumed.
		err = nil
	} else if err != nil {
		return
	}

	token.SetFlag(f)

	return
}

func nextTokenSet(index int, input []rune) (newIndex int, t Token, err error) {
	var r, prev rune
	var escaped bool

	runes := make([]rune, 0, 30)
	complete, prevExists := false, false

	newIndex, r, escaped, err = nextRune(index, input)
	if err != nil {
		return
	}

	for ; !complete && err != errEndOfInput; newIndex, r, escaped, err = nextRune(newIndex, input) {
		if err != nil {
			return
		}

		if !escaped {
			// Handle symbols.
			switch r {
			case '-':
				if !prevExists {
					err = errors.Wrap(ErrInvalidInput, invalidInputMessage(newIndex, "range character '-' must be preceded by a Unicode character"))
					return
				}
				if newIndex >= len(input)-1 {
					err = errors.Wrap(ErrInvalidInput, invalidInputMessage(newIndex, "range character '-' must be followed by a Unicode character"))
					return
				}

				// Get the next rune to know the extent of the range.
				newIndex, r, escaped, err = nextRune(newIndex, input)

				if !escaped {
					if r == ']' || r == '-' {
						err = errors.Wrap(ErrInvalidInput, invalidInputMessage(newIndex, "range character '-' cannot be followed by a special symbol"))
						return
					}
				}
				if r < prev {
					err = errors.Wrap(ErrInvalidInput, invalidInputMessage(newIndex, "range is out of order: '%s' comes before '%s' in Unicode", string(r), string(prev)))
				}

				for x := prev; x <= r; x++ {
					runes = append(runes, x)
				}

				prevExists = false

			case ']':
				complete = true

			// Nothing special to do.
			default:
				runes = append(runes, r)
				prev, prevExists = r, true
			}
		} else {
			// Nothing special to do.
			runes = append(runes, r)
			prev, prevExists = r, true
		}
	}

	if !complete {
		err = errors.Wrap(ErrInvalidInput, invalidInputMessage(newIndex, "found [ without matching ]"))
	} else {
		t = NewSet(runes)
	}

	return
}

func nextFlag(index int, input []rune) (newIndex int, f Flag, err error) {
	var escaped, ok bool
	var r rune

	newIndex, r, escaped, err = nextRune(index, input)
	if err != nil {
		return
	}

	if !escaped {
		// Revert back to index for later consumption.
		if f, ok = flagRunes[r]; !ok {
			newIndex = index
		}
	} else {
		// Revert back to index for later consumption.
		newIndex = index
	}

	return
}

func nextRune(index int, input []rune) (newIndex int, r rune, escaped bool, err error) {
	if index >= len(input) {
		err = errEndOfInput
		return
	}

	if m, ok := modifierRunes[input[index]]; ok {
		switch m {

		case ModifierBackslash:
			if index < len(input)-1 {
				newIndex, r, escaped = index+2, input[index+1], true
			} else if index == len(input)-1 {
				err = errors.Wrap(ErrInvalidInput, invalidInputMessage(index, "input ends with a \\ (escape) character"))
			}
		default:
			err = errors.Wrapf(ErrBadImplementation, "encountered unhandled modifier: %v", m)
		}
	} else {
		newIndex, r, escaped = index+1, input[index], false
	}

	return
}

func invalidInputMessage(index int, str string, a ...interface{}) string {
	return fmt.Sprintf("pos:%d:%s", index, fmt.Sprintf(str, a...))
}
