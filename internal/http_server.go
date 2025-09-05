package internal

import (
	"log"
	"net/http"
	"strings"
)

func SetupHttpServer() {
	http.Handle("/", http.FileServer(http.Dir("./resources/assets/public")))
	http.HandleFunc("/api/v1/", func(response http.ResponseWriter, request *http.Request) {
		handlePathing(strings.ReplaceAll(request.URL.String(), "/api/v1/", ""), response, request)
	})

	registerEndpoints()

	err := http.ListenAndServeTLS(":8080", "secrets/public.cert", "secrets/private.key", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func handlePathing(urlPart string, rsp http.ResponseWriter, rqst *http.Request) {
	result := endpoint_handlers[urlPart]
	if result == nil {
		rsp.WriteHeader(http.StatusBadRequest)
	}

	rsp.WriteHeader(result(urlPart, rsp, rqst))
}
