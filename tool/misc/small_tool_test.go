package misc

import "testing"

func TestGetIpAddr(t *testing.T) {
	type args struct {
		Ip string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestGetIpAddr",
			args: args{
				Ip: "122.224.197.37",
			},
			want: "中国",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetIpAddr(tt.args.Ip); got != tt.want {
				t.Errorf("GetIpAddr() = %v, want %v", got, tt.want)
			}
		})
	}
}
