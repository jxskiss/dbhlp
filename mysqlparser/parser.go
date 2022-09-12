package mysqlparser

import (
	"fmt"
	"log"

	"github.com/pingcap/tidb/parser"
	"github.com/pingcap/tidb/parser/ast"
	"github.com/pingcap/tidb/parser/test_driver"
)

var _ test_driver.ValueExpr

func ParseTables(sql string, config *Config) ([]*Table, error) {
	stmts, warns, err := parser.New().Parse(sql, config.Charset, config.Collation)
	if err != nil {
		return nil, err
	}
	for _, warn := range warns {
		log.Printf("WARN: %v", warn)
	}

	var tables []*Table
	for _, stmt := range stmts {
		stmt, ok := stmt.(*ast.CreateTableStmt)
		if !ok {
			continue
		}

		table := &Table{
			Name:   stmt.Table.Name.String(),
			Stmt:   stmt,
			config: config,
		}
		for _, opt := range stmt.Options {
			if opt.Tp == ast.TableOptionComment {
				table.Comment = opt.StrValue
			}
		}

		// columns
		for _, col := range stmt.Cols {
			var tmp = &Column{
				Name:   col.Name.Name.String(),
				Table:  table,
				Def:    col,
				config: config,
			}
			for _, opt := range col.Options {
				if opt.Tp == ast.ColumnOptionComment {
					tmp.Comment = opt.StrValue
				}
			}
			table.Columns = append(table.Columns, tmp)
		}

		// make sure a single column primary key is specified
		if table.PrimaryKey() == "" {
			return nil, fmt.Errorf("primary key not specified for table %s", table.Name)
		}

		tables = append(tables, table)
	}

	return tables, nil
}
