package internal

import "net/http"

var endpoint_handlers = make(map[string]func(string, http.ResponseWriter, *http.Request) int)

func registerEndpoints() {

}
