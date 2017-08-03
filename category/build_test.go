package category

import (
	"testing"
)

func TestBuild(t *testing.T) {
	tree := loadTree()
	assignKeys(tree)
	genSQLFile(tree)
}
