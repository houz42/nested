package nested

import (
	"log"
	"testing"
)

func TestAddRoot(t *testing.T) {
	err := AddRootNode(1, "Clothing")
	if err != nil {
		t.Error(err)
	}
}

func TestInserting(t *testing.T) {
	err := AddRootNode(1, "Clothing")
	if err != nil {
		t.Error(err)
	}

	err = AddNodeByParent(2, "Men' s", 1)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(3, "Women' s", 2)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(4, "Suits", 2)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(5, "Slacks", 4)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(6, "Jackets", 5)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(7, "Dresses", 3)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(8, "Evening Gowns", 7)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(9, "Sun Dresses", 8)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(10, "Skirts", 7)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(11, "Blouses", 10)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(12, "Shoes", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestQueryDetail(t *testing.T) {
	node, err := GetNodeDetail(5)
	if err != nil {
		t.Error(err)
	}
	log.Print(node)
}

func TestChildren(t *testing.T) {
	nodes, err := GetChildren(3)
	if err != nil {
		t.Error(err)
	}
	log.Print(nodes)
}

func TestDescendants(t *testing.T) {
	nodes, err := GetDescendants(3)
	if err != nil {
		t.Error(err)
	}
	log.Print(nodes)
}

func TestRemoveOneNode(t *testing.T) {
	err := RemoveOneNode(4)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveSunTree(t *testing.T) {
	err := RemoveNodeAndDescendants(2)
	if err != nil {
		t.Error(err)
	}
}
