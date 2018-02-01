package gintersect

import (
	"fmt"
	"testing"
)

func TestNonEmptyIntersections(t *testing.T) {
	tests := map[string][]string{
		"abcd":        []string{"abcd", "....", "[a-d]*"},
		"pqrs":        []string{".qrs", "p.rs", "pq.s", "pqr."},
		".*":          []string{"asdklfj", "jasdfh", "asdhfajfh", "asdflkasdfjl"},
		"d*":          []string{"[abcd][abcd]", "d[a-z]+", ".....", "[d]*"},
		"[a-p]+":      []string{"[p-z]+", "apapapaapapapap", ".*", "abcdefgh*"},
		"abcd[a-c]z+": []string{"abcd[b-d][yz]*", "abcdazzzz", "abcdbzzz", "abcdcz"},
		".*\\\\":      []string{".*", "asdfasdf\\\\"}, // Escaped \ character.
		".a.a":        []string{"b.b.", "c.c.", "d.d.", "e.e."},
		".*.*.*.*.*.*.*.*.*.*.*.*.*.*.*": []string{".*.*.*.*.*.*.*.*.*.*.*"},
	}

	for lhs, rhss := range tests {
		for _, rhs := range rhss {
			ne, err := NonEmpty(lhs, rhs)
			if err != nil {
				t.Error(err)
			}

			if !ne {
				t.Errorf("lhs: %s, rhs: %s should be non-empty", lhs, rhs)
			}
		}
	}
}

func TestEmptyIntersections(t *testing.T) {
	tests := map[string][]string{
		"abcd":      []string{"lsdfhda", "abcdla", "asdlfk", "ksdfj"},
		"[a-d]+":    []string{"xyz", "p+", "[e-f]+"},
		"[0-9]*":    []string{"[a-z]", ".\\*"},
		"mamama.*":  []string{"dadada.*", "nanana.*"},
		".*mamama":  []string{".*dadada", ".*nanana"},
		".xyz.":     []string{"paaap", ".*pqr.*"},
		"ab+":       []string{"a", "b", "abc"},
		".*.*.*.*f": []string{".*.*.*.*g"},
	}

	for lhs, rhss := range tests {
		for _, rhs := range rhss {
			ne, err := NonEmpty(lhs, rhs)
			if err != nil {
				t.Error(err)
			}

			if ne {
				t.Errorf("lhs: %s, rhs: %s should be non-empty", lhs, rhs)
			}
		}
	}
}

func BenchmarkContinuousDotStarNonEmpty(b *testing.B) {
	lhs, rhs := "", ""
	dotStar := ".*"
	for i := 1; i <= 15; i++ {
		lhs = lhs + dotStar
		rhs = rhs + dotStar

		b.Run(fmt.Sprintf("with-%d-stars", i), func(b *testing.B) {
			_, err := NonEmpty(lhs, rhs)
			if err != nil {
				b.Error(err)
			}
		})
	}
}

func BenchmarkContinuousDotStarEmpty(b *testing.B) {
	lhsPrefix, rhsPrefix := "", ""
	dotStar := ".*"
	for i := 1; i <= 15; i++ {
		lhsPrefix = lhsPrefix + dotStar
		rhsPrefix = rhsPrefix + dotStar

		lhs, rhs := lhsPrefix+"c", rhsPrefix+"d"

		b.Run(fmt.Sprintf("with-%d-stars", i), func(b *testing.B) {
			_, err := NonEmpty(lhs, rhs)
			if err != nil {
				b.Error(err)
			}
		})
	}
}
