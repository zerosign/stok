package tokenize

import "context"

type Position struct {
	Start, Len int
}

type rawTokenMaps struct {
	indices map[string][]Position
	words   []interface{}
}

type tokenMaps struct {
	indices map[string][]Position
	words   []string
}

func TokenIndexer() Visitor {
	return &tokenMaps{
		indices: make(map[string][]Position),
		words:   make([]string, 0),
	}
}

type DeferredApplyEx = func(context.Context, ...string)

// Collect
type CollectEx = func(context.Context) ([]bool, error)

func (tm *tokenMaps) BulkIntersectWithEx(
	ctx context.Context,
	defferedContainsFn DeferredApplyEx,
	collectFn CollectEx,
	batch int,
) ([]Position, error) {

	intersects := make([]Position, 0)
	wordsLen := len(tm.words)

	for ii := 0; ii < wordsLen; ii += batch {
		lastIdx := ii + batch

		if lastIdx >= wordsLen {
			lastIdx = wordsLen
		}

		defferedContainsFn(ctx, tm.words[ii:lastIdx]...)
	}

	res, err := collectFn(ctx)

	if err != nil {
		return nil, err
	}

	for ii, isIntersect := range res {

		if !isIntersect {
			continue
		}

		word := tm.words[ii]

		intersects = append(intersects, tm.indices[word]...)
	}

	return intersects, nil

}

func (tm *tokenMaps) Intersect(sets []string) []Position {

	intersects := make([]Position, 0)

	for _, str := range sets {
		if pos, ok := tm.indices[str]; ok {
			intersects = append(intersects, pos...)
		}
	}

	return intersects
}

func ReplaceWithPosition(ref string, char rune, positions ...Position) string {

	runes := []rune(ref)

	for _, position := range positions {
		for idx := position.Start; idx <= position.Start+position.Len; idx++ {
			runes[idx] = rune('*')
		}
	}

	return string(runes)
}

func (tm *tokenMaps) VisitWord(start, len int, word string) {
	pos, exists := tm.indices[word]

	if !exists {
		tm.indices[word] = []Position{
			{
				Start: start,
				Len:   len,
			},
		}

		tm.words = append(tm.words, word)

	} else {
		pos = append(pos, Position{
			Start: start,
			Len:   len,
		})
	}
}

func (tm *tokenMaps) Words() []string {
	return tm.words
}

func RawTokenIndexer() Visitor {
	return &rawTokenMaps{
		indices: make(map[string][]Position),
		words:   make([]interface{}, 0),
	}
}

func (tm *rawTokenMaps) VisitWord(start, len int, word string) {
	pos, exists := tm.indices[word]

	if !exists {
		tm.indices[word] = []Position{
			{
				Start: start,
				Len:   len,
			},
		}

		tm.words = append(tm.words, word)

	} else {
		pos = append(pos, Position{
			Start: start,
			Len:   len,
		})
	}
}

func (tm *rawTokenMaps) RawWords() []interface{} {
	return tm.words
}
