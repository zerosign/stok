package tokenize

import (
	"log"
	"os"
	"testing"
)

func BenchTokenize(b *testing.B) {
	log.SetOutput(os.Stderr)
}
