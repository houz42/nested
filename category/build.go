package category

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strconv"
)

const (
	infoTableName = "category_info"
	treeTableName = "category_tree"
	dataFile      = "./data/categories.json"
	infoFile      = "./data/categoryInfo.sql"
	treeFile      = "./data/categoryTree.sql"
	infoInsertSQL = "INSERT INTO " + infoTableName + "(id, name, introduction, is_delete) VALUES("
	treeInsertSQL = "INSERT INTO " + treeTableName + "(id, name, pid, depth, lft, rgt) VALUES("
)

type category struct {
	Status int32  `json:"status,omitempty"`
	Leaf   int32  `json:"leaf,omitempty"`
	Name   string `json:"name,omitempty"`
	SPUID  int64  `json:"spuid,omitempty"`
	Spell  string `json:"spell,omitempty"`
	SID    int64  `json:"sid,omitempty,string"`
	PID    int64  `json:"pid,omitempty,string"`
	Sub    []*category
	Left   int32
	Right  int32
	Depth  int32
}

func loadTree() *category {
	file, err := os.Open(dataFile)
	if err != nil {
		log.Fatal("os.Open error: ", err)
	}
	defer file.Close()

	// put new node into catMap and Sub filed of its parent
	catMap := make(map[int64]*category)
	var root category
	root.SID = 0
	catMap[0] = &root

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var cat category
		err := json.Unmarshal([]byte(scanner.Text()), &cat)
		if err != nil {
			log.Print("json.Unmarshal error: ", err)
		}
		catMap[cat.SID] = &cat
		p := catMap[cat.PID]
		if p.Sub == nil {
			p.Sub = make([]*category, 0)
		}
		p.Sub = append(p.Sub, &cat)
	}
	log.Printf("got %d categories", len(catMap))
	return &root
}

// number the nodes according a tree traversal
func assignKeys(tree *category) {
	start := int32(0)
	start = indexTree(tree, start)
}

func indexTree(root *category, start int32) int32 {
	start++
	root.Left = start
	for _, sub := range root.Sub {
		start = indexTree(sub, start)
	}
	start++
	root.Right = start
	return start
}

// generate database table initial inserting sql queries
func genSQLFile(root *category) {
	info, err := os.Create(infoFile)
	if err != nil {
		log.Panic("os.Create error: ", err)
	}
	defer info.Close()

	tree, err := os.Create(treeFile)
	if err != nil {
		log.Panic("os.Create error: ", err)
	}
	defer tree.Close()

	for _, p := range root.Sub {
		genSQL(info, tree, p, 1)
	}
}

func genSQL(info, tree *os.File, cat *category, depth int32) {
	sql := bytes.NewBufferString(infoInsertSQL)
	sql.WriteString(i64toa(cat.SID))
	sql.WriteString(", '")
	sql.WriteString(cat.Name)
	sql.WriteString("', '")
	sql.WriteString(cat.Spell)
	sql.WriteString("', ")
	sql.WriteString(itoa(cat.Status))
	sql.WriteString(");\n")

	_, err := info.Write(sql.Bytes())
	if err != nil {
		log.Fatal("info.Write error: ", err, " when writting category: ", *cat)
	}

	sql.Reset()
	sql.WriteString(treeInsertSQL)
	sql.WriteString(i64toa(cat.SID))
	sql.WriteString(", '")
	sql.WriteString(cat.Name)
	sql.WriteString("', ")
	sql.WriteString(i64toa(cat.PID))
	sql.WriteString(", ")
	sql.WriteString(itoa(cat.Depth))
	sql.WriteString(", ")
	sql.WriteString(itoa(cat.Left))
	sql.WriteString(", ")
	sql.WriteString(itoa(cat.Right))
	sql.WriteString(");\n")

	_, err = tree.Write(sql.Bytes())
	if err != nil {
		log.Fatal("tree.Write error: ", err, " when writting category: ", *cat)
	}

	for _, sub := range cat.Sub {
		genSQL(info, tree, sub, depth+1)
	}
}

func i64toa(i int64) string {
	return strconv.FormatInt(i, 10)
}

func itoa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}
