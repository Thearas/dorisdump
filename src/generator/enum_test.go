package generator

import (
	"reflect"
	"strings"
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"

	"github.com/Thearas/dodo/src/parser"
)

func TestEnumGen_Gen(t *testing.T) {
	type fields struct {
		Enum    []any
		Weights []float32
	}
	tests := []struct {
		name   string
		fields fields
		want   any
	}{
		{
			name: "simple",
			fields: fields{
				Enum:    []any{1, 2, 3.0},
				Weights: []float32{0, 0, 1},
			},
			want: 3.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &EnumGen{
				Enum:    tt.fields.Enum,
				Weights: tt.fields.Weights,
			}
			if got := g.Gen(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EnumGen.Gen() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewEnumGenRule(t *testing.T) {
	type args struct {
		dataType string
		r        GenRule
	}
	tests := []struct {
		name    string
		args    args
		want    Gen
		wantErr bool
	}{
		{
			name: "simple",
			args: args{
				dataType: "int",
				r:        MustYAMLUmarshal("{enum: [1, 2, 3], weights: [0.4, 0.5, 0.1]}"),
			},
			want: &EnumGen{
				Enum:    []any{1, 2, 3},
				Weights: []float32{0.4, 0.5, 0.1},
			},
			wantErr: false,
		},
		{
			name: "complex",
			args: args{
				dataType: "varchar(100)",
				r: MustYAMLUmarshal(`
enum:
    - length: 5
    - length: {min: 5, max: 10}
    - format: "int to str: {{%d}}"
      gen:
          enum: [1, 2, 3]
    - format: "{{%d}}"
      gen:
          ref: t1.col1
weights: [0.4, 0.4, 0.1, 0.1]
 `),
			},
			wantErr: false,
		},
		{
			name: "err: less weights",
			args: args{
				dataType: "int",
				r:        MustYAMLUmarshal("{enum: [1, 2, 3], weights: [0.4, 0.5]}"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewParser(tt.name, tt.args.dataType)
			dataType := p.DataType()
			assert.NoError(t, p.ErrListener.LastErr)
			got, err := NewEnumGenerator(NewTypeVisitor(tt.name, nil), dataType, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEnumGenRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if strings.HasPrefix(tt.name, "complex") {
				// inject values to ref t1.col1
				refgen := getColumnRefGen("t1", "col1")
				refgen.AddRefVals(lo.ToAnySlice(lo.Range(996))...)
				enum := got.(*EnumGen).Enum
				for _, v := range enum {
					for range 100 {
						assert.IsType(t, "", v.(Gen).Gen())
					}
				}
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEnumGenRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
