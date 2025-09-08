package internal

import (
	"log"
	"net/http"
	"strings"
)

func SetupHttpServer(store *FileStore) {

	http.Handle("/", http.FileServer(http.Dir("./resources/assets/public")))
	http.HandleFunc("/api/v1/", func(response http.ResponseWriter, request *http.Request) {
		handlePathing(store, strings.ReplaceAll(request.URL.String(), "/api/v1/", ""), response, request)
	})

	registerEndpoints()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handlePathing(store *FileStore, urlPart string, rsp http.ResponseWriter, rqst *http.Request) {
	result, ok := endpoint_handlers[urlPart]
	if !ok {
		rsp.WriteHeader(http.StatusBadRequest)
		return
	}

	rsp.WriteHeader(result(store, urlPart, rsp, rqst))
}
