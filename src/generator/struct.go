package generator

import (
	"encoding/json"

	"github.com/samber/lo"
)

var _ Gen = &StructGen{}

type StructGen struct {
	Fields []*StructFieldGen
}

func (g *StructGen) AddChild(name string, child Gen) {
	g.Fields = append(g.Fields, &StructFieldGen{Name: name, Value: child})
}

func (g *StructGen) Gen() any {
	field2Data := lo.SliceToMap(g.Fields, func(field *StructFieldGen) (string, any) {
		return field.Name, field.Gen()
	})

	return json.RawMessage(MustJSONMarshal(field2Data))
}

type StructFieldGen struct {
	Name  string
	Value Gen
}

func (g *StructFieldGen) Gen() any {
	return g.Value.Gen()
}
