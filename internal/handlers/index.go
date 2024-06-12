package handlers

import "net/http"

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "html/static/index.html")
}
