package nested

import (
	"bytes"
	"errors"
	"log"
)

const (
	tblName           = "nested"
	selectSQL         = "SELECT id, node, pid, depth, lft, rgt FROM " + tblName + " WHERE "
	selectChildrenSQL = "SELECT child.id, child.node, child.pid, child.depth, child.lft, child.rgt FROM " + tblName + " AS child, " + tblName + " AS parent WHERE "
	selectParentsSQL  = "SELECT parent.id, parent.node, parent.pid, parent.depth, parent.lft, parent.rgt FROM " + tblName + " AS child, " + tblName + " AS parent WHERE "
	moveByParentSQL   = "UPDATE " + tblName + " SET lft=CASE WHEN lft>=? THEN lft+2 ELSE lft END, rgt=CASE WHEN rgt>=? THEN rgt+2 ELSE rgt END"
	moveBySiblingSQL  = "UPDATE " + tblName + " SET lft=CASE WHEN lft>? THEN lft+2 ELSE lft END, rgt=CASE WHEN rgt>? THEN rgt+2 ELSE rgt END"
	moveOnDeleteSQL   = "UPDATE " + tblName + " SET lft=CASE WHEN lft>? THEN lft-? ELSE lft END, rgt=CASE WHEN rgt>? THEN rgt-? ELSE rgt END"
	insertSQL         = "INSERT INTO " + tblName + "(id, node, pid, depth, lft, rgt) VALUES(?,?,?,?,?,?)"
	deleteSQL         = "DELETE FROM " + tblName + " WHERE "
)

// Node detail with path from root to node
type Node struct {
	ID          int32
	Node        string
	ParentID    int32
	Depth       int32
	Path        []int32
	PathName    []string
	NumChildren int32
}

// GetNodeDetail with path from root
func GetNodeDetail(id int32) (*Node, error) {
	log.Println("GetNodeDetail for: ", id)

	var sql bytes.Buffer
	sql.WriteString(selectParentsSQL)
	sql.WriteString("child.id=? AND child.lft BETWEEN parent.lft AND parent.rgt ORDER BY lft ASC")
	log.Println("select parents sql: ", sql.String(), ", args: ", id)

	rows, err := query(sql.String(), id)
	if err != nil {
		log.Panicln("query error: ", err)
	}
	if len(rows) < 1 {
		log.Println("got none")
		return nil, nil
	}

	path := make([]int32, 0, len(rows))
	pathName := make([]string, 0, len(rows))
	for _, r := range rows {
		path = append(path, atoi(r["id"]))
		pathName = append(pathName, r["node"])
	}

	r := rows[len(rows)-1]
	node := &Node{
		ID:          atoi(r["id"]),
		Node:        r["node"],
		ParentID:    atoi(r["pid"]),
		Depth:       atoi(r["depth"]),
		Path:        path,
		PathName:    pathName,
		NumChildren: (atoi(r["rgt"]) - atoi(r["lft"]) - 1) / 2,
	}
	log.Printf("got node detail %+v", *node)
	return node, nil
}

// GetChildren returns all immediate children of node
func GetChildren(id int32) ([]Node, error) {
	log.Println("GetChildren for: ", id)

	var sql bytes.Buffer
	sql.WriteString(selectSQL)
	sql.WriteString("pid=?")
	log.Println("select children sql: ", sql.String(), ", args: ", id)

	rows, err := query(sql.String(), id)
	if err != nil {
		log.Panicln("db.query error: ", err)
	}

	children := make([]Node, 0, len(rows))
	for _, r := range rows {
		children = append(children, Node{
			ID:          atoi(r["id"]),
			Node:        r["node"],
			ParentID:    atoi(r["pid"]),
			Depth:       atoi(r["depth"]),
			NumChildren: (atoi(r["rgt"]) - atoi(r["lft"]) - 1) / 2,
		})
	}
	log.Printf("got children: %+v", children)
	return children, nil
}

// GetDescendants returns sub tree of node
func GetDescendants(id int32) ([]Node, error) {
	log.Println("GetDescendants for: ", id)

	var sql bytes.Buffer
	sql.WriteString(selectChildrenSQL)
	sql.WriteString("parent.id=? AND child.lft BETWEEN parent.lft AND parent.rgt")
	log.Println("select descendants sql: ", sql.String(), ", args: ", id)

	rows, err := query(sql.String(), id)
	if err != nil {
		log.Panic("db.query error: ", err)
	}

	descendants := make([]Node, 0, len(rows))
	for _, r := range rows {
		descendants = append(descendants, Node{
			ID:          atoi(r["id"]),
			Node:        r["node"],
			ParentID:    atoi(r["pid"]),
			Depth:       atoi(r["depth"]),
			NumChildren: (atoi(r["rgt"]) - atoi(r["lft"]) - 1) / 2,
		})
	}
	log.Printf("got descendants: %+v", descendants)
	return descendants, nil
}

// func GetNodesByDepth(depth int32)([]Node, error)

