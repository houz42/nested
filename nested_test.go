package nested

import (
	"database/sql"
	"log"
	"testing"
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open("driverName", "dataSourceName")
	if err != nil {
		log.Panic(err)
	}
}

func TestAddRoot(t *testing.T) {
	err := AddRootNode(db, 1, "Clothing")
	if err != nil {
		t.Error(err)
	}
}

func TestInserting(t *testing.T) {
	err := AddRootNode(db, 1, "Clothing")
	if err != nil {
		t.Error(err)
	}

	err = AddNodeByParent(db, 2, "Men' s", 1)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(db, 3, "Women' s", 2)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(db, 4, "Suits", 2)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(db, 5, "Slacks", 4)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(db, 6, "Jackets", 5)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(db, 7, "Dresses", 3)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeByParent(db, 8, "Evening Gowns", 7)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(db, 9, "Sun Dresses", 8)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(db, 10, "Skirts", 7)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(db, 11, "Blouses", 10)
	if err != nil {
		t.Error(err)
	}
	err = AddNodeBySibling(db, 12, "Shoes", 1)
	if err != nil {
		t.Error(err)
	}
}

func TestQueryDetail(t *testing.T) {
	node, err := GetNodeDetail(db, 5)
	if err != nil {
		t.Error(err)
	}
	log.Print(node)
}

func TestChildren(t *testing.T) {
	nodes, err := GetChildren(db, 3)
	if err != nil {
		t.Error(err)
	}
	log.Print(nodes)
}

func TestDescendants(t *testing.T) {
	nodes, err := GetDescendants(db, 3)
	if err != nil {
		t.Error(err)
	}
	log.Print(nodes)
}

func TestRemoveOneNode(t *testing.T) {
	err := RemoveOneNode(db, 4)
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveSunTree(t *testing.T) {
	err := RemoveNodeAndDescendants(db, 2)
	if err != nil {
		t.Error(err)
	}
}
