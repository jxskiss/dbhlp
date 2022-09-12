package mysqlparser

import (
	"github.com/jxskiss/gopkg/v2/utils/strutil"
	"github.com/pingcap/tidb/parser/ast"
)

type Table struct {
	Name    string
	Comment string
	Columns []*Column
	Stmt    *ast.CreateTableStmt

	config *Config
}

func (t *Table) GetColumn(name string) *Column {
	for _, col := range t.Columns {
		if col.Name == name {
			return col
		}
	}
	return nil
}

func (t *Table) PrimaryKey() string {
	for _, con := range t.Stmt.Constraints {
		if con.Tp == ast.ConstraintPrimaryKey {
			if len(con.Keys) == 1 {
				return con.Keys[0].Column.Name.String()
			}
		}
	}
	return ""
}

func (t *Table) PKFieldName() string {
	return strutil.ToCamelCase(t.PrimaryKey())
}

func (t *Table) PKVarName() string {
	return strutil.ToLowerCamelCase(t.PrimaryKey())
}

func (t *Table) TypeName() string {
	name := t.config.getTableTypeName(t)
	if name != "" {
		return name
	}
	return strutil.ToCamelCase(t.Name)
}

func (t *Table) VarName() string {
	return strutil.ToLowerCamelCase(t.TypeName())
}

func (t *Table) TableNameConst() string {
	return "tableName_" + t.TypeName()
}

func (t *Table) DaoImplName() string {
	return t.VarName() + "DAOImpl"
}

func (t *Table) Queries() []string {
	return t.config.getQueries(t)
}

func (t *Table) PkgPrefix() string {
	return t.config.getModelPkgPrefix()
}
