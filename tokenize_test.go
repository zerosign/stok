package tokenize

import (
	"log"
	"os"
	"reflect"
	"testing"
)

func TestTokenizeWords(t *testing.T) {
	log.SetOutput(os.Stderr)

	statements := []string{
		"hello,    world  test, ジャパンナレッジ, ",
		"",
		",",
		"\t test \t \t \t \t         test     test",
	}

	expectedResults := []interface{}{
		[]string{
			"hello",
			"world",
			"test",
			"ジャパンナレッジ",
		},
		EmptyStatement,
		EmptyStatement,
		[]string{
			"test",
			"test",
			"test",
		},
	}

	runeTokens := RuneSets([]rune{' ', '\t', ','})

	for testIdx, statement := range statements {
		expected := expectedResults[testIdx]

		visitor := WordVisitor()
		err := TokenizeString(statement, visitor, runeTokens)

		switch ex := expected.(type) {
		case error:
			if err != ex {
				t.Fatal("err should be : ", ex, " but got ", err)
			}
			return
		default:
		}

		if err != nil {
			t.Fatal(err)
		}

		words, err := GetWords(visitor)

		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(words, expected) {
			t.Log(words)
			t.Fatal(words, " not equal with ", expected)
		}
	}
}
