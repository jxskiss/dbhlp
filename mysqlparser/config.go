package mysqlparser

import "path/filepath"

type Config struct {
	Charset   string
	Collation string
	DAOPkg    string
	ModelPkg  string

	Table struct {
		TypeName func(t *Table) string
		Queries  func(t *Table) []string
	}

	Column struct {
		GoType     func(c *Column) string
		IsBitmap   func(c *Column) bool
		IsBool     func(c *Column) bool
		IsProtobuf func(c *Column) bool
		PBType     func(c *Column) string
		IsJSON     func(c *Column) bool
		JSONType   func(c *Column) string
	}
}

func (c *Config) getQueries(t *Table) []string {
	if c == nil || c.Table.Queries == nil {
		return c.DefaultQueries(t)
	}

	var result []string
	cfgQueries := c.Table.Queries(t)
	for _, q := range cfgQueries {
		if q == "@default" {
			result = append(result, c.DefaultQueries(t)...)
			continue
		}
		result = append(result, q)
	}
	return result
}

func (c *Config) DefaultQueries(t *Table) []string {
	return []string{
		"Get", "GetWhere",
		"MGet", "MGetWhere",
		"Create", "Update",
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
