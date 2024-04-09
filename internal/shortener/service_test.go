package shortener

import (
	"testing"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/repo/inmemory"
)

const BaseURL = "http://localhost:8080"

func BenchmarkService(b *testing.B) {
	//db, err := pkg.InitDatabasePg(PgDsn)
	//if err != nil {
	//	panic(err)
	//}

	repo, err := inmemory.New("/tmp/short-url-repo.json") //pg.New(db)
	if err != nil {
		panic(err)
	}

	shortener := New(BaseURL, repo)

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
