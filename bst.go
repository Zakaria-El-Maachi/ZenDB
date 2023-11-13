package main

import "errors"

type Pair struct {
	key   string
	value string
}

type TreeNode struct {
	elem  Pair
	color bool
	marker bool 
	left  *TreeNode
	right *TreeNode
	size  int
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
		t.elem.value = p.value
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

func (t *TreeNode) delete(key string) (*TreeNode, error) {
	found := t.search(key)
	if found == nil {
		return nil, errors.New("Key Not Found")
	}
	found  23 
	return found, nil
}
