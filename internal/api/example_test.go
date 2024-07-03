package api

import (
	"github.com/Azzonya/go-shortener/internal/repo/pg"
	"github.com/Azzonya/go-shortener/internal/shortener"
	"github.com/Azzonya/go-shortener/pkg"
)

func Example() {
	db, err := pkg.InitDatabasePg("postgresql://postgres:postgres@localhost:5437/postgresdb")
	if err != nil {
		panic(err)
	}

	repo := pg.New(db)

	shortenerService := shortener.New("http://localhost:8594", repo)

	api := New(shortenerService, "jwt_secret", "", false, nil)

	api.Start("localhost:8594", "localhost:8082")
}
