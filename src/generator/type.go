package generator

import (
	"fmt"

	"github.com/spf13/cast"

	"github.com/Thearas/dodo/src/parser"
)

var _ Gen = &TypeGen{}

type TypeGen struct {
	GenRule GenRule

	gen Gen
}

func (g *TypeGen) Gen() any {
	return g.gen.Gen()
}

func NewTypeGenerator(colpath string, _ parser.IDataTypeContext, r GenRule) (Gen, error) {
	p := parser.NewParser(colpath, cast.ToString(r["type"]))
	dataType := p.DataType()
	if p.ErrListener.LastErr != nil {
		return nil, fmt.Errorf("parse type generator failed for column '%s', err: %v", colpath, p.ErrListener.LastErr)
	}

	return newTypeGenerator(colpath, r, dataType), nil
}

func newTypeGenerator(colpath string, r GenRule, dataType parser.IDataTypeContext) Gen {
	g := &TypeGen{
		GenRule: r,
	}

	visitor := NewTypeVisitor(colpath, r)
	g.gen = visitor.GetTypeGen(dataType)

	return g
}
