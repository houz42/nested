package nested

import (
	"bytes"
	"database/sql"
	"errors"
)

var (
	tblName           string
	selectSQL         = "SELECT id, node, pid, depth, lft, rgt FROM " + tblName + " WHERE "
	selectChildrenSQL = "SELECT child.id, child.node, child.pid, child.depth, child.lft, child.rgt FROM " + tblName + " AS child, " + tblName + " AS parent WHERE "
	selectParentsSQL  = "SELECT parent.id, parent.node, parent.pid, parent.depth, parent.lft, parent.rgt FROM " + tblName + " AS child, " + tblName + " AS parent WHERE "
	moveOnAddSQL      = "UPDATE " + tblName + " SET lft=CASE WHEN lft>? THEN lft+2 ELSE lft END, rgt=CASE WHEN rgt>? THEN rgt+2 ELSE rgt END"
	moveOnDeleteSQL   = "UPDATE " + tblName + " SET lft=CASE WHEN lft>? THEN lft-? ELSE lft END, rgt=CASE WHEN rgt>? THEN rgt-? ELSE rgt END"
	moveOnLevelUpSQL  = "UPDATE " + tblName + " SET lft=lft-1, rgt=rgt-1, depth=depth-1 WHERE lft BETWEEN ? AND ?"
	updatePIDSQL      = "UPDATE " + tblName + " AS child, " + tblName + " AS parent SET child.pid=parent.pid WHERE child.pid=parent.id AND child.lft BETWEEN ? AND ?"
	insertSQL         = "INSERT INTO " + tblName + "(id, node, pid, depth, lft, rgt) VALUES(?,?,?,?,?,?)"
	deleteSQL         = "DELETE FROM " + tblName + " WHERE "
)

// Node detail with path from root to node
type Node struct {
	ID          int64
	Node        string
	ParentID    int64
	Depth       int32
	Path        []int64
	PathName    []string
	NumChildren int32
}

// GetNodeDetail with path from root
func GetNodeDetail(db *sql.DB, id int64) (*Node, error) {
	var sql bytes.Buffer
	sql.WriteString(selectParentsSQL)
	sql.WriteString("child.id=? AND child.lft BETWEEN parent.lft AND parent.rgt ORDER BY lft ASC")

	rows, err := query(db, sql.String(), id)
	if err != nil {
		return nil, err
	}
	if len(rows) < 1 {
		return nil, nil
	}

	path := make([]int64, 0, len(rows))
	pathName := make([]string, 0, len(rows))
	for _, r := range rows {
		path = append(path, atoi64(r["id"]))
		pathName = append(pathName, r["node"])
	}

	r := rows[len(rows)-1]
	node := &Node{
		ID:          atoi64(r["id"]),
		Node:        r["node"],
		ParentID:    atoi64(r["pid"]),
		Depth:       atoi(r["depth"]),
		Path:        path,
		PathName:    pathName,
		NumChildren: (atoi(r["rgt"]) - atoi(r["lft"]) - 1) / 2,
	}
	return node, nil
}

// GetChildren returns all immediate children of node
func GetChildren(db *sql.DB, id int64) ([]Node, error) {
	var sql bytes.Buffer
	sql.WriteString(selectSQL)
	sql.WriteString("pid=?")

	rows, err := query(db, sql.String(), id)
	if err != nil {
		return nil, err
	}

	children := make([]Node, 0, len(rows))
	for _, r := range rows {
		children = append(children, Node{
			ID:          atoi64(r["id"]),
			Node:        r["node"],
			ParentID:    atoi64(r["pid"]),
			Depth:       atoi(r["depth"]),
			NumChildren: (atoi(r["rgt"]) - atoi(r["lft"]) - 1) / 2,
		})
	}
	return children, nil
}

// GetDescendants returns sub tree of node
func GetDescendants(db *sql.DB, id int64) ([]Node, error) {
	var sql bytes.Buffer
	sql.WriteString(selectChildrenSQL)
	sql.WriteString("parent.id=? AND child.lft BETWEEN parent.lft AND parent.rgt")

	rows, err := query(db, sql.String(), id)
	if err != nil {
		return nil, err
	}

	descendants := make([]Node, 0, len(rows))
	for _, r := range rows {
		descendants = append(descendants, Node{
			ID:          atoi64(r["id"]),
			Node:        r["node"],
			ParentID:    atoi64(r["pid"]),
			Depth:       atoi(r["depth"]),
			NumChildren: (atoi(r["rgt"]) - atoi(r["lft"]) - 1) / 2,
		})
	}
	return descendants, nil
}

// func GetNodesByDepth(depth int32)([]Node, error)

// AddRootNode adds a new root. There could be more than one root node, and the new root will be the left most one,
// or AddNodeBySibling should be used to insert a new root after another one.
func AddRootNode(db *sql.DB, id int64, name string) error {
	// move all other nodes to right, if exits
	var sql bytes.Buffer
	sql.WriteString(moveOnAddSQL)
	_, err := db.Exec(sql.String(), 0, 0)
	if err != nil {
		return err
	}
	sql.Reset()

	// insert root
	sql.WriteString(insertSQL)
	args := []interface{}{id, name, 0, 1, 1, 2}

	result, err := db.Exec(sql.String(), args...)
	if err != nil {
		return nil
	}
	affected, _ := result.RowsAffected()
	if affected < 1 {
		return errors.New("nested: inserting root affected none")
	}

	return nil
}

