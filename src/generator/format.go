package generator

import (
	"fmt"
	"io"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/valyala/fasttemplate"
)

var _ Gen = &FormatGen{}

type FormatGen struct {
	Format string
	inner  Gen

	template *fasttemplate.Template
}

func (g *FormatGen) Gen() any {
	var result any = nil
	if g.inner != nil {
		result = g.inner.Gen()
	}

	formatted, err := g.template.ExecuteFuncStringWithErr(func(w io.Writer, tag string) (int, error) {
		if strings.HasPrefix(tag, "%") {
			// inject result
			return w.Write(fmt.Appendf(nil, tag, result))
		}
		tagF, ok := FormatTags[tag]
		if !ok {
			return 0, fmt.Errorf("unknown format tag '%s'", tag)
		}

		result := tagF.Call(nil)[0].Interface()
		if result == nil {
			return w.Write([]byte(`\N`))
		} else if s, ok := result.(string); ok {
			return w.Write([]byte(s))
		}
		return w.Write(fmt.Append(nil, result))
	})
	if err != nil {
		logrus.Errorf("format execute templace failed, err: %v\n", err)
	}

	return formatted
}

func NewFormatGenerator(format string, inner Gen) (Gen, error) {
	t, err := fasttemplate.NewTemplate(format, "{{", "}}")
	if err != nil {
		return nil, err
	}

	return &FormatGen{
		Format:   format,
		inner:    inner,
		template: t,
	}, nil
}
