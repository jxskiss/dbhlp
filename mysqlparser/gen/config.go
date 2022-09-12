package main

import (
	"github.com/jxskiss/gopkg/v2/confr"
	"github.com/jxskiss/gopkg/v2/easy"

	parser "github.com/jxskiss/dbhlp/mysqlparser"
)

func readConfig(args *Args) *Config {
	cfg := &Config{
		DisableFormat: args.DisableFormat,
	}
	if args.ConfigFile != "" {
		err := confr.Load(cfg, args.ConfigFile)
		assertNil(err)
	}

	if args.DAOPkg != "" {
		cfg.DAOPkg = args.DAOPkg
	}
	if args.ModelPkg != "" {
		cfg.ModelPkg = args.ModelPkg
	}
	if cfg.DAOPkg == "" {
		panic("dao-pkg is neither given in YAML config file nor command line args")
	}
	return cfg
}

type Config struct {
	Charset   string `yaml:"charset"`
	Collation string `yaml:"collation"`
	DAOPkg    string `yaml:"dao_pkg"`
	ModelPkg  string `yaml:"model_pkg"`

	Queries map[string][]string `yaml:"queries"`

	ColumnsConfig map[string]ColumnsConfig `yaml:"columns_config"`

	DisableFormat bool `yaml:"-"`
}

type ColumnsConfig struct {
	BitmapCols   []string          `yaml:"bitmap_cols"`
	BoolCols     []string          `yaml:"bool_cols"`
	JSONCols     map[string]string `yaml:"json_cols"`
	ProtobufCols map[string]string `yaml:"protobuf_cols"`
}

func (c *Config) ToParserConfig() *parser.Config {
	pc := &parser.Config{
		Charset:   c.Charset,
		Collation: c.Collation,
		DAOPkg:    c.DAOPkg,
		ModelPkg:  c.ModelPkg,
	}

	pc.Table.Queries = c.getQueries
	pc.Column.GoType = c.getGoType
	pc.Column.IsBitmap = c.isBitmap
	pc.Column.IsBool = c.isBool
	pc.Column.IsProtobuf = c.isProtobuf
	pc.Column.PBType = c.getProtobufType
	pc.Column.IsJSON = c.isJSON
	pc.Column.JSONType = c.getJSONType

	return pc
}

func (c *Config) getQueries(t *parser.Table) []string {
	return c.Queries[t.Name]
}

func (c *Config) getGoType(col *parser.Column) string {
	return ""
}

func (c *Config) isBitmap(col *parser.Column) bool {
	return easy.Index(c.ColumnsConfig[col.Table.Name].BitmapCols, col.Name) >= 0
}

func (c *Config) isBool(col *parser.Column) bool {
	return easy.Index(c.ColumnsConfig[col.Table.Name].BoolCols, col.Name) >= 0
}

func (c *Config) isProtobuf(col *parser.Column) bool {
	_, ok := c.ColumnsConfig[col.Table.Name].ProtobufCols[col.Name]
	return ok
}

func (c *Config) getProtobufType(col *parser.Column) string {
	return c.ColumnsConfig[col.Table.Name].ProtobufCols[col.Name]
}

func (c *Config) isJSON(col *parser.Column) bool {
	_, ok := c.ColumnsConfig[col.Table.Name].JSONCols[col.Name]
	return ok
}

func (c *Config) getJSONType(col *parser.Column) string {
	return c.ColumnsConfig[col.Table.Name].JSONCols[col.Name]
}
