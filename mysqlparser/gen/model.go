package main

import (
	"bytes"
	"go/format"
	"log"
	"text/template"

	parser "github.com/jxskiss/dbhlp/mysqlparser"
)

func generateModels(cfg *Config, tables []*parser.Table) {
	var code []byte

	dirName := getFileName(cfg.ModelPkg, "")
	err := mkdirIfNotExists(dirName, 0)
	assertNil(err)

	for _, t := range tables {
		code = generateModelCode(cfg, t)
		if len(code) == 0 {
			continue
		}

		modelFile := getFileName(cfg.ModelPkg, t.Name+"_model_gen.go")
		log.Printf("writing model file: %s", modelFile)
		err = writeFile(modelFile, code, 0644)
		assertNil(err)
	}
}

func generateModelCode(cfg *Config, table *parser.Table) []byte {
	var err error
	var buf bytes.Buffer

	pkgName := getBasePkgName(cfg.ModelPkg)
	err = headerTmpl.Execute(&buf, map[string]interface{}{
		"PkgName": pkgName,
	})
	assertNil(err)

	err = modelTmpl.ExecuteTemplate(&buf, "model", table)
	assertNil(err)

	err = modelTmpl.ExecuteTemplate(&buf, "getterSetter", table)
	assertNil(err)

	code := buf.Bytes()
	if !cfg.DisableFormat {
		code, err = format.Source(code)
		assertNil(err)
	}
	return code
}

// -------- templates -------- //

var modelTmpl = &template.Template{}

func init() {
	mustParse := func(name, text string) {
		template.Must(modelTmpl.New(name).Parse(text))
	}

	mustParse("model", `
type {{ .TypeName }} struct {

{{- range .Columns }}
	{{ .GoName }} {{ .GoType }} {{ .GoTag }} // {{ .DBType }}
{{- end }}
}

type {{ .TypeName }}List []*{{ .TypeName }}

func (p {{ .TypeName }}List) To{{ .PKFieldName }}Map() map[int64]*{{ .TypeName }} {
	out := make(map[int64]*{{ .TypeName }}, len(p))
	for _, x := range p {
		out[x.{{ .PKFieldName }}] = x
	}
	return out
}

func (p {{ .TypeName }}List) Pluck{{ .PKFieldName }}s() []int64 {
	out := make([]int64, 0, len(p))
	for _, x := range p {
		out = append(out, x.{{ .PKFieldName }})
	}
	return out
}
`)

	mustParse("getterSetter", `
{{ $table := . }}
{{ range .Columns }}

{{ if .IsProtobuf }}
func {{ .GetterFuncName }}(buf []byte) (interface{}, error) {
	var err error
	x := &{{ .PBType }}{}
	if len(buf) > 0 {
		err = proto.Unmarshal(buf, x)
	}
	return x, err
}

func (p *{{ $table.TypeName }}) Get{{ .GoName }}() (*{{ .PBType }}, error) {
	out, err := p.{{ .GoName }}.Get({{ .GetterFuncName }})
	if err != nil {
		log.Printf("failed unmarshal {{ .PBType }}, {{ $table.PrimaryKey }}= %v, err= %v", p.{{ $table.PKFieldName }}, err)
		return nil, err
	}
	return out.(*{{ .PBType }}), nil
}

func (p *{{ $table.TypeName }}) Set{{ .GoName }}({{ .VarName }} *{{ .PBType }}) {
	buf, _ := proto.Marshal({{ .VarName }})
	p.{{ .GoName }}.Set(buf, {{ .VarName }})
}
{{ end }}

{{ if .IsJSON }}
func {{ .GetterFuncName }}(buf []byte) (interface{}, error) {
	var err error
	x := &{{ .JSONType }}{}
	if len(buf) > 0 {
		err = json.Unmarshal(buf, x)
	}
	return x, err
}

func (p *{{ $table.TypeName }}) Get{{ .GoName }}() (*{{ .JSONType }}, error) {
	out, err := p.{{ .GoName }}.Get({{ .GetterFuncName }})
	if err != nil {
		log.Printf("failed unmarshal {{ .JSONType }}, {{ $table.PrimaryKey }}= %v, err= %v", p.{{ $table.PKFieldName }}, err)
		return nil, err
	}
	return out.(*{{ .JSONType }}), nil
}

func (p *{{ $table.TypeName }}) Set{{ .GoName }}({{ .VarName }} *{{ .JSONType }}) {
	buf, _ := json.Marshal({{ .VarName }})
	p.{{ .GoName }}.Set(buf, {{ .VarName }})
}
{{ end }}

{{ end }}`)
}
