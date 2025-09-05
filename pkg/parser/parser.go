package parser

import (
	"fmt"
	"strings"

	sql "github.com/xwb1989/sqlparser"
)

type Query struct {
	Fields   []string
	Resource string
	Where    string
	OrderBy  []Order
	Limit    *int
	Offset   *int
}

type Order struct {
	Field string
	Desc  bool
}

func Parse(q string) (*Query, error) {
	stmt, err := sql.Parse(q)
	if err != nil {
		return nil, err
	}
	sel, ok := stmt.(*sql.Select)
	if !ok {
		return nil, fmt.Errorf("only SELECT statements are supported")
	}

	res := &Query{}

	if sel.From == nil || len(sel.From) != 1 {
		return nil, fmt.Errorf("exactly one FROM resource is required")
	}

	tblExpr := sel.From[0]
	tbl, ok := tblExpr.(*sql.AliasedTableExpr)
	if !ok {
		return nil, fmt.Errorf("unsupported FROM expression")
	}
	name := sql.String(tbl.Expr)
	name = strings.Trim(name, "`")
	res.Resource = strings.ToLower(name)

	if len(sel.SelectExprs) == 1 {
		if _, ok := sel.SelectExprs[0].(*sql.StarExpr); ok {
			res.Fields = []string{"*"}
		}
	}
	if len(res.Fields) == 0 {
		for _, e := range sel.SelectExprs {
			ae, ok := e.(*sql.AliasedExpr)
			if !ok {
				return nil, fmt.Errorf("unsupported select expression")
			}
			col := sql.String(ae.Expr)
			col = strings.TrimSpace(strings.Trim(col, "`"))
			res.Fields = append(res.Fields, col)
		}
	}

	if sel.Where != nil {
		res.Where = sql.String(sel.Where.Expr)
	}

	for _, ob := range sel.OrderBy {
		field := sql.String(ob.Expr)
		res.OrderBy = append(res.OrderBy, Order{Field: field, Desc: ob.Direction == sql.DescScr})
	}

	return res, nil
}
