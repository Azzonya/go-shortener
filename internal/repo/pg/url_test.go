package pg

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/pkg"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

var PgDsn = "postgresql://postgres:postgres@localhost:5437/postgresdb"

func TestNew(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}
	type args struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name string
		args args
		want *St
	}{
		{
			name: "new db test",
			args: args{
				db: db,
			},
			want: &St{
				db: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := New(tt.args.db)
			assert.NotNil(t, got)
		})
	}
}

func TestSt_Add(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
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
			name: "add pg test",
			fields: fields{
				db: db,
			},
			args: args{
				originalURL: "example.com",
				shortURL:    "tst",
				userID:      "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}
			source := rand.NewSource(time.Now().UnixNano())
			randomGenerator := rand.New(source)
			randomNumber := randomGenerator.Intn(10000)

			if err := s.Add(tt.args.originalURL+fmt.Sprint(randomNumber), tt.args.shortURL+fmt.Sprint(randomNumber), tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("Add() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_CreateShortURLs(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
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
			name: "pg create short URLs test",
			fields: fields{
				db: db,
			},
			args: args{
				urls: []*entities.ReqURL{
					{
						OriginalURL: "example.com",
						ShortURL:    "tst",
						ID:          "1",
					},
					{
						OriginalURL: "youtube.com",
						ShortURL:    "tst1",
						ID:          "1",
					},
				},
				userID: "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}

			for _, v := range tt.args.urls {
				source := rand.NewSource(time.Now().UnixNano())
				randomGenerator := rand.New(source)
				randomNumber := randomGenerator.Intn(10000)
				v.ShortURL += fmt.Sprint(randomNumber)
				v.OriginalURL += fmt.Sprint(randomNumber)
				v.ID += fmt.Sprint(randomNumber)
			}

			if err := s.CreateShortURLs(tt.args.urls, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("CreateShortURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_DeleteURLs(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
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
			name: "delete URLs test",
			fields: fields{
				db: db,
			},
			args: args{
				urls: []string{
					"testdelete",
				},
				userID: "999",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}

			source := rand.NewSource(time.Now().UnixNano())
			randomGenerator := rand.New(source)
			randomNumber := randomGenerator.Intn(10000)

			urls := []*entities.ReqURL{
				{
					OriginalURL: "testdelete.com" + fmt.Sprint(randomNumber),
					ShortURL:    "testdelete" + fmt.Sprint(randomNumber),
				},
			}

			if err := s.CreateShortURLs(urls, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("CreateShortURLs() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := s.DeleteURLs(tt.args.urls, tt.args.userID); (err != nil) != tt.wantErr {
				t.Errorf("DeleteURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_GetByOriginalURL(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
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
			name: "get by original URL pg test",
			fields: fields{
				db: db,
			},
			args: args{
				originalURL: "getbyoriginal.url",
			},
			want:  "getbyoriginal",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}

			s.Add(tt.args.originalURL, tt.want, "1")

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
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
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
			name: "get by short URL pg test",
			fields: fields{
				db: db,
			},
			args: args{
				"getbyshorturl",
			},
			want:  "getbyshorturl.com",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}

			s.Add(tt.want, tt.args.shortURL, "1")

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
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "initialize pg test",
			fields: fields{
				db: db,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}
			if err := s.Initialize(); (err != nil) != tt.wantErr {
				t.Errorf("Initialize() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_ListAll(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
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
			name: "list all pg test",
			fields: fields{
				db: db,
			},
			args: args{
				"1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}
			_, err := s.ListAll(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestSt_Ping(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "ping pg test",
			fields: fields{
				db: db,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}
			if err := s.Ping(); (err != nil) != tt.wantErr {
				t.Errorf("Ping() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSt_SyncData(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "sync data pg test",
			fields: fields{
				db: db,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}
			s.SyncData()
		})
	}
}

func TestSt_TableExist(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "table exist pg test",
			fields: fields{
				db: db,
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}
			if got := s.TableExist(); got != tt.want {
				t.Errorf("TableExist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSt_URLDeleted(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		shortURL string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		{
			name: "URL deleted pg test",
			fields: fields{
				db: db,
			},
			args: args{
				"urldeleted",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}

			s.Add("urldeleted.com", tt.args.shortURL, "1")

			if got := s.URLDeleted(tt.args.shortURL); got != tt.want {
				t.Errorf("URLDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSt_Update(t *testing.T) {
	db, err := pkg.InitDatabasePg(PgDsn)
	if err != nil {
		panic(err)
	}

	type fields struct {
		db *pgxpool.Pool
	}
	type args struct {
		originalURL string
		shortURL    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "update pg test",
			fields: fields{
				db: db,
			},
			args: args{
				originalURL: "updatepg.com",
				shortURL:    "update",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &St{
				db: tt.fields.db,
			}

			source := rand.NewSource(time.Now().UnixNano())
			randomGenerator := rand.New(source)
			randomNumber := randomGenerator.Intn(10000)

			s.Add(tt.args.originalURL, tt.args.shortURL, "1")

			if err := s.Update(tt.args.originalURL, tt.args.shortURL+fmt.Sprint(randomNumber)); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
