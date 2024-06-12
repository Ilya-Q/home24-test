package handlers

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"

	"github.com/Ilya-Q/home24-test/internal/analyze"
	"github.com/Ilya-Q/home24-test/internal/check"
	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
)

type AnalysisPageResults struct {
	URL           string
	Doctype       string
	Title         string
	HeadingCounts analyze.HeadingCounter
	HasLoginForm  bool
	LinkCounts    check.LinkCounts
}

var analysisPageTemplate = template.Must(template.ParseFiles("html/templates/analysis.html"))

func AnalysisFormHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		writeError(w,
			http.StatusBadRequest,
			fmt.Sprintf("Form could not be parsed: %v", err),
		)
		return
	}
	url := r.Form.Get("url")
	if url == "" {
		writeError(w,
			http.StatusBadRequest,
			"URL not passed or empty",
		)
		return
	}
	req, err := http.NewRequestWithContext(r.Context(), http.MethodGet, url, nil)
	if err != nil {
		writeError(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Couldn't create request: %v", err),
		)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		writeError(w,
			http.StatusBadRequest, // Not sure if this is really a 400
			fmt.Sprintf("URL '%s' not reachable: %v", url, err),
		)
		return
	}

	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		writeError(w,
			http.StatusBadRequest,
			fmt.Sprintf("URL '%s' returned an invalid Content-Type '%s': %v",
				url,
				resp.Header.Get("Content-Type"),
				err,
			),
		)
		return
	}
	if contentType != "text/html" {
		writeError(w,
			http.StatusBadRequest,
			fmt.Sprintf("URL '%s' does not point to an HTML resource", url),
		)
		return
	}

	br, err := charset.NewReader(resp.Body, resp.Header.Get("Content-Type"))
	if err != nil {
		writeError(w,
			http.StatusInternalServerError,
			fmt.Sprintf("Couldn't create encoding sniffer: %v", err),
		)
		return
	}

	root, err := html.Parse(br)
	if err != nil {
		writeError(w,
			http.StatusBadRequest,
			fmt.Sprintf("HTML could not be parsed: %v", err),
		)
		return
	}

	tg := new(analyze.TitleGetter)
	dt := new(analyze.DoctypeGetter)
	hc := new(analyze.HeadingCounter)
	ex := new(analyze.LinkExtractor)
	lfd := new(analyze.LoginFormDetector)
	analyze.Walk(root, []analyze.HTMLVisitor{tg, dt, hc, ex, lfd})

	lc := check.CheckLinks(r.Context(), resp.Request.URL, ex.Links)

	analysisPageTemplate.Execute(w, &AnalysisPageResults{
		URL:           url,
		Doctype:       dt.Doctype,
		Title:         tg.Title,
		HeadingCounts: *hc,
		HasLoginForm:  lfd.LoginFormFound,
		LinkCounts:    lc,
	})
}
