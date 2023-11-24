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

func insert(t **TreeNode, p Pair) {
	if *t == nil {
		*t = &TreeNode{
			elem:  p,
			left:  nil,
			right: nil,
			size:  1,
		}
	} else if (*t).elem.key == p.key {
		(*t).elem = p
	} else if (*t).elem.key > p.key {
		(*t).size++
		insert(&((*t).left), p)
	} else {
		(*t).size++
		insert(&((*t).right), p)
	}
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
