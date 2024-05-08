package inmemory

import (
	"github.com/Azzonya/go-shortener/internal/entities"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		args    args
		want    *St
		wantErr bool
	}{
		{
			name: "New inmemory test",
			args: args{
				filePath: "/tmp/short-url-repo.json",
			},
			want: &St{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.filePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.filePath, tt.want.filePath) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSt_Add(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	type args struct {
		originalURL string
		shortURL    string
		userID      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Add inmemory test",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			args: args{
				shortURL:    "tst",
				originalURL: "test.com",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				URLMap:   tt.fields.URLMap,
				filePath: tt.fields.filePath,
				lastID:   tt.fields.lastID,
			}
			if err := s.Add(tt.args.originalURL, tt.args.shortURL, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_CreateShortURLs(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	type args struct {
		urls   []*entities.ReqURL
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create short urls inmemory test",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				URLMap:   tt.fields.URLMap,
				filePath: tt.fields.filePath,
				lastID:   tt.fields.lastID,
			}
			if err := s.CreateShortURLs(tt.args.urls, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("CreateShortURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_DeleteURLs(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	type args struct {
		urls   []string
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Create short urls inmemory test",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			args:    args{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				URLMap:   tt.fields.URLMap,
				filePath: tt.fields.filePath,
				lastID:   tt.fields.lastID,
			}
			if err := s.DeleteURLs(tt.args.urls, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_GetByOriginalURL(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	type args struct {
		originalURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "Get by original url inmemory test",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			args: args{
				"test.com",
			},
			want:  "tst",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				URLMap:   tt.fields.URLMap,
				filePath: tt.fields.filePath,
				lastID:   tt.fields.lastID,
			}
			if err := s.Add(tt.args.originalURL, "tst", ""); (err != nil) != false {
				t.Errorf("Add() error = %v, wantErr %v", err, false)
			}

			got, got1 := s.GetByOriginalURL(tt.args.originalURL)
			if got != tt.want {
				t.Errorf("GetByOriginalURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetByOriginalURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSt_GetByShortURL(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		{
			name: "Get by short url inmemory test",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			args: args{
				"tst",
			},
			want:  "test.com",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				URLMap:   tt.fields.URLMap,
				filePath: tt.fields.filePath,
				lastID:   tt.fields.lastID,
			}
			if err := s.Add(tt.want, tt.args.shortURL, ""); (err != nil) != false {
				t.Errorf("Add() error = %v, wantErr %v", err, false)
			}

			got, got1 := s.GetByShortURL(tt.args.shortURL)
			if got != tt.want {
				t.Errorf("GetByShortURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetByShortURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestSt_Initialize(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Initialize inmemory test",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				filePath: tt.fields.filePath,
			}
			if err := s.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_ListAll(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	type args struct {
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.ReqListAll
		wantErr bool
	}{
		{
			name: "ListAll inmemory test",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
			args: args{
				"",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				URLMap:   tt.fields.URLMap,
				filePath: tt.fields.filePath,
				lastID:   tt.fields.lastID,
			}
			_, err := s.ListAll(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSt_SyncData(t *testing.T) {
	type fields struct {
		URLMap   map[string]string
		filePath string
		lastID   int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "SyncData inmemory",
			fields: fields{
				filePath: "/tmp/short-url-repo.json",
				URLMap:   make(map[string]string),
				lastID:   0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				URLMap:   tt.fields.URLMap,
				filePath: tt.fields.filePath,
				lastID:   tt.fields.lastID,
			}
			s.SyncData()
		})
	}
}
