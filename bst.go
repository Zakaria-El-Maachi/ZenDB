package main

type Pair struct {
	key   string
	value string
}

type TreeNode struct {
	elem   Pair // Key Value
	color  bool //RB Tree Colors
	marker bool //Deleted or Just Set
	left   *TreeNode
	right  *TreeNode
	size   int //Size of the SubTree
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
		elem:   list[s/2],
		color:  c,
		marker: true,
		left:   BuildTree(list[:s/2], !c),
		size:   s,
	}
	if s/2+1 < s {
		root.right = BuildTree(list[s/2+1:], !c)
	}
	return root
}

func (t *TreeNode) insert(p Pair, marker bool) {
	if t == nil {
		t.elem = p
		t.marker = marker
		t.left = nil
		t.right = nil
		t.size = 1
	}
	if t.elem.key == p.key {
		t.elem.value = p.value
		t.marker = marker
	} else if t.elem.key < p.key {
		t.left.insert(p, marker)
	} else {
		t.right.insert(p, marker)
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
