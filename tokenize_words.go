package tokenize

import "fmt"

type words struct {
	inner []string
}

func (w *words) VisitWord(start, len int, word string) {
	w.inner = append(w.inner, word)
}

func WordVisitor() Visitor {
	return &words{
		inner: make([]string, 0),
	}
}

func GetWords(visitor Visitor) ([]string, error) {

	switch v := visitor.(type) {
	case *words:
		return v.inner, nil
	default:
		return nil, fmt.Errorf("invalid type of visitor, should be `words`")
	}

}
