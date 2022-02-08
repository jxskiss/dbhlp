package mysqlparser

import "path/filepath"

type Config struct {
	Charset   string
	Collation string
	ModelPkg  string
	DAOPkg    string

	Table struct {
		TypeName func(t *Table) string
		Queries  func(t *Table) []string
	}

	Column struct {
		GoType       func(c *Column) string
		IsProtobuf   func(c *Column) bool
		ProtobufType func(c *Column) string
		IsJSON       func(c *Column) bool
		JSONType     func(c *Column) string
	}
}

func (c *Config) getQueries(t *Table) []string {
	if c != nil && c.Table.Queries != nil {
		return c.Table.Queries(t)
	}
	return c.DefaultQueries(t)
}

func (c *Config) DefaultQueries(t *Table) []string {
	return []string{
		"Get", "GetWhere",
		"MGet", "MGetWhere",
		"Update",
	}
}

func (c *Config) getTableTypeName(t *Table) string {
	name := ""
	if c.Table.TypeName != nil {
		name = c.Table.TypeName(t)
	}
	return name
}

func (c *Config) getColumnGoType(col *Column) string {
	typeName := ""
	if c.Column.GoType != nil {
		typeName = c.Column.GoType(col)
	}
	return typeName
}

func (c *Config) getModelPkgPrefix() string {
	pkgPrefix := ""
	if c.ModelPkg != c.DAOPkg {
		modelPkgName := getBasePkgName(c.ModelPkg)
		pkgPrefix = modelPkgName + "."
	}
	return pkgPrefix
}

func getBasePkgName(fullPkgName string) string {
	return filepath.Base(fullPkgName)
}
