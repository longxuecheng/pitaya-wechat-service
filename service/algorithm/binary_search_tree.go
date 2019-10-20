package algorithm

import "errors"

type BinaryNode struct {
	Element int64
	Left    *BinaryNode
	Right   *BinaryNode
}

func NewBinaryNode(x int64) *BinaryNode {
	return &BinaryNode{
		Element: x,
	}
}

type BinarySearchTree struct {
	Root *BinaryNode
}

func NewBinarySearchTree(node *BinaryNode) *BinarySearchTree {
	return &BinarySearchTree{
		Root: node,
	}
}
func (t *BinarySearchTree) IsEmpty() bool {
	return t.Root == nil
}

func (t *BinarySearchTree) MakeEmpty() {
	t.Root = nil
}

func (t *BinarySearchTree) Contains(x int64) bool {
	return t.contains(x, t.Root)
}

func (t *BinarySearchTree) contains(x int64, node *BinaryNode) bool {
	if node == nil {
		return false
	}
	if x == node.Element {
		return true
	}
	if x < node.Element {
		return t.contains(x, node.Left)
	}
	return t.contains(x, node.Right)
}

func (t *BinarySearchTree) FindMin() (*BinaryNode, error) {
	if t.IsEmpty() {
		return nil, errors.New("Empty tree")
	}
	return t.findMin(t.Root), nil
}

// findMin use recursive method
func (t *BinarySearchTree) findMin(node *BinaryNode) *BinaryNode {
	if node.Left != nil {
		return t.findMin(node.Left)
	}
	return node
}

func (t *BinarySearchTree) FindMax() (*BinaryNode, error) {
	if t.IsEmpty() {
		return nil, errors.New("Empty tree")
	}
	return t.findMax(t.Root), nil
}

// none recursive finding
func (t *BinarySearchTree) findMax(node *BinaryNode) *BinaryNode {
	if node != nil {
		for {
			if node.Right != nil {
				node = node.Right
			} else {
				break
			}
		}
	}
	return node
}

func (t *BinarySearchTree) Insert(x int64) *BinaryNode {
	return t.insert(x, t.Root)
}

func (t *BinarySearchTree) insert(x int64, node *BinaryNode) *BinaryNode {
	if node == nil {
		return &BinaryNode{
			Element: x,
		}
	}
	if x < node.Element {
		node.Left = t.insert(x, node.Left)
	}

	if x > node.Element {
		node.Right = t.insert(x, node.Right)
	}
	return node
}
