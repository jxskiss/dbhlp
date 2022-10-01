package main

import (
	"bytes"
	"fmt"
	"go/format"
	"log"
	"os"
	"text/template"

	parser "github.com/jxskiss/dbhlp/mysqlparser"
)

func generateDAOs(cfg *Config, tables []*parser.Table) {

	dirName := getFileName(cfg.DAOPkg, "")
	err := mkdirIfNotExists(dirName, 0)
	assertNil(err)

	var code []byte
	for _, t := range tables {
		code = generateDAOCode(cfg, t)

		daoPkgName := getBasePkgName(cfg.DAOPkg)
		customDAOFile := getFileName(cfg.DAOPkg, t.Name+"_store.go")
		genDAOFile := getFileName(cfg.DAOPkg, t.Name+"_store_gen.go")

		if !cfg.DisableFormat {
			code, err = formatCode(genDAOFile, code)
			assertNil(err)
		}
		if len(code) == 0 {
			continue
		}

		log.Printf("writing dao file: %s", genDAOFile)
		err = writeFile(genDAOFile, code, 0644)
		assertNil(err)

		touchCustomDAOFile(customDAOFile, daoPkgName, t)
	}
}

var customDAOTmpl = `
package %s

type %sCustomMethods interface {
}
`

func touchCustomDAOFile(filename string, pkgName string, table *parser.Table) {
	if _, err := os.Stat(filename); err == nil || !os.IsNotExist(err) {
		return
	}

	code := []byte(fmt.Sprintf(customDAOTmpl, pkgName, table.VarName()))
	code, _ = format.Source(code)

	log.Printf("writing dao file: %s", filename)
	err := writeFile(filename, code, 0644)
	assertNil(err)
}

func generateDAOCode(cfg *Config, t *parser.Table) []byte {
	var err error
	var buf bytes.Buffer

	pkgName := getBasePkgName(cfg.DAOPkg)
	headerData := map[string]interface{}{
		"PkgName": pkgName,
	}
	if cfg.DAOPkg != cfg.ModelPkg {
		headerData["ModelPkg"] = cfg.ModelPkg
	}
	err = headerTmpl.Execute(&buf, headerData)
	assertNil(err)

	daoMethods := getDAOMethods(cfg, t)
	err = storeTmpl.ExecuteTemplate(&buf, "dao", map[string]interface{}{
		"Table":   t,
		"Methods": daoMethods,
	})
	assertNil(err)

	queries := t.Queries()
	for _, q := range queries {
		var tmpl string
		var data interface{}
		switch q {
		case "Get", "GetWhere", "MGet", "MGetWhere", "Create", "Update":
			tmpl = q
			data = t
		default:
			cq := parser.ParseQuery(t, q)
			if cq.IsMany() {
				tmpl = "customMGet"
			} else {
				tmpl = "customGet"
			}
			data = cq
		}
		err = storeTmpl.ExecuteTemplate(&buf, tmpl, data)
		assertNil(err)
	}

	return buf.Bytes()
}

func getDAOMethods(cfg *Config, t *parser.Table) (methods []string) {
	pkgPrefix := ""
	if cfg.ModelPkg != cfg.DAOPkg {
		modelPkgName := getBasePkgName(cfg.ModelPkg)
		pkgPrefix = modelPkgName + "."
	}

	queries := t.Queries()
	for _, q := range queries {
		switch q {
		case "Get":
			sig := fmt.Sprintf("Get(ctx context.Context, %s int64, opts ...dbhlp.Opt) (*%s%s, error)",
				t.PKVarName(), pkgPrefix, t.TypeName())
			methods = append(methods, sig)
		case "GetWhere":
			sig := fmt.Sprintf("GetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) (*%s%s, error)",
				pkgPrefix, t.TypeName())
			methods = append(methods, sig)
		case "MGet":
			sig := fmt.Sprintf("MGet(ctx context.Context, %sList []int64, opts ...dbhlp.Opt) (%s%sList, error)",
				t.PKVarName(), pkgPrefix, t.TypeName())
			methods = append(methods, sig)
		case "MGetWhere":
			sig := fmt.Sprintf("MGetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) (%s%sList, error)",
				pkgPrefix, t.TypeName())
			methods = append(methods, sig)
		case "Create":
			sig := fmt.Sprintf("Create(ctx context.Context, %s *%s%s, opts ...dbhlp.Opt) error",
				t.VarName(), pkgPrefix, t.TypeName())
			methods = append(methods, sig)
		case "Update":
			sig := fmt.Sprintf("Update(ctx context.Context, %s int64, updates map[string]interface{}, opts ...dbhlp.Opt) error",
				t.PKVarName())
			methods = append(methods, sig)
		default:
			cq := parser.ParseQuery(t, q)
			if cq.IsMany() {
				sig := fmt.Sprintf("%s(ctx context.Context, %s, opts ...dbhlp.Opt) (%s%sList, error)",
					cq.Name, cq.ArgList(), pkgPrefix, t.TypeName())
				methods = append(methods, sig)
			} else {
				sig := fmt.Sprintf("%s(ctx context.Context, %s, opts ...dbhlp.Opt) (*%s%s, error)",
					cq.Name, cq.ArgList(), pkgPrefix, t.TypeName())
				methods = append(methods, sig)
			}
		}
	}
	return
}

