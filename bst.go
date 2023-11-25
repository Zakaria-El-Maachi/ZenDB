package main

type Color bool

const (
	Red   Color = true
	Black Color = false
)

type Pair struct {
	marker bool //Deleted or Just Set
	key    string
	value  string
}

type TreeNode struct {
	elem  Pair  // Key Value
	color Color //RB Tree Colors
	left  *TreeNode
	right *TreeNode
	size  int //Size of the SubTree
}

func NewTree(list []Pair) *TreeNode {
	return BuildTree(list, false)
}

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

func insert(t **TreeNode, p Pair) int {
	var addedLength int
	*t, addedLength = insertRB(*t, p)
	(*t).color = Black // Ensure the root is always black
	return addedLength
}

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

func isRed(node *TreeNode) bool {
	if node == nil {
		return false
	}
	return node.color == Red
}

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

func flipColors(h *TreeNode) {
	h.color = Red
	h.left.color = Black
	h.right.color = Black
}

func size(node *TreeNode) int {
	if node == nil {
		return 0
	}
	return node.size
}

func (t *TreeNode) search(key string) *TreeNode {
	if t == nil {
		return nil
	}
	if t.elem.key == key {
		return t
	}
	if t.elem.key > key {
		return t.left.search(key)
	} else {
		return t.right.search(key)
	}
}

func (t *TreeNode) traverse() []Pair {
	if t == nil {
		return make([]Pair, 0)
	}
	return append(append(t.left.traverse(), t.elem), t.right.traverse()...)
}

func (t *TreeNode) getMaxOffset(p []Pair) int {
	if len(p) == 0 {
		return 0
	}
	if len(p) == 1 {
		return 14
	}
	index := 14
	for i := 0; i < len(p)-1; i++ {
		if p[i].marker {
			index += len(p[i].key) + 2 + len(p[i].value)
		} else {
			index += len(p[i].key)
		}
	}
	return index + 3*(len(p)-1)
}
