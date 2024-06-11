package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/html"

	"github.com/Ilya-Q/home24-test/internal/analyze"
	"github.com/gorilla/mux"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	url, err := io.ReadAll(r.Body)

	if err != nil {
		fmt.Fprintf(w, "%v\n", err)
		return
	}

	target, _ := http.Get(string(url))

	root, _ := html.Parse(target.Body)

	visitors := []analyze.HTMLVisitor{
		new(analyze.TitleGetter),
		new(analyze.DoctypeGetter),
		new(analyze.HeadingCounter),
	}
	analyze.Walk(root, visitors)
	for _, v := range visitors {
		fmt.Fprintf(w, "%+v\n", v)
	}

}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", TestHandler)

	http.ListenAndServe("localhost:8080", r)
}
