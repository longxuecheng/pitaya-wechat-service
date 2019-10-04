package algorithm

import (
	"testing"
)

func TestBuildBinarySearchTree(t *testing.T) {
	tree := NewBinarySearchTree(NewBinaryNode(10))
	tree.Insert(9)
	tree.Insert(8)
	tree.Insert(7)
	tree.Insert(11)
	tree.Insert(22)
	tree.Insert(19)

	maxNode, err := tree.FindMax()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Max node value is %d", maxNode.Element)
	minNode, err := tree.FindMin()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Min node value is %d", minNode.Element)
}
