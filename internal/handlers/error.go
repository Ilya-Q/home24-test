package handlers

import (
	"fmt"
	"html/template"
	"net/http"
)

var errorPageTemplate = template.Must(template.ParseFiles("html/templates/error.html"))

type ErrorPageInfo struct {
	Status  string
	Message string
}

func writeError(w http.ResponseWriter, statusCode int, message string) {
	w.WriteHeader(statusCode)
	errorPageTemplate.Execute(w, &ErrorPageInfo{
		fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		message,
	})
}

func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w,
		http.StatusNotFound,
		fmt.Sprintf("Path '%s' is not present on this server", r.URL),
	)
}

func MethodNotAllowedHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w,
		http.StatusMethodNotAllowed,
		fmt.Sprintf("Method '%s' is not allowed for path '%s'", r.Method, r.URL),
	)
}
