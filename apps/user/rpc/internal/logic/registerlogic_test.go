package logic

import (
	"GoChat/apps/user/rpc/user"
	"context"
	"testing"
)

func TestRegisterLogic_Register(t *testing.T) {
	type args struct {
		in *user.RegisterReq
	}
	tests := []struct {
		name      string
		args      args
		wantPrint bool
		wantErr   bool
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{in: &user.RegisterReq{
				Phone:    "12345678900",
				Nickname: "user1",
				Password: "123456",
				Avatar:   "png.jpg",
				Sex:      1,
			}},
			wantPrint: true,
			wantErr:   false,
		},
		{
			name: "test2",
			args: args{in: &user.RegisterReq{
				Phone:    "12345678900",
				Nickname: "user2",
				Password: "123456",
				Avatar:   "png.jpg",
				Sex:      1,
			}},
			wantPrint: true,
			wantErr:   false,
		},
		{
			name: "test3",
			args: args{in: &user.RegisterReq{
				Phone:    "12345678901",
				Nickname: "user3",
				Password: "",
				Avatar:   "png.jpg",
				Sex:      1,
			}},
			wantPrint: true,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewRegisterLogic(context.Background(), svcCtx)
			got, err := l.Register(tt.args.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantPrint {
				t.Log(tt.name, got)
			}
		})
	}
}
