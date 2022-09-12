package main

import (
	"os"
	"path/filepath"

	"github.com/jxskiss/gopkg/v2/easy"
	"github.com/jxskiss/mcli"

	parser "github.com/jxskiss/dbhlp/mysqlparser"
)

type Args struct {
	SQLFile       string `cli:"#R, -f, -sql-file, SQL file name"`
	ConfigFile    string `cli:"-c, -config-file, YAML config file name"`
	ModelPkg      string `cli:"-model-pkg, full name of generated model package"`
	DAOPkg        string `cli:"#R, -dao-pkg, full name of generated dao package"`
	DisableFormat bool   `cli:"#H, -disable-format, don't format go code"`
}

func main() {
	args := &Args{}
	mcli.Parse(args)

	// TODO: config
	config := &Config{args: args}

	sqlText, err := os.ReadFile(args.SQLFile)
	easy.PanicOnError(err)

	tables, err := parser.ParseTables(string(sqlText), config.ToParserConfig())
	easy.PanicOnError(err)

	generateModels(args, tables)

	generateDAOs(args, tables)
}

func getBasePkgName(fullPkgName string) string {
	return filepath.Base(fullPkgName)
}

func getFileName(pkg, basename string) string {
	goPath := os.Getenv("GOPATH")
	if goPath == "" {
		panic("cannot get GOPATH from env")
	}
	filename := filepath.Join(goPath, "src", pkg, basename)
	return filepath.Clean(filename)
}

func assertNil(args ...interface{}) {
	for _, arg := range args {
		if err, ok := arg.(error); ok && err != nil {
			panic(err)
		}
	}
}
