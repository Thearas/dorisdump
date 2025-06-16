package generator

import (
	"bytes"
	"math/rand/v2"

	"github.com/goccy/go-json"
)

type MapGen struct {
	Key, Value     Gen
	LenMin, LenMax int
}

func (g *MapGen) SetKeyGen(k Gen) {
	g.Key = k
}

func (g *MapGen) SetValueGen(v Gen) {
	g.Value = v
}

func (g *MapGen) Gen() any {
	len := rand.IntN(g.LenMax-g.LenMin+1) + g.LenMin

	var b bytes.Buffer
	b.WriteByte('{')
	for i := 0; i < len; i++ {
		key := g.Key.Gen()
		value := g.Value.Gen()
		b.Write(MustJsonMarshal(key))
		b.WriteByte(':')
		b.Write(MustJsonMarshal(value))
		if i < len-1 {
			b.WriteByte(',')
		}
	}
	b.WriteByte('}')
	return json.RawMessage(b.Bytes())
}
