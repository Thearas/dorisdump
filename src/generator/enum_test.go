package generator

import (
	"reflect"
	"testing"
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
		r GenRule
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
				r: MustYamlUmarshal("{enum: [1, 2, 3], weights: [0.4, 0.5, 0.1]}"),
			},
			want: &EnumGen{
				Enum:    []any{1, 2, 3},
				Weights: []float32{0.4, 0.5, 0.1},
			},
			wantErr: false,
		},
		{
			name: "err: less weights",
			args: args{
				r: MustYamlUmarshal("{enum: [1, 2, 3], weights: [0.4, 0.5]}"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewEnumGenerator("", tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewEnumGenRule() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEnumGenRule() = %v, want %v", got, tt.want)
			}
		})
	}
}
