package session

import (
	"context"
	"reflect"
	"testing"

	"github.com/Azzonya/go-shortener/internal/user"
)

func TestGetUser(t *testing.T) {
	testUser := &user.User{ID: "testID"}
	parentContext := context.Background()
	ctx := context.WithValue(parentContext, ctxKeyUID, testUser)

	type args struct {
		c context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *user.User
		wantErr bool
	}{
		{
			name: "Get user from context",
			args: args{
				ctx,
			},
			want:    testUser,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetUser(tt.args.c)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetUser() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetUserFromContext(t *testing.T) {
	testUser := &user.User{ID: "testID"}
	ctx := context.WithValue(context.Background(), ctxKeyUID, testUser)

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name   string
		args   args
		wantU  *user.User
		wantOk bool
	}{
		{
			name: "get from context",
			args: args{
				ctx,
			},
			wantU:  testUser,
			wantOk: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotU, gotOk := GetUserFromContext(tt.args.ctx)
			if !reflect.DeepEqual(gotU, tt.wantU) {
				t.Errorf("GetUserFromContext() gotU = %v, want %v", gotU, tt.wantU)
			}
			if gotOk != tt.wantOk {
				t.Errorf("GetUserFromContext() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}

func TestSetUserContext(t *testing.T) {
	testUser := &user.User{ID: "testID"}
	parentContext := context.Background()

	type args struct {
		parent context.Context
		u      *user.User
	}
	tests := []struct {
		name string
		args args
		want context.Context
	}{
		{
			name: "Set user context",
			args: args{
				parent: parentContext,
				u:      testUser,
			},
			want: context.WithValue(parentContext, ctxKeyUID, testUser),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetUserContext(tt.args.parent, tt.args.u); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetUserContext() = %v, want %v", got, tt.want)
			}
		})
	}
}
