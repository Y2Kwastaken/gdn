package internal

import (
	"net/http"
	"os"
)

var endpoint_handlers = make(map[string]func(*FileStore, string, http.ResponseWriter, *http.Request) int)

func registerEndpoints() {
	endpoint_handlers["photo"] = putPhoto
}

func putPhoto(store *FileStore, urlPart string, rspn http.ResponseWriter, rqst *http.Request) int {
	if rqst.Header.Get("X-API-Key") != os.Getenv("ADMIN_SECRET") {
		return http.StatusUnauthorized
	}

	return http.StatusOK
}
