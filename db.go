package nested

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
)

const (
	driverName     = "mysql"
	dataSourceName = "root:root@tcp(10.24.248.100:3307)/froad_xmall_test?charset=utf8"
	// maxOpenConns   int
)

var db *sql.DB

func init() {
	var err error
	db, err = sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Panicln("sql.Open error: ", err)
	}
	db.Ping()
}

func query(sql string, args ...interface{}) (records []map[string]string, err error) {
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
