package generator

import (
	"testing"

	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
)

func TestRefGenerator(t *testing.T) {
	type args struct {
		in0 string
		r   GenRule
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "ref_col1 1",
			args:    args{r: GenRule{"ref": "table1.col1", "limit": 10}},
			wantErr: false,
		},
		{
			name:    "ref_other",
			args:    args{r: GenRule{"ref": "table1.other_col", "limit": 10}},
			wantErr: false,
		},
		{
			name:    "ref_col1 2",
			args:    args{r: GenRule{"ref": "table1.col1", "limit": 20}},
			wantErr: false,
		},
		{
			name:    "ref_col1 3",
			args:    args{r: GenRule{"ref": "table1.col1", "limit": 100}},
			wantErr: false,
		},
	}
	refGens := []*RefGen{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewRefGenerator(NewTypeVisitor(tt.name, nil), nil, tt.args.r)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewRefGenerator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			refGens = append(refGens, got.(*RefGen))
		})
	}

	shardColGens := GetTableRefGen("table1")
	assert.Len(t, shardColGens, 2)

	refCol1 := getColumnRefGen("table1", "col1")
	assert.Equal(t, 100, refCol1.Limit)
	refCol1_2 := refGens[2]
	assert.Equal(t, 20, refCol1_2.Limit)
	refColOther := getColumnRefGen("table1", "other_col")
	assert.Equal(t, 10, refColOther.Limit)

	refCol1.AddRefVals()
	refCol1.AddRefVals(1, 2, 3, 4)
	refCol1.AddRefVals(lo.ToAnySlice(lo.Range(996))...)
	refCol1.AddRefVals(lo.ToAnySlice(lo.RangeFrom(1000, 100))...)

	refCol1_3 := refGens[3]
	assert.Len(t, *refCol1_3.refValsPtr, len(*refCol1.refValsPtr))
	assert.Len(t, *refCol1_3.refValsPtr, 100)
	for range 100 {
		assert.Less(t, refCol1_3.Gen(), 1100)
	}
}
