package main

import (
	"go/format"
	"os"
	"path/filepath"

	"github.com/jxskiss/gopkg/v2/easy"
	"github.com/jxskiss/mcli"
	"golang.org/x/tools/imports"

	parser "github.com/jxskiss/dbhlp/mysqlparser"
)

type Args struct {
	SQLFile       string `cli:"#R, -f, -sql-file, SQL file name"`
	ConfigFile    string `cli:"-c, -config-file, YAML config file name"`
	DAOPkg        string `cli:"-dao-pkg, full name of generated dao package, overwrite value in YAML config file"`
	ModelPkg      string `cli:"-model-pkg, full name of generated model package, overwrite value in YAML config file"`
	DisableFormat bool   `cli:"-disable-format, don't format go code"`
}

func main() {
	args := &Args{}
	mcli.Parse(args)
	config := readConfig(args)

	sqlText, err := os.ReadFile(args.SQLFile)
	easy.PanicOnError(err)

	tables, err := parser.ParseTables(string(sqlText), config.ToParserConfig())
	easy.PanicOnError(err)

	generateModels(config, tables)

	generateDAOs(config, tables)
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

func mkdirIfNotExists(dir string, perm os.FileMode) error {
	if perm == 0 {
		perm = 0755
	}
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return os.MkdirAll(dir, perm)
	} else {
		return err
	}
}

func writeFile(name string, data []byte, perm os.FileMode) error {

	//log.Printf("MOCK: writeFile: %v", name)
	//return nil

	return os.WriteFile(name, data, perm)
}

func formatCode(filename string, code []byte) ([]byte, error) {
	code, err := format.Source(code)
	if err != nil {
		return nil, err
	}
	code, err = imports.Process(filename, code, &imports.Options{
		Comments:  true,
		TabIndent: true,
		TabWidth:  4,
	})
	if err != nil {
		return nil, err
	}
	return code, nil
}
