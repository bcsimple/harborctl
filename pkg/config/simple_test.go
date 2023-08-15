package config

import (
	"fmt"
	"testing"
)

func TestNewSimpleConnectInfo(t *testing.T) {
	tests := []struct {
		name string
		want *SimpleConnectInfo
	}{
		// TODO: Add test cases.
		{
			name: "simple1",
			want: &SimpleConnectInfo{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewSimpleConnectInfo()
			fmt.Println(got)
		})
	}
}
