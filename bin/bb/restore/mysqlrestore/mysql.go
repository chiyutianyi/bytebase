// mysqlrestore is a library for restoring MySQL database schemas and data provided by bytebase.com.
package mysqlrestore

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/bytebase/bytebase/bin/bb/connect"
)

// Restore restores the schema of a MySQL instance.
func Restore(conn *connect.MysqlConnect, sc *bufio.Scanner) error {
	s := ""
	delimiter := false
	for sc.Scan() {
		line := sc.Text()

		execute := false
		switch {
		case s == "" && line == "":
			continue
		case strings.HasPrefix(line, "--"):
			continue
		case line == "DELIMITER ;;":
			delimiter = true
			continue
		case line == "DELIMITER ;" && delimiter:
			delimiter = false
			execute = true
		case strings.HasSuffix(line, ";"):
			s = s + line + "\n"
			if !delimiter {
				execute = true
			}
		default:
			s = s + line + "\n"
			continue
		}
		if execute {
			_, err := conn.DB.Exec(s)
			if err != nil {
				return fmt.Errorf("execute query %q failed: %v", s, err)
			}
			s = ""
		}
	}
	if err := sc.Err(); err != nil {
		return err
	}

	return nil
}
