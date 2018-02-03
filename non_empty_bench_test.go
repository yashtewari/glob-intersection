package gintersect

import (
	"fmt"
	"testing"
)

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
