package main

import (
	"reflect"
	"testing"
)

func Test_parseRows(t *testing.T) {
	type args struct {
		csv string
	}
	tests := []struct {
		name string
		args args
		want [][]any
	}{
		{name: "all strings, no special chars",
			args: args{csv: `a,b,c,d,e,f`},
			want: [][]any{
				{"a", "b", "c", "d", "e", "f"},
			},
		},
		{name: "ints and strings, no special chars",
			args: args{csv: `a,b,c,1,2,3`},
			want: [][]any{
				{"a", "b", "c", 1, 2, 3},
			},
		},
		{name: "ints, floats, and strings, no special chars",
			args: args{csv: `a,b,c,1,2,3,4.5,6.7E-1`},
			want: [][]any{
				{"a", "b", "c", 1, 2, 3, 4.5, 6.7e-1},
			},
		},
		{name: "booleans, ints, strings, no special chars",
			args: args{csv: `a,b,c,1,2,3,false,true`},
			want: [][]any{
				{"a", "b", "c", 1, 2, 3, false, true},
			},
		},
		{name: "escaped newline",
			args: args{csv: `a,b,c\
1,2,3`},
			want: [][]any{
				{"a", "b", "c\n1", 2, 3},
			},
		},
		{name: "quoted newline",
			args: args{csv: `a,b,"c
1",2,3`},
			want: [][]any{
				{"a", "b", "\"c\n1\"", 2, 3},
			},
		},
		{name: "boolean values quoted, unquoted",
			args: args{csv: `"false",true,"b
1","true",false`},
			want: [][]any{
				{"\"false\"", true, "\"b\n1\"", "\"true\"", false},
			},
		},
		{name: "escaped escape char",
			args: args{csv: `"\\a\\b",1,3,4,abc`},
			want: [][]any{
				{"\"\\a\\b\"", 1, 3, 4, "abc"},
			},
		},
		{name: "two sequential double-quote chars",
			args: args{csv: `"" I like quotes "",1,3,4,abc`},
			want: [][]any{
				{"\" I like quotes \"", 1, 3, 4, "abc"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ParseRows(tt.args.csv); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRows() = %#v, want %#v", got, tt.want)
			}
		})
	}
}
