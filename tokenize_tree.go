package tokenize

type Edge struct {
	// XXX: should we use []rune instead ?
	label []rune
	next  *Node
}

type Node struct {
	isLeaf bool
	edges  map[rune]Edge
}

func EmptyNode() *Node {
	return &Node{
		isLeaf: false,
		edges:  make(map[rune]Edge),
	}
}

func EmptyLeafNode() *Node {
	return &Node{
		isLeaf: true,
		edges:  make(map[rune]Edge),
	}
}

func Min(lhs, rhs int) int {
	if lhs < rhs {
		return lhs
	}
	return rhs
}

// find the first index in runes where the rune is mismatch
func FindFirstMismatch(word, other []rune) (int, bool) {
	length := Min(len(word), len(other))

	for ii := 1; ii < length; ii++ {
		if word[ii] != other[ii] {
			return ii, true
		}
	}

	return 0, false
}

func EmptyEdge(label string) Edge {
	return Edge{
		label: []rune(label),
		next: &Node{
			isLeaf: true,
			edges:  make(map[rune]Edge),
		},
	}
}

func (n *Node) InsertRaw(raws []rune, next *Node) {
	pivot := raws[0]

	n.edges[pivot] = Edge{
		label: raws,
		next:  next,
	}

}

func (n *Node) Insert(label string, next Node) {
	n.InsertRaw([]rune(label), &next)
}

// returns *Edge so that it could mutate the edge
func (n *Node) EdgeMut(ptr rune) (*Edge, bool) {
	edge, found := n.edges[ptr]

	if !found {
		return nil, false
	}

	// TODO: could we assume the edge always exists ?

	return &edge, true
}

// NOTE: let's make it immutable rather than need to keep track deleted thingy
//
// XXX: dude, let's just rebuilds it since it's rarely being added or removed anyway
type Tree struct {
	root *Node
}

func NewEmptyWordSet() *Tree {
	return &Tree{root: &Node{
		isLeaf: false,
		edges:  make(map[rune]Edge),
	}}
}

func NewWordSet(words ...string) *Tree {
	tree := NewEmptyWordSet()

	for _, word := range words {
		tree.Insert(word)
	}

	return tree
}

func (t *Tree) Search(statement string) (int, bool) {

	node := t.root
	cursor := 0
	raws := []rune(statement)
	rawsLen := len(raws)

	// traverse the node for given prefix
	// shortcircuit if :
	// - prefix doesn't exists
	// -
	for cursor < rawsLen {
		current := raws[cursor]

		_, exists := node.EdgeMut(current)

		// if there's prefix that doesn't exists let just return false with the indexes
		if !exists {
			return cursor, false
		}

		// TODO:
	}

	return cursor, node.isLeaf
}

func (t *Tree) Insert(statement string) {
	node := t.root
	cursor := 0

	raws := []rune(statement)
	rawsLen := len(raws)

	for cursor < rawsLen {
		current := raws[cursor]

		edge, exists := node.EdgeMut(current)
		substr := raws[:cursor]

		// base case: if edge not exists then let's create it and append the new edge then finish
		//
		if !exists {
			// need to handle if edge is not exists for the current pivot
			node.edges[current] = EmptyEdge(string(substr))
			break
		}

		rawLabel := []rune(edge.label)

		// some idioms to make (probably) easier to understand
		//
		// internal word : the word that are being kept by our self
		// external word : the word that we want to insert

		// if edge is exists then we need to prepend to existing edge
		pivot, isMismatch := FindFirstMismatch(
			substr,
			rawLabel,
		)

		// if mismatch then split based on pivot
		//
		// if no mismatch then one of the words contains the other word
		// there are 2 conditions
		//
		// - if internal word that are the longer one, then we need to do the switches between internal & external word,
		//   then insert internal word next to external word
		//
		// - if external word that are the longer one, then we need
		//
		if isMismatch {
			// split label into [prefix|suffix]
			prefix := rawLabel[:pivot]
			suffix := rawLabel[pivot:]

			// assign prefix to old label
			edge.label = prefix

			// doing linked list insert in place shenanigans
			prev := edge.next
			edge.next = EmptyNode()

			// assign suffix for new appended node/edge
			// XXX: TODO refine the API
			edge.next.InsertRaw(suffix, prev)

		} else if !isMismatch && len(substr) == len(rawLabel) {
			// if both string not mismatch and does have same length then both string are equal
			// if internal word & exteral word is equal then let just set current edge next node to be leaf instead
			// this will be reseted anyway if it founds something elses (logic below)
			edge.next.isLeaf = true
			break

		} else if !isMismatch && len(substr) > len(rawLabel) {
			// if internal word is smaller than external word, then set pivot to len label to skips N chars
			pivot = len(edge.label)

		} else if !isMismatch && len(substr) < len(rawLabel) {
			// do linked list insert shenanigans, since internal word should be less than external word (as prefix)

			diffIdx := len(substr)
			suffix := rawLabel[diffIdx:]

			swappedNode := edge.next

			insertedNode := EmptyLeafNode()
			insertedNode.InsertRaw(suffix, swappedNode)

			edge.label = substr
			edge.next = insertedNode

			break
		}

		// iterate to next node
		node = edge.next
		// skips N scanned chars
		cursor += pivot
	}

}
