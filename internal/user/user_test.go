package user

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		want    *User
		wantErr bool
	}{
		{
			name: "new user",
			want: &User{
				ID:  "smth",
				new: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New()
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.new != tt.want.new {
				t.Errorf("New() got = %v, want %v", got.new, tt.want.new)
			}
		})
	}
}

func TestNewWithID(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name string
		args args
		want *User
	}{
		{
			name: "with id",
			args: args{
				id: "5",
			},
			want: &User{
				ID:  "5",
				new: false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWithID(tt.args.id); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWithID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUser_IsNew(t *testing.T) {
	type fields struct {
		ID  string
		new bool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "check is new - false",
			fields: fields{
				ID:  "1",
				new: false,
			},
			want: false,
		},
		{
			name: "check is new - true",
			fields: fields{
				ID:  "2",
				new: true,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &User{
				ID:  tt.fields.ID,
				new: tt.fields.new,
			}
			if got := u.IsNew(); got != tt.want {
				t.Errorf("IsNew() = %v, want %v", got, tt.want)
			}
		})
	}
}
