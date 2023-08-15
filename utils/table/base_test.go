package table

import (
	"testing"
)

func TestNewTableInformation(t *testing.T) {
	type args struct {
		style string
	}
	tests := []struct {
		name  string
		args  args
		want  TableInformation
		title []string
		data  [][]string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				style: "kube",
			},
			title: []string{"name", "age"},
			data: [][]string{
				{
					"zhangsan",
					"20",
				},
				{
					"wangwu",
					"21",
				},
			},
		},
		{
			name: "test2",
			args: args{
				style: "table",
			},
			title: []string{"name", "age"},
			data: [][]string{
				{
					"zhangsan",
					"20",
				},
				{
					"wangwu",
					"21",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			NewTableInformation(tt.args.style).SetTitles(tt.title).SetData(tt.data).Output()
		})
	}
}
