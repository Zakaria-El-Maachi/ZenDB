package main

// Color represents the color of a node in a Red-Black Tree.
type Color bool

const (
	Red   Color = true
	Black Color = false
)

// Pair represents a key-value pair in the tree.
type Pair struct {
	marker bool   // Deleted or Just Set
	key    string // Key of the pair
	value  string // Value associated with the key
}

// TreeNode represents a node in a Red-Black Tree.
type TreeNode struct {
	elem  Pair  // Key-Value Pair
	color Color // RB Tree Colors
	left  *TreeNode
	right *TreeNode
	size  int // Size of the SubTree
}

// NewTree creates a new Red-Black Tree from a list of pairs.
func NewTree(list []Pair) *TreeNode {
	return BuildTree(list, Black)
}

// BuildTree builds a Red-Black Tree from a list of pairs with a given color.
func BuildTree(list []Pair, c Color) *TreeNode {
	s := len(list)
	if s == 0 {
		return nil
	}
	root := &TreeNode{
		elem:  list[s/2],
		color: c,
		left:  BuildTree(list[:s/2], !c),
		size:  s,
	}
	if s/2+1 < s {
		root.right = BuildTree(list[s/2+1:], !c)
	}
	return root
}

// Insert inserts a new pair into the Red-Black Tree.
func Insert(t **TreeNode, p Pair) int {
	var addedLength int
	*t, addedLength = insertRB(*t, p)
	(*t).color = Black // Ensure the root is always black
	return addedLength
}

// insertRB inserts a new pair into the Red-Black Tree and returns the modified tree and the added length.
func insertRB(root *TreeNode, p Pair) (*TreeNode, int) {
	if root == nil {
		return &TreeNode{
			elem:  p,
			color: Red, // New nodes are always red
			left:  nil,
			right: nil,
			size:  1,
		}, len(p.key) + len(p.value)
	}

	if root.elem.key == p.key {
		diffLength := len(p.value) - len(root.elem.value)
		root.elem = p
		return root, diffLength
	}

	var addedLength int
	if p.key < root.elem.key {
		root.left, addedLength = insertRB(root.left, p)
	} else {
		root.right, addedLength = insertRB(root.right, p)
	}

	if isRed(root.right) && !isRed(root.left) {
		root = rotateLeft(root)
	}
	if isRed(root.left) && isRed(root.left.left) {
		root = rotateRight(root)
	}
	if isRed(root.left) && isRed(root.right) {
		flipColors(root)
	}

	root.size = size(root.left) + size(root.right) + 1
	return root, addedLength
}

// isRed checks if a node is red.
func isRed(node *TreeNode) bool {
	if node == nil {
		return false
	}
	return node.color == Red
}

// rotateLeft performs a left rotation on the Red-Black Tree.
func rotateLeft(h *TreeNode) *TreeNode {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = Red
	x.size = h.size
	h.size = size(h.left) + size(h.right) + 1
	return x
}

// rotateRight performs a right rotation on the Red-Black Tree.
func rotateRight(h *TreeNode) *TreeNode {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = Red
	x.size = h.size
	h.size = size(h.left) + size(h.right) + 1
	return x
}

// flipColors flips the colors of nodes in the Red-Black Tree.
func flipColors(h *TreeNode) {
	h.color = Red
	h.left.color = Black
	h.right.color = Black
}

// size returns the size of the subtree rooted at the given node.
func size(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return node.size
}

// Search searches for a key in the Red-Black Tree.
func (t *TreeNode) Search(key string) *TreeNode {
	if t == nil {
		return nil
	}
	if t.elem.key == key {
		return t
	}
	if t.elem.key > key {
		return t.left.Search(key)
	} else {
		return t.right.Search(key)
	}
}

// Traverse performs an in-order traversal of the Red-Black Tree and returns a list of pairs.
func (t *TreeNode) Traverse() []Pair {
	if t == nil {
		return make([]Pair, 0)
	}
	return append(append(t.left.Traverse(), t.elem), t.right.Traverse()...)
}

// Bloom Filter from a pair slice
func GetBloom(p []Pair) *BloomFilter {
	bloom := NewBloomFilter(29, 10)
	for _, element := range p {
		bloom.Add([]byte(element.key))
	}
	return bloom
}
