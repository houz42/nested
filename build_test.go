package nested

import (
	"log"
	"runtime/debug"
	"testing"
)

func TestBuild(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			log.Print(string(debug.Stack()))
		}
	}()
	loadAddress()
	trees := buildTrees()
	log.Print("len of beijing areas:", len(trees[0].SubAreas))
	log.Print("len of tianjin areas: ", len(trees[1].SubAreas))
	log.Print("len of hebei cities: ", len(trees[2].SubAreas))
	log.Printf("tree with %d roots", len(trees))

	assignKeys(trees)
	log.Printf("key from %d to %d", trees[0].Left, trees[len(trees)-1].Right)

	genSQLFile(trees)
}

func TestCode(t *testing.T) {
	p := getProvince("120106010000")
	if p != "120000" {
		t.Error(p)
	}
}
