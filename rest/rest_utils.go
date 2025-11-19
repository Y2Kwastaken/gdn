package rest

import (
	"net/http"

	"github.com/Y2Kwastaken/gdn/internal"
)

type FileStore = internal.FileStore
type IdResponse = internal.IdResponse
type Metadata = internal.Metadata

func wstd(rspn http.ResponseWriter, code int) {
	rspn.WriteHeader(code)
}

func werr(rspn http.ResponseWriter, code int) {
	http.Error(rspn, http.StatusText(code), code)
}
