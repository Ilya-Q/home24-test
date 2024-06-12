package check

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
)

type LinkCounts struct {
	InternalLinks     uint64
	ExternalLinks     uint64
	InaccessibleLinks uint64
}

func dedupeLinks(counts *LinkCounts, baseUrl *url.URL, links []string) map[string]uint64 {
	deduped := make(map[string]uint64)
	for _, link := range links {
		parsedLink, err := baseUrl.Parse(link)
		if !strings.HasPrefix(parsedLink.Scheme, "http") {
			continue
		}
		if err != nil {
			log.Printf("Link '%s' couldn't be parsed: %v", link, err)
			counts.InaccessibleLinks++
			continue
		}
		if parsedLink.Host == baseUrl.Host {
			counts.InternalLinks++
		} else {
			counts.ExternalLinks++
		}
		parsedLink.Fragment = ""
		deduped[parsedLink.String()]++
	}
	return deduped
}

func CheckLinks(ctx context.Context, baseUrl *url.URL, links []string) LinkCounts {
	var counts LinkCounts
	deduped := dedupeLinks(&counts, baseUrl, links)

	var inaccessible atomic.Uint64
	var wg sync.WaitGroup
	for l, c := range deduped {
		wg.Add(1)
		go func() {
			defer wg.Done()
			// it's okay to capture l and c in this closure
			// because Go finally fixed the loop variable gotcha in 1.22
			req, err := http.NewRequestWithContext(ctx, http.MethodGet, l, nil)
			if err != nil {
				log.Printf("Couldn't create request for '%s': %v", l, err)
				inaccessible.Add(c)
				return
			}
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Printf("GET failed for '%s': %v", l, err)
				inaccessible.Add(c)
				return
			}
			// we only care about the status code, not the body
			// so let's close it immediately, without a defer
			resp.Body.Close()
			if resp.StatusCode/100 != 2 {
				log.Printf("Non-200 response for '%s': %v", l, resp.Status)
				inaccessible.Add(c)
				return
			}

		}()
	}
	wg.Wait()
	counts.InaccessibleLinks = inaccessible.Load()
	return counts
}
