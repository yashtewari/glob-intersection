// Package gintersect provides methods to check whether the intersection of several globs matches a non-empty set of strings.
package gintersect

// Glob represent a glob.
type Glob []Token

// TokenType is the type of a Token.
type TokenType uint

const (
	TTCharacter TokenType = iota
	TTDot
	TTSet
)

// Flag applies to a token.
type Flag uint

const (
	FlagNone = iota
	FlagPlus
	FlagStar
)

// Token is the element that makes up a Glob.
type Token interface {
	Type() TokenType
	Flag() Flag
	SetFlag(Flag)
}

type token struct {
	ttype TokenType
	flag  Flag
}

func (t token) Type() TokenType {
	return t.ttype
}

func (t token) Flag() Flag {
	return t.flag
}

func (t *token) SetFlag(f Flag) {
	t.flag = f
}

// character is a specific rune.
type character struct {
	token
	r rune
}

func NewCharacter(r rune) Token {
	return &character{
		token: token{ttype: TTCharacter},
		r:     r,
	}
}

func (c character) Rune() rune {
	return c.r
}

// dot is any character.
type dot struct {
	token
}

func NewDot() Token {
	return &dot{
		token: token{ttype: TTDot},
	}
}

// set is a set of characters (similar to regexp character class).
type set struct {
	token
	runes map[rune]bool
}

func NewSet(runes []rune) Token {
	m := map[rune]bool{}
	for _, r := range runes {
		m[r] = true
	}
	return &set{
		token: token{ttype: TTSet},
		runes: m,
	}
}

func (s set) Runes() map[rune]bool {
	return s.runes
}
