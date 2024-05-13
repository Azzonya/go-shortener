package main

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/app"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	fmt.Println("Build version:", getValueOrDefault(buildVersion))
	fmt.Println("Build date:", getValueOrDefault(buildDate))
	fmt.Println("Build commit:", getValueOrDefault(buildCommit))

	app.Start()
}

func getValueOrDefault(value string) string {
	if value == "" {
		return "N/A"
	}
	return value
}
