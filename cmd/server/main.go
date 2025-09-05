package main

import (
	"github.com/Y2Kwastaken/gdn/internal"
	_ "modernc.org/sqlite"
)

func main() {
	store := internal.NewObjectStore("http://localhost:9000")
	store.Connect("admin", "password")

	internal.SetupHttpServer()
}
