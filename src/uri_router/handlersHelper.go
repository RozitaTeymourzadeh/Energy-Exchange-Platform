package uri_router

import "net/http"

//for return code 503
func returnCode503(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Error", http.StatusServiceUnavailable)
	http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
}

//for return code 500
func returnCode500(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Server Error", http.StatusInternalServerError)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

//for return code 204
func returnCode204(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Block does not exists", http.StatusNoContent)
	http.Error(w, http.StatusText(http.StatusNoContent), http.StatusNoContent)
}
