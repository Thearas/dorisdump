package generator

import (
	"math/rand/v2"

	"github.com/goccy/go-json"
	"github.com/samber/lo"
)

type ArrayGen struct {
	Element        Gen
	LenMin, LenMax int
}

func (g *ArrayGen) SetElementGen(elem Gen) {
	g.Element = elem
}

func (g *ArrayGen) Gen() any {
	length := rand.IntN(g.LenMax-g.LenMin+1) + g.LenMin

	elementData := lo.RepeatBy(length, func(_ int) any {
		return g.Element.Gen()
	})

	return json.RawMessage(MustJSONMarshal(elementData))
}
