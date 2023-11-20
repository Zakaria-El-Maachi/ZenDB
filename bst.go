package main

type Pair struct {
	marker bool //Deleted or Just Set
	key    string
	value  string
}

type TreeNode struct {
	elem  Pair // Key Value
	color bool //RB Tree Colors
	left  *TreeNode
	right *TreeNode
	size  int //Size of the SubTree
}

func NewTree(list []Pair) *TreeNode {
	return BuildTree(list, true)
}

func BuildTree(list []Pair, c bool) *TreeNode {
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

func (t *TreeNode) insert(p Pair) {
	if t == nil {
		t.elem = p
		t.left = nil
		t.right = nil
		t.size = 1
	}
	if t.elem.key == p.key {
		t.elem = p
	} else if t.elem.key < p.key {
		t.left.insert(p)
	} else {
		t.right.insert(p)
	}
}

func (t *TreeNode) search(key string) *TreeNode {
	if t == nil {
		return nil
	}
	if t.elem.key == key {
		return t
	}
	if t.elem.key < key {
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
	index := 0
	for i := 0; i < len(p)-1; i++ {
		index += len(p[i].key) + len(p[i].value)
	}
	return index + (len(p)-1)*4
}