// -------- templates -------- //

var storeTmpl = &template.Template{}

func init() {
	mustParse := func(name, text string) {
		template.Must(storeTmpl.New(name).Parse(text))
	}

	mustParse("dao", `
const {{ .Table.TableNameConst }} = "{{ .Table.Name }}"

type {{ .Table.TypeName }}DAO interface {
	{{- range .Methods }}
	{{ . }}
	{{- end }}
	{{ .Table.VarName }}CustomMethods
}

func Get{{ .Table.TypeName }}DAO(conn dbhlp.DBConn) {{ .Table.TypeName }}DAO {
	return &{{ .Table.DaoImplName }}{
		db: conn,
	}
}

type {{ .Table.DaoImplName }} struct {
	db *gorm.DB
}
`)

	mustParse("Get", `
func (p *{{ .DaoImplName }}) Get(ctx context.Context, {{ .PKVarName }} int64, opts ...dbhlp.Opt) (*{{ .PkgPrefix }}{{ .TypeName }}, error) {
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .TableNameConst }}
	var out = &{{ .PkgPrefix }}{{ .TypeName }}{}
	err := conn.WithContext(ctx).Table(tableName).Where("{{ .PrimaryKey }} = ?", {{ .PKVarName }}).First(out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}
`)

	mustParse("GetWhere", `
func (p *{{ .DaoImplName }}) GetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) (*{{ .PkgPrefix }}{{ .TypeName }}, error) {
	params, opts := dbhlp.SplitOpts(paramsAndOpts)
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .TableNameConst }}
	var out = &{{ .PkgPrefix }}{{ .TypeName }}{}
	err := conn.WithContext(ctx).Table(tableName).Where(where, params...).First(out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}
`)

	mustParse("MGet", `
func (p *{{ .DaoImplName }}) MGet(ctx context.Context, {{ .PKVarName }}List []int64, opts ...dbhlp.Opt) ({{ .PkgPrefix }}{{ .TypeName }}List, error) {
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .TableNameConst }}
	var out {{ .PkgPrefix }}{{ .TypeName }}List
	err := conn.WithContext(ctx).Table(tableName).Where("{{ .PrimaryKey }} in (?)", {{ .PKVarName }}List).Find(&out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}
`)
	mustParse("MGetWhere", `
func (p *{{ .DaoImplName }}) MGetWhere(ctx context.Context, where string, paramsAndOpts ...interface{}) ({{ .PkgPrefix }}{{ .TypeName }}List, error) {
	params, opts := dbhlp.SplitOpts(paramsAndOpts)
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .TableNameConst }}
	var out {{ .PkgPrefix }}{{ .TypeName }}List
	err := conn.WithContext(ctx).Table(tableName).Where(where, params...).Find(&out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}
`)

	mustParse("Create", `
func (p *{{ .DaoImplName }}) Create(ctx context.Context, {{ .VarName }} *{{ .PkgPrefix }}{{ .TypeName }}, opts ...dbhlp.Opt) error {
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .TableNameConst }}
	err := conn.WithContext(ctx).Table(tableName).Create({{ .VarName }}).Error
	if err != nil {
		return errors.AddStack(err)
	}
	return nil
}
`)

	mustParse("Update", `
func (p *{{ .DaoImplName }}) Update(ctx context.Context, {{ .PKVarName }} int64, updates map[string]interface{}, opts ...dbhlp.Opt) error {
	if len(updates) == 0 {
		return errors.New("programming error: empty updates map")
	}
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .TableNameConst }}
	err := conn.WithContext(ctx).Table(tableName).Where("{{ .PrimaryKey }} = ?", {{ .PKVarName }}).Updates(updates).Error
	if err != nil {
		return errors.AddStack(err)
	}
	return nil
}
`)

	mustParse("customGet", `
func (p *{{ .Table.DaoImplName }}) {{ .Name }}(ctx context.Context, {{ .ArgList }}, opts ...dbhlp.Opt) (*{{ .Table.PkgPrefix }}{{ .Table.TypeName }}, error) {
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .Table.TableNameConst }}
	var out = &{{ .Table.PkgPrefix }}{{ .Table.TypeName }}{}
	err := conn.WithContext(ctx).Table(tableName).
		Where({{ .Where }}).
		First(out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}
`)

	mustParse("customMGet", `
func (p *{{ .Table.DaoImplName }}) {{ .Name }}(ctx context.Context, {{ .ArgList }}, opts ...dbhlp.Opt) ({{ .Table.PkgPrefix }}{{ .Table.TypeName }}List, error) {
	conn := dbhlp.GetSession(p.db, opts...)
	tableName := {{ .Table.TableNameConst }}
	var out {{ .Table.PkgPrefix }}{{ .Table.TypeName }}List
	err := conn.WithContext(ctx).Table(tableName).
		Where({{ .Where }}).
		Find(&out).Error
	if err != nil {
		return nil, errors.AddStack(err)
	}
	return out, nil
}
`)
}
