package config

import "testing"

func TestConnectConfiguration_SetConnectInfo(t *testing.T) {
	type fields struct {
		APIVersion     string
		Kind           string
		Harbors        HarborConnectInfos
		CurrentContext string
	}
	type args struct {
		info *HarborConnectInfo
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name:   "test1",
			fields: fields{},
			args: args{
				info: &HarborConnectInfo{
					Name:     "harbor1",
					Host:     "172.16.69.50:1121",
					User:     "admin",
					Password: "Harbor12345",
				},
			},
		},

		{
			name:   "test2",
			fields: fields{},
			args: args{
				info: &HarborConnectInfo{
					Name:     "harbor2",
					Host:     "10.0.0.1:1121",
					User:     "admin",
					Password: "admin1",
				},
			},
		},

		{
			name:   "test3",
			fields: fields{},
			args: args{
				info: &HarborConnectInfo{
					Name:     "harbor3",
					Host:     "10.0.0.1:1121",
					User:     "admin",
					Password: "admin",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//初始化
			c := NewConnectConfiguration()
			c.SetConnectInfo(tt.args.info)
		})
	}
}

func TestConnectConfiguration_SetConnectInfoAlias(t *testing.T) {
	type fields struct {
		APIVersion     string
		Kind           string
		Harbors        HarborConnectInfos
		CurrentContext string
	}
	type args struct {
		name      string
		aliasName string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				name:      "harbor1",
				aliasName: "1",
			},
			wantErr: false,
		},
		{
			name: "test1",
			args: args{
				name:      "harbor1",
				aliasName: "1",
			},
			wantErr: false,
		},

		{
			name: "test1",
			args: args{
				name:      "harbor2",
				aliasName: "2",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConnectConfiguration()
			if err := c.SetConnectInfoAlias(tt.args.name, tt.args.aliasName); (err != nil) != tt.wantErr {
				t.Errorf("SetConnectInfoAlias() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestConnectConfiguration_SetConnectInfoContext(t *testing.T) {
	type fields struct {
		APIVersion     string
		Kind           string
		Harbors        HarborConnectInfos
		CurrentContext string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "context1",
			args: args{
				name: "harbor1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConnectConfiguration()
			c.SetConnectInfoContext(tt.args.name)
		})
	}
}

func TestConnectConfiguration_DelConnectInfo(t *testing.T) {
	type fields struct {
		APIVersion     string
		Kind           string
		Harbors        HarborConnectInfos
		CurrentContext string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "del1",
			args: args{
				name: "harbor2",
			},
		},
		{
			name: "del1",
			args: args{
				name: "harbor1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewConnectConfiguration()
			c.DelConnectInfo(tt.args.name)
		})
	}
}
