package main

import parser "github.com/jxskiss/dbgen/mysqlparser"

type Config struct {
	args *Args
}

func (c *Config) ToParserConfig() *parser.Config {
	pc := &parser.Config{
		ModelPkg: c.args.ModelPkg,
		DAOPkg:   c.args.DAOPkg,
	}
	return pc
}
