package shortener

import (
	"reflect"
	"testing"

	"github.com/Azzonya/go-shortener/internal/repo"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
)

const BaseURL = "http://localhost:8080"
const fileStoragePath = "/tmp/short-url-repo.json"

func BenchmarkService(b *testing.B) {
	//db, err := pkg.InitDatabasePg(PgDsn)
	//if err != nil {
	//	panic(err)
	//}

	repoTest, err := inmemory.New("/tmp/short-url-repo.json") //pg.New(db)
	if err != nil {
		panic(err)
	}

	shortener := New(BaseURL, repoTest)

	urls := []*entities.ReqURL{
		{
			ID:          "1",
			OriginalURL: "blab2la.com",
			ShortURL:    "",
		},
		{
			ID:          "1",
			OriginalURL: "ssdsd3s.kz",
			ShortURL:    "",
		},
	}

	b.ResetTimer()

	b.Run("generate_shortURL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			shortener.GenerateShortURL()
		}
	})

	b.Run("shorten_urls", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			shortener.ShortenURLs(urls, "1")
		}
	})

	b.Run("shorten_and_save_link", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			shortener.ShortenAndSaveLink("ysadsadas", "1")
		}
	})

	b.Run("get_one_by_originalURL", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			shortener.GetOneByOriginalURL("blab2la.com")
		}
	})
}

func TestNew(t *testing.T) {
	repoTest, err := inmemory.New(fileStoragePath)
	if err != nil {
		panic(err)
	}

	type args struct {
		baseURL string
		repo    *inmemory.St
	}
	tests := []struct {
		name string
		args args
		want *Shortener
	}{
		{
			name: "New shortener test",
			args: args{
				baseURL: BaseURL,
				repo:    repoTest,
			},
			want: &Shortener{
				repo:    repoTest,
				baseURL: BaseURL,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.baseURL, tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_DeleteURLs(t *testing.T) {
	repoTest, err := inmemory.New(fileStoragePath)
	if err != nil {
		panic(err)
	}

	type fields struct {
		repo    *inmemory.St
		baseURL string
	}
	type args struct {
		urls   []string
		userID string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Generate short test",
			fields: fields{
				baseURL: BaseURL,
				repo:    repoTest,
			},
			args: args{
				urls:   []string{},
				userID: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			s.DeleteURLs(tt.args.urls, tt.args.userID)
		})
	}
}

func TestShortener_GenerateShortURL(t *testing.T) {
	repoTest, err := inmemory.New(fileStoragePath)
	if err != nil {
		panic(err)
	}

	type fields struct {
		repo    *inmemory.St
		baseURL string
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "Generate short test",
			fields: fields{
				baseURL: BaseURL,
				repo:    repoTest,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			if got := s.GenerateShortURL(); got == "" {
				t.Errorf("GenerateShortURL() = %v, want not empty", got)
			}
		})
	}
}

func TestShortener_GetOneByOriginalURL(t *testing.T) {
	type fields struct {
		repo    repo.Repo
		baseURL string
	}
	type args struct {
		url string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			got, got1 := s.GetOneByOriginalURL(tt.args.url)
			if got != tt.want {
				t.Errorf("GetOneByOriginalURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetOneByOriginalURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestShortener_GetOneByShortURL(t *testing.T) {
	type fields struct {
		repo    repo.Repo
		baseURL string
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
		want1  bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			got, got1 := s.GetOneByShortURL(tt.args.key)
			if got != tt.want {
				t.Errorf("GetOneByShortURL() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetOneByShortURL() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestShortener_IsDeleted(t *testing.T) {
	type fields struct {
		repo    repo.Repo
		baseURL string
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			if got := s.IsDeleted(tt.args.shortURL); got != tt.want {
				t.Errorf("IsDeleted() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_ListAll(t *testing.T) {
	type fields struct {
		repo    repo.Repo
		baseURL string
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			got, err := s.ListAll(tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_PingDB(t *testing.T) {
	type fields struct {
		repo    repo.Repo
		baseURL string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			if err := s.PingDB(); (err != nil) != tt.wantErr {
				t.Errorf("PingDB() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestShortener_ShortenAndSaveLink(t *testing.T) {
	type fields struct {
		repo    repo.Repo
		baseURL string
	}
	type args struct {
		originalURL string
		userID      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			got, err := s.ShortenAndSaveLink(tt.args.originalURL, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortenAndSaveLink() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ShortenAndSaveLink() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestShortener_ShortenURLs(t *testing.T) {
	repoTest, err := inmemory.New(fileStoragePath)
	if err != nil {
		panic(err)
	}

	type fields struct {
		repo    *inmemory.St
		baseURL string
	}
	type args struct {
		urls   []*entities.ReqURL
		userID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*entities.ReqURL
		wantErr bool
	}{
		{
			name: "New shortener test",
			fields: fields{
				baseURL: BaseURL,
				repo:    repoTest,
			},
			args: args{
				urls: []*entities.ReqURL{
					{
						OriginalURL: "test.com",
					},
				},
			},
			want: []*entities.ReqURL{
				{
					OriginalURL: "test.com",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Shortener{
				repo:    tt.fields.repo,
				baseURL: tt.fields.baseURL,
			}
			got, err := s.ShortenURLs(tt.args.urls, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ShortenURLs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("ShortenURLs() got = %v, want %v", got, tt.want)
			}
		})
	}
}
