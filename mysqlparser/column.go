package mysqlparser

import (
	"fmt"
	"strings"

	"github.com/jxskiss/gopkg/v2/strutil"
	"github.com/pingcap/parser/ast"
)

type Column struct {
	Name    string
	Comment string
	Table   *Table
	Def     *ast.ColumnDef

	config *Config
}

func (c *Column) GoName() string {
	return strutil.ToCamelCase(c.Name)
}

func (c *Column) VarName() string {
	return strutil.ToLowerCamelCase(c.Name)
}

func (c *Column) ListVarName() string {
	return strutil.ToLowerCamelCase(c.Name + "_list")
}

func (c *Column) GetterFuncName() string {
	return fmt.Sprintf("getter_%s_%s", c.Table.TypeName(), c.GoName())
}

func (c *Column) GoTag() string {
	pkTag := ""
	if c.Name == c.Table.PrimaryKey() {
		pkTag = ";primaryKey"
	}
	return fmt.Sprintf("`db:\"%s\" gorm:\"column:%s%s\"`", c.Name, c.Name, pkTag)
}

func (c *Column) GoType() string {
	configType := c.config.getColumnGoType(c)
	if configType != "" {
		return configType
	}
	defaultType := c.defaultGoType()
	if c.Table.PrimaryKey() == c.Name && defaultType == "int" {
		defaultType = "int64"
	}
	return defaultType
}

func (c *Column) DBType() string {
	return c.Def.Tp.String()
}

func (c *Column) defaultGoType() string {
	dbType := c.Def.Tp.String()
	if strings.HasPrefix(dbType, "bigint") {
		return "int64"
	}
	if strings.HasPrefix(dbType, "tinyint") ||
		strings.HasPrefix(dbType, "smallint") ||
		strings.HasPrefix(dbType, "int") {
		return "int"
	}
	if strings.HasPrefix(dbType, "float") ||
		strings.HasPrefix(dbType, "double") {
		return "float64"
	}
	if strings.HasPrefix(dbType, "decimal") {
		return "string"
	}
	if strings.HasPrefix(dbType, "varchar") ||
		strings.HasPrefix(dbType, "tinytext") ||
		strings.HasPrefix(dbType, "text") ||
		strings.HasPrefix(dbType, "mediumtext") ||
		strings.HasPrefix(dbType, "longtext") {
		return "string"
	}
	if strings.HasPrefix(dbType, "datetime") ||
		strings.HasPrefix(dbType, "timestamp") ||
		strings.HasPrefix(dbType, "date") {
		return "time.Time"
	}
	if strings.HasPrefix(dbType, "year") {
		return "int"
	}
	if strings.HasPrefix(dbType, "varbinary") ||
		strings.Contains(dbType, "blob") {
		return "[]byte"
	}
	panic(fmt.Sprintf("unsupported db type %q", dbType))
}

func (c *Column) IsProtobuf() bool {
	if c.config.Column.IsProtobuf != nil {
		return c.config.Column.IsProtobuf(c)
	}
	return false
}

func (c *Column) PBType() string {
	if c.config.Column.PBType != nil {
		return c.config.Column.PBType(c)
	}
	return ""
}

func (c *Column) IsJSON() bool {
	if c.config.Column.IsJSON != nil {
		return c.config.Column.IsJSON(c)
	}
	return false
}

func (c *Column) JSONType() string {
	if !c.IsJSON() {
		return ""
	}
	if c.config.Column.JSONType != nil {
		return c.config.Column.JSONType(c)
	}
	return "gemap.Map"
}
