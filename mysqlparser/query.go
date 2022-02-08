package mysqlparser

import (
	"fmt"
	"regexp"
	"strings"
)

type Operator int

const (
	Equal Operator = iota
	NotEqual
	LessThan
	LessThanOrEqual
	GreaterThan
	GreaterThanOrEqual
	In
	NotIn
)

func (op Operator) format() string {
	switch op {
	case Equal:
		return "="
	case NotEqual:
		return "<>"
	case LessThan:
		return "<"
	case LessThanOrEqual:
		return "<="
	case GreaterThan:
		return ">"
	case GreaterThanOrEqual:
		return ">="
	case In:
		return "in"
	case NotIn:
		return "not in"
	}
	panic(fmt.Sprintf("unknown operator: %v", op))
}

var opTable = map[string]Operator{
	"eq":  Equal,
	"ne":  NotEqual,
	"lt":  LessThan,
	"lte": LessThanOrEqual,
	"gt":  GreaterThan,
	"gte": GreaterThanOrEqual,
	"in":  In,
	"nin": NotIn,
}

type QueryArg struct {
	Col *Column
	Op  Operator
}

func (a QueryArg) GoType() string {
	if a.Op == In || a.Op == NotIn {
		return "[]" + a.Col.GoType()
	}
	return a.Col.GoType()
}

func (a QueryArg) VarName() string {
	if a.Op == In || a.Op == NotIn {
		return a.Col.ListVarName()
	}
	return a.Col.VarName()
}

func (a QueryArg) Placeholder() string {
	if a.Op == In || a.Op == NotIn {
		return "(?)"
	}
	return "?"
}

type Query struct {
	Name string
	Args []QueryArg
}

func (q *Query) IsMGet() bool {
	return strings.HasPrefix(q.Name, "MGet")
}

func (q *Query) ArgList() string {
	var buf strings.Builder
	for i, arg := range q.Args {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%s %s, ", arg.VarName(), arg.GoType()))
	}
	return buf.String()
}

func (q *Query) Where() string {
	var buf strings.Builder
	buf.WriteString(`"`)
	for i, arg := range q.Args {
		if i > 0 {
			buf.WriteString(" and ")
		}
		buf.WriteString(fmt.Sprintf("%s %s %s", arg.Col.Name, arg.Op.format(), arg.Placeholder()))
	}
	buf.WriteString(`", `)
	for i, arg := range q.Args {
		if i > 0 {
			buf.WriteString(", ")
		}
		buf.WriteString(arg.VarName())
	}
	return buf.String()
}

var queryRE = regexp.MustCompile(`^(\w+)\(([^)]+)\)$`)

var argsRE = regexp.MustCompile(`^(\w+(?:\.\w+)?)(?:,\s*(\w+(?:\.\w+)?))*$`)

func ParseQuery(t *Table, q string) *Query {
	match := queryRE.FindStringSubmatch(q)
	if len(match) < 3 {
		panic(fmt.Sprintf("query definition is invalid: %q", q))
	}
	name := match[1]
	argsStr := match[2]
	argsMatch := argsRE.FindStringSubmatch(argsStr)
	if len(argsMatch) < 1 {
		panic(fmt.Sprintf("query definition is invalid: %q", q))
	}
	var args []QueryArg
	for i := 1; i < len(argsMatch); i++ {
		parts := strings.SplitN(argsMatch[i], ".", 2)
		colName := parts[0]
		col := t.GetColumn(colName)
		if col == nil {
			panic(fmt.Sprintf("query definition is invalid: %q", q))
		}
		opStr := "eq"
		if len(parts) > 1 {
			opStr = parts[1]
		}
		op, valid := opTable[opStr]
		if !valid {
			panic(fmt.Sprintf("query definition is invalid: %q", q))
		}
		args = append(args, QueryArg{
			Col: col,
			Op:  op,
		})
	}
	return &Query{
		Name: name,
		Args: args,
	}
}
