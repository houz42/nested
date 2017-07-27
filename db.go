package nested

import (
	"database/sql"
	"log"
	"strconv"
)

// SetTableName for query strings
func SetTableName(name string) {
	tblName = name
}

func query(db *sql.DB, sql string, args ...interface{}) (records []map[string]string, err error) {
	rows, err := db.Query(sql, args...)
	if err != nil {
		log.Println("Query failed, sql:", sql, " error:", err)
		return nil, err
	}
	defer rows.Close()
	columns, _ := rows.Columns()
	if len(columns) == 0 {
		return records, nil
	}
	values := make([][]byte, len(columns))
	scanArgs := make([]interface{}, len(columns))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			log.Println("rows.Scan failed, sql:", sql, " error:", err)
			return nil, err
		}
		record := make(map[string]string)
		for i, col := range values {
			if col == nil {
				record[columns[i]] = ""
			} else {
				record[columns[i]] = string(col)
			}
		}
		records = append(records, record)
	}
	return records, nil
}

func atoi(s string) int32 {
	i, _ := strconv.ParseInt(s, 10, 32)
	return int32(i)
}

func itoa(i int32) string {
	return strconv.FormatInt(int64(i), 10)
}

func atoi64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}