// AddRootNode adds a new root. There could be more than one root node, and the new root will be the left most one,
// or AddNodeBySibling should be used to insert a new root after another one.
func AddRootNode(id int32, name string) error {
	log.Println("AddRootNode for id: ", id, ", name: ", name)

	// move all other nodes to right, if exits
	var sql bytes.Buffer
	sql.WriteString(moveByParentSQL)
	log.Println("move nodes sql: ", sql.String(), ", args: ", 0, 0)
	_, err := db.Exec(sql.String(), 0, 0)
	if err != nil {
		log.Panicln("db.Exec error: ", err)
	}
	sql.Reset()

	// insert root
	sql.WriteString(insertSQL)
	args := []interface{}{id, name, 0, 1, 1, 2}
	log.Println("insert root sql: ", sql.String(), ", args: ", args)

	result, err := db.Exec(sql.String(), args...)
	if err != nil {
		log.Panicln("db.Exec error: ", err)
	}
	affected, _ := result.RowsAffected()
	if affected != 1 {
		return errors.New("insert root affected none")
	}

	log.Println("insert root done")
	return nil
}

// AddNodeByParent adds a new node with certain parent, new node will be the last child of the parent.
func AddNodeByParent(id int32, name string, parentID int32) error {
	log.Println("AddNodeByParent for id: ", id, ", name: ", name, " of parent: ", parentID)

	// query parent
	var sql bytes.Buffer
	sql.WriteString(selectSQL)
	sql.WriteString("id=?")
	log.Println("select parent sql: ", sql.String(), " id: ", parentID)

	rows, err := query(sql.String(), parentID)
	if err != nil {
		log.Panicln("db.query error: ", err)
	}
	if len(rows) < 1 {
		return errors.New("query parent got none")
	}
	parentRight := atoi(rows[0]["rgt"])
	parentDepth := atoi(rows[0]["depth"])
	sql.Reset()

	// moves nodes on the right
	sql.WriteString(moveByParentSQL)
	log.Println("move right sql: ", sql.String(), " on right of: ", parentRight)

	_, err = db.Exec(sql.String(), parentRight, parentRight)
	if err != nil {
		log.Panicln("db.Exec error: ", err)
	}
	sql.Reset()

	// insert new node
	sql.WriteString(insertSQL)
	args := []interface{}{id, name, parentID, parentDepth + 1, parentRight, parentRight + 1}
	log.Println("insert new node sql: ", sql.String(), ", args: ", args)

	r, err := db.Exec(sql.String(), args...)
	if err != nil {
		log.Panicln("db.Exec error: ", err)
	}
	row, _ := r.RowsAffected()
	if row != 1 {
		return errors.New("insert affected none")
	}
	return nil
}

// AddNodeBySibling add a new node right after sibling
func AddNodeBySibling(id int32, name string, siblingID int32) error {
	log.Println("AddNodeBySibling for id: ", id, ", name: ", name, ", with sibling: ", siblingID)

	var sql bytes.Buffer

	// query sibling
	sql.WriteString(selectSQL)
	sql.WriteString("id=?")
	log.Println("select sibling sql: ", sql.String(), " id: ", siblingID)

	rows, err := query(sql.String(), siblingID)
	if err != nil {
		log.Panicln("db.query error: ", err)
	}
	if len(rows) < 1 {
		log.Panicln("query sibling got none: ", siblingID)
	}
	siblingRight := atoi(rows[0]["rgt"])
	siblingDepth := atoi(rows[0]["depth"])
	parentID := atoi(rows[0]["pid"])
	sql.Reset()

	// moves nodes on the right
	sql.WriteString(moveBySiblingSQL)
	log.Println("move right sql: ", sql.String(), " on right of: ", siblingRight)

	_, err = db.Exec(sql.String(), siblingRight, siblingRight)
	if err != nil {
		log.Panicln("db.Exec error: ", err)
	}
	sql.Reset()

	// insert new node
	sql.WriteString(insertSQL)
	args := []interface{}{id, name, parentID, siblingDepth, siblingRight + 1, siblingRight + 2}
	log.Println("insert new node sql: ", sql.String(), ", args: ", args)

	r, err := db.Exec(sql.String(), args...)
	if err != nil {
		log.Panicln("db.Exec error: ", err)
	}
	row, _ := r.RowsAffected()
	if row != 1 {
		return errors.New("insert affected none")
	}
	return nil
}

// RemoveNodes removes node and all its descendants
func RemoveNodes(id string) error {
	log.Println("RemoveNode: ", id)

	// query node
	var sql bytes.Buffer
	sql.WriteString(selectSQL)
	sql.WriteString("id=?")
	log.Println("select sql: ", sql.String(), " id: ", id)

	rows, err := query(sql.String(), id)
	if err != nil {
		log.Panicln("db.query error: ", err)
	}
	if len(rows) < 1 {
		return errors.New("query got none")
	}

	left := atoi(rows[0]["lft"])
	right := atoi(rows[0]["rgt"])
	width := right - left + 1
	sql.Reset()

	// delete descendants
	sql.WriteString(deleteSQL)
	sql.WriteString("lft BETWEEN ? AND ?")
	log.Println("delete sql: ", sql.String(), ", args: ", left, right)

	result, err := db.Exec(sql.String(), left, right)
	if err != nil {
		log.Panicln("db.Exec error: ", err)
	}
	affected, _ := result.RowsAffected()
	if affected < 1 {
		return errors.New("delete affected none")
	}
	sql.Reset()

	// move keys
	sql.WriteString(moveOnDeleteSQL)
	log.Println("move keys on delete sql: ", sql.String(), " args: ", right, width, right, width)

	_, err = db.Exec(sql.String(), right, width, right, width)
	if err != nil {
		log.Panic("db.Exec error: ", err)
	}
	return nil
}
