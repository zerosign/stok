package tokenize

// For given statement go string, given list of runes (runeTokens) to split, split
// the statement into several words
//
// - convert string statement into `[]rune`
// - loop through rune in statement (runeIdx, runeValue)
func TokenizeString(
	statement string,
	visitor Visitor,
	runeTokens RuneSet,
) error {

	if len(statement) == 0 {
		return EmptyStatement
	}

	raws := []rune(statement)

	if len(raws) == 0 {
		return EmptyStatement
	}

	cursor := -1

	runeIdx := 0

	// NOTE: this loops { byteIdx, runeValue } = range statement
	// since we don't use byteIdx but want to use rawRuneIdx instead
	//
	for _, runeValue := range statement {

		_, isRuneMatch := runeTokens[runeValue]

		localCursor := cursor

		if !isRuneMatch && cursor == -1 {
			// NOTE: assume cursor == -1 means the cursor being resets
			//
			// assign cursor to runeIdx iffs :
			// - first char in statement and isRuneMatch is false
			// - not first char in statement but isRuneMatch is false & cursor == -1
			//
			cursor = runeIdx

		} else if isRuneMatch && cursor != -1 {

			// NOTE: resets cursor but not reset local cursor
			cursor = -1

			// NOTE: if rune is matched and cursor is not -1
			// cursor already scans some of the words
			//
			// we need to capture last scanned runes from
			// |cursor ... runeIdx|
			//

			word := raws[localCursor:runeIdx]
			visitor.VisitWord(localCursor, runeIdx, string(word))
		}

		// log.Println(
		// 	runeIdx,
		// 	string(runeValue),
		// 	isRuneMatch,
		// 	raws[runeIdx] == runeValue,
		// 	words,
		// )

		runeIdx += 1
	}

	// NOTE: this is for handling leftovers if last token isn't match with runeTokens
	//
	if cursor != -1 {
		word := raws[cursor:runeIdx]
		visitor.VisitWord(cursor, runeIdx, string(word))
	}

	return nil
}
