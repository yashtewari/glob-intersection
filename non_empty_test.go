package gintersect

import (
	"testing"
)

var (
	nonEmptyIntersections = map[string][]string{
		"abcd":        []string{"abcd", "....", "[a-d]*"},
		"pqrs":        []string{".qrs", "p.rs", "pq.s", "pqr."},
		".*":          []string{"asdklfj", "jasdfh", "asdhfajfh", "asdflkasdfjl"},
		"d*":          []string{"[abcd][abcd]", "d[a-z]+", ".....", "[d]*"},
		"[a-p]+":      []string{"[p-z]+", "apapapaapapapap", ".*", "abcdefgh*"},
		"abcd[a-c]z+": []string{"abcd[b-d][yz]*", "abcdazzzz", "abcdbzzz", "abcdcz"},
		".*\\\\":      []string{".*"}, // Escaped \ character.
		".a.a":        []string{"b.b.", "c.c.", "d.d.", "e.e."},
	}

	emptyIntersections = map[string][]string{
		"abcd":     []string{"lsdfhda", "abcdla", "asdlfk", "ksdfj"},
		"[a-d]+":   []string{"xyz", "p+", "[e-f]+"},
		"[0-9]*":   []string{"[a-z]", "\\*"},
		"mamama.*": []string{"dadada.*", "nanana.*"},
		".*mamama": []string{".*dadada", ".*nanana"},
		".xyz.":    []string{"paaap", ".*pqr.*"},
	}
)

func TestNonEmptyIntersections(t *testing.T) {
	for lhs, rhss := range nonEmptyIntersections {
		for _, rhs := range rhss {
			ne, err := NonEmpty(lhs, rhs)
			if err != nil {
				// TODO(yash): All errors should show breaking input.
				t.Errorf("error: %v", err)
			}

			if !ne {
				t.Errorf("lhs: %s, rhs: %s should be non-empty", lhs, rhs)
			}
		}
	}
}

func TestEmptyIntersections(t *testing.T) {
	for lhs, rhss := range emptyIntersections {
		for _, rhs := range rhss {
			ne, err := NonEmpty(lhs, rhs)
			if err != nil {
				// TODO(yash): All errors should show breaking input.
				t.Errorf("error: %v", err)
			}

			if ne {
				t.Errorf("lhs: %s, rhs: %s should be non-empty", lhs, rhs)
			}
		}
	}
}
