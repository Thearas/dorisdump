package generator

import (
	"fmt"

	"github.com/spf13/cast"

	"github.com/Thearas/dodo/src/parser"
)

var _ Gen = &TypeGen{}

type TypeGen struct {
	Type    string
	GenRule GenRule

	gen Gen
}

func (g *TypeGen) Gen() any {
	return g.gen.Gen()
}

func NewTypeGenerator(colpath string, r GenRule) (Gen, error) {
	ty := cast.ToString(r["type"])
	g := &TypeGen{
		Type:    ty,
		GenRule: r,
	}

	p := parser.NewParser(colpath, ty)
	dataType := p.DataType()
	if p.ErrListener.LastErr != nil {
		return nil, fmt.Errorf("parse type generator failed for column '%s', err: %v", colpath, p.ErrListener.LastErr)
	}

	visitor := NewTypeVisitor(colpath, r)
	g.gen = visitor.GetTypeGen(dataType)

	return g, nil
}
