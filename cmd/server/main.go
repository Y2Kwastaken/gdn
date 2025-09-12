package main

import (
	"log"

	"github.com/Y2Kwastaken/gdn/internal"
	_ "modernc.org/sqlite"
)

func main() {
	err := internal.LoadEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := internal.NewDBConnection("file:gdn_main.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	err = db.SetupTables()
	if err != nil {
		log.Fatal(err)
	}

	store := internal.NewObjectStore("localhost:9000")
	store.Database = db
	err = store.Connect("admin", "password")
	if err != nil {
		log.Fatal(err)
	}

	internal.SetupHttpServer(store)
}
