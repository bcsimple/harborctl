package encrypt

import (
	"fmt"
	"testing"
)

func TestDecrypt(t *testing.T) {
	type args struct {
		encryptedDataStr string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				encryptedDataStr: "zVfmsys=",
			},
		},
		{
			name: "test2",
			args: args{
				encryptedDataStr: "zVfmsysUPQ==",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.encryptedDataStr)
			fmt.Println(got, err)
		})
	}
}

func TestEncrypt(t *testing.T) {
	type args struct {
		rawData string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				rawData: "admin",
			},
		},
		{
			name: "test1",
			args: args{
				rawData: "admin-1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Encrypt(tt.args.rawData)
			fmt.Println(got, err)
		})
	}
}
