package gintersect

func NonEmpty(lhs string, rhs string) (bool, error) {
	g1, err := NewGlob(lhs)
	if err != nil {
		return false, err
	}

	g2, err := NewGlob(rhs)
	if err != nil {
		return false, err
	}
	// TODO(yash): Go over flow of error messages to see if they make sense.

	return intersectNormal(g1, g2), nil
}

func intersectNormal(g1, g2 Glob) bool {
	var i, j int
	for i, j = 0, 0; i < len(g1) && j < len(g2); i, j = i+1, j+1 {
		if g1[i].Flag() == FlagNone && g2[j].Flag() == FlagNone {
			if !Match(g1[i], g2[j]) {
				return false
			}
		} else {
			return intersectSpecial(g1[i:], g2[j:])
		}
	}

	if i == len(g1) && j == len(g2) {
		return true
	}

	return false
}

func intersectSpecial(g1, g2 Glob) bool {
	if g1[0].Flag() != FlagNone { // If g1 starts with a Token having a Flag.
		switch g1[0].Flag() {
		case FlagPlus:
			return intersectPlus(g1, g2)
		case FlagStar:
			return intersectStar(g1, g2)
		}
	} else { // If g2 starts with a Token having a Flag.
		switch g2[0].Flag() {
		case FlagPlus:
			return intersectPlus(g2, g1)
		case FlagStar:
			return intersectStar(g2, g1)
		}
	}

	return false
}

func intersectPlus(plussed, other Glob) bool {
	if !Match(plussed[0], other[0]) {
		return false
	}
	return intersectStar(plussed, other[1:])
}

func intersectStar(starred, other Glob) bool {
	// starToken, nextToken are the token having FlagStar and the one that follows immediately after, respectively.
	var starToken, nextToken Token

	starToken = starred[0]
	if len(starred) > 1 {
		nextToken = starred[1]
	}

	for i, t := range other {
		// Start gobbling up tokens in other while they match starToken.
		if nextToken != nil && Match(t, nextToken) {
			// When a token in other matches the token after starToken, stop gobbling and try to match the two all the way.
			allTheWay := intersectNormal(starred[1:], other[i:])
			// If they match all the way, the Globs intersect.
			if allTheWay {
				return true
			} else {
				// If they don't match all the way, then the current token from other should still match starToken.
				if !Match(t, starToken) {
					return false
				}
			}
		} else {
			// Only move forward if this token can be gobbled up by starToken.
			if !Match(t, starToken) {
				return false
			}
		}
	}

	// If there was no token following starToken, and everything from other was gobbled, the Globs intersect.
	if nextToken == nil {
		return true
	}

	// If everything from other was gobbles but there was a nextToken to match, they don't intersect.
	return false
}
