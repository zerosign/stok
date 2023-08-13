package tokenize

type Visitor interface {
	VisitWord(
		start, len int,
		word string,
	)
}

type RuneSet map[rune]struct{}

func RuneSets(runes []rune) RuneSet {
	sets := make(RuneSet)

	for _, r := range runes {
		sets[r] = struct{}{}
	}

	return sets
}
