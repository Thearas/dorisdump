package generator

import (
	"errors"
	"fmt"
	"math/rand/v2"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/samber/lo"
	"github.com/spf13/cast"

	"github.com/Thearas/dodo/src/parser"
)

var _ Gen = &EnumGen{}

type EnumGen struct {
	Enum    []any     `yaml:"enum,omitempty"`
	Weights []float32 `yaml:"weights,omitempty"`
}

func (g *EnumGen) Gen() any {
	v := g.gen()
	if gr, ok := v.(Gen); ok {
		return gr.Gen()
	}
	return v
}

//nolint:revive
func (g *EnumGen) gen() any {
	if len(g.Weights) == 0 {
		return g.Enum[gofakeit.IntN(len(g.Enum))]
	}

	weight := rand.Float32()
	for i, w := range g.Weights {
		weight -= w
		if weight < 0 {
			return g.Enum[i]
		}
	}

	panic("EnumGen.Gen(): unreachable")
}

func NewEnumGenerator(visitor *TypeVisitor, dataType parser.IDataTypeContext, r GenRule) (Gen, error) {
	enum_ := r["enum"]
	if enum_ == nil {
		enum_ = cast.ToStringSlice(r["enums"])
	}
	enum, ok := enum_.([]any)
	if !ok || len(enum) == 0 {
		return nil, errors.New("enum is empty")
	}
	for i, v := range enum {
		gr, ok := v.(GenRule)
		if !ok {
			continue
		}
		enum[i] = visitor.GetChildGen(fmt.Sprintf("enum.%d", i), dataType, gr)
	}

	weights_ := r["weights"]
	if weights_ == nil {
		weights_ = r["weight"]
	}
	if weights_ == nil {
		return &EnumGen{Enum: enum}, nil
	}
	ws, ok := weights_.([]any)
	if !ok || len(ws) == 0 {
		return nil, fmt.Errorf("weights should be a [float], but got: %T", r["weights"])
	}
	weights := lo.Map(ws, func(w any, _ int) float32 { return cast.ToFloat32(w) })
	if len(weights) != len(enum) {
		return nil, errors.New("enum length not equals to weights length")
	}
	if lo.Sum(weights) != 1 {
		return nil, errors.New("sum of weights should be 1")
	}

	return &EnumGen{
		Enum:    enum,
		Weights: weights,
	}, nil
}