// AddNodeByParent adds a new node with certain parent, new node will be the last child of the parent.
func AddNodeByParent(db *sql.DB, id int64, name string, parentID int64) error {
	// query parent
	var sql bytes.Buffer
	sql.WriteString(selectSQL)
	sql.WriteString("id=?")

	rows, err := query(db, sql.String(), parentID)
	if err != nil {
		return err
	}
	if len(rows) < 1 {
		return errors.New("nested: adding node with parent does not exist")
	}
	parentRight := atoi(rows[0]["rgt"])
	parentDepth := atoi(rows[0]["depth"])
	sql.Reset()

	// moves nodes on the right to right by 2,
	sql.WriteString(moveOnAddSQL)

	_, err = db.Exec(sql.String(), parentRight, parentRight-1) //  move right index of parent to right by 2
	if err != nil {
		return err
	}
	sql.Reset()

	// insert new node
	sql.WriteString(insertSQL)
	args := []interface{}{id, name, parentID, parentDepth + 1, parentRight, parentRight + 1}

	r, err := db.Exec(sql.String(), args...)
	if err != nil {
		return err
	}
	row, _ := r.RowsAffected()
	if row != 1 {
		return errors.New("nested: inserting affected none")
	}
	return nil
}

// AddNodeBySibling add a new node right after sibling
func AddNodeBySibling(db *sql.DB, id int64, name string, siblingID int64) error {
	var sql bytes.Buffer

	// query sibling
	sql.WriteString(selectSQL)
	sql.WriteString("id=?")

	rows, err := query(db, sql.String(), siblingID)
	if err != nil {
		return err
	}
	if len(rows) < 1 {
		return errors.New("nested: adding node with sibling does not exist")
	}
	siblingRight := atoi(rows[0]["rgt"])
	siblingDepth := atoi(rows[0]["depth"])
	parentID := atoi(rows[0]["pid"])
	sql.Reset()

	// moves nodes on the right to right by 2
	sql.WriteString(moveOnAddSQL)

	_, err = db.Exec(sql.String(), siblingRight, siblingRight)
	if err != nil {
		return err
	}
	sql.Reset()

	// insert new node
	sql.WriteString(insertSQL)
	args := []interface{}{id, name, parentID, siblingDepth, siblingRight + 1, siblingRight + 2}

	r, err := db.Exec(sql.String(), args...)
	if err != nil {
		return err
	}
	row, _ := r.RowsAffected()
	if row != 1 {
		return errors.New("nested: inserting affected none")
	}
	return nil
}

// RemoveNodeAndDescendants removes node and all its descendants -- it removes the whole subtree.
func RemoveNodeAndDescendants(db *sql.DB, id int64) error {
	// query deleting node
	var sql bytes.Buffer
	sql.WriteString(selectSQL)
	sql.WriteString("id=?")

	rows, err := query(db, sql.String(), id)
	if err != nil {
		return err
	}
	if len(rows) < 1 {
		return errors.New("nested: deleting node does not exist")
	}

	left := atoi(rows[0]["lft"])
	right := atoi(rows[0]["rgt"])
	width := right - left + 1
	sql.Reset()

	// delete node and all its descendants
	sql.WriteString(deleteSQL)
	sql.WriteString("lft BETWEEN ? AND ?")

	result, err := db.Exec(sql.String(), left, right)
	if err != nil {
		return err
	}
	affected, _ := result.RowsAffected()
	if affected < 1 {
		return errors.New("nested: deleting node affected none")
	}
	sql.Reset()

	// move all node on the right to left
	sql.WriteString(moveOnDeleteSQL)

	_, err = db.Exec(sql.String(), right, width, right, width)
	if err != nil {
		return err
	}
	return nil
}

// RemoveOneNode removes one node and move all its descentants 1 level up -- it removes the certain node from the tree only.
func RemoveOneNode(db *sql.DB, id int64) error {
	// query deleting node
	var sql bytes.Buffer
	sql.WriteString(selectSQL)
	sql.WriteString("id=?")

	rows, err := query(db, sql.String(), id)
	if err != nil {
		return err
	}
	if len(rows) < 1 {
		return errors.New("nested: deleting node does not exist")
	}
	sql.Reset()

	left := atoi(rows[0]["lft"])
	right := atoi(rows[0]["rgt"])

	// update pid of its descendants
	sql.WriteString(updatePIDSQL)

	_, err = db.Exec(sql.String(), left, right)
	if err != nil {
		return err
	}
	sql.Reset()

	// delete node
	sql.WriteString(deleteSQL)
	sql.WriteString("id=?")

	r, err := db.Exec(sql.String(), id)
	if err != nil {
		return err
	}
	affected, _ := r.RowsAffected()
	if affected < 1 {
		return errors.New("nested: deleting node affected none")
	}
	sql.Reset()

	// move all its descentants left and up 1 step
	sql.WriteString(moveOnLevelUpSQL)

	_, err = db.Exec(sql.String(), left, right) // could affect none
	if err != nil {
		return err
	}
	sql.Reset()

	// move all other nodes on the right to left by 2 steps
	sql.WriteString(moveOnDeleteSQL)

	_, err = db.Exec(sql.String(), right, 2, right, 2) // could affect none
	if err != nil {
		return err
	}

	return nil
}
