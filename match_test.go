package gintersect

import (
	"testing"
)

var (
	matching, nonMatching map[Token][]Token
)

func init() {
	initializeTestSamples()

	matching = map[Token][]Token{
		testCharacters['a']: []Token{testCharacters['a'], testLowerAlphaSet, testLowerAlphaSetPlus, testDot},
		testCharacters['P']: []Token{testUpperAlphaSetStar, testDot},
		testDotPlus:         []Token{testDotStar, testSymbolSet, testNumSetPlus},
		testSymbolSet:       []Token{testCharacters['.'], testCharacters['+'], NewSet([]rune{'.', 'x'})},
		testNumSetPlus:      []Token{testCharacters['0'], testCharacters['9'], testDotStar, NewSet([]rune{'~', 'T', '4'})},
	}

	nonMatching = map[Token][]Token{
		testCharacters['d']: []Token{testCharacters['D'], testCharacters['b'], testNumSet},
		testNumSetPlus:      []Token{testCharacters['.'], testCharacters['g'], testSymbolSetPlus, testLowerAlphaSet},
		testUpperAlphaSet:   []Token{testCharacters['5'], testCharacters['j'], testSymbolSetStar, testLowerAlphaSetPlus},
	}
}

func TestMatching(t *testing.T) {
	for t1, t2s := range matching {
		for _, t2 := range t2s {
			matches, err := Match(t1, t2)
			if err != nil {
				t.Errorf("matching %s and %s gives error: %v", t1.String(), t2.String(), err)
			}

			if !matches {
				t.Errorf("expected %s and %s to match, but they didn't", t1.String(), t2.String())
			}
		}
	}
}

func TestNonMatching(t *testing.T) {
	for t1, t2s := range nonMatching {
		for _, t2 := range t2s {
			matches, err := Match(t1, t2)
			if err != nil {
				t.Errorf("matching %s and %s gives error: %v", t1.String(), t2.String(), err)
			}

			if matches {
				t.Errorf("expected %s and %s not to match, but they did", t1.String(), t2.String())
			}
		}
	}
}
