# A Simple Website Analyzer

---

## Building and Running

The solution consists of a single Go service without any other components.
It can therefore be built and run simply by doing the following from the root of the repo:
```
go build ./cmd/main.go
./main
```

Or just:
```
go run ./cmd/main.go
```

You may also need to fetch the dependencies with e.g. `go mod tidy`.

By default, it will start listening on [`localhost:8080`](http://localhost:8080), which can be modified through the configuration file in `./cfg/config.yaml`.

## Decisions made

# HTML Verision Detection

The closest thing to an explicit version marker that an HTML document has is the DOCTYPE, so we display that.
Very old documents may not have a DOCTYPE at all, making it impossible to figure out the intended HTML version.
In this case, we could try to guess the version from the HTML elements included on the page, but that's not ideal either.
In a real scenario, I would probably talk to the stakeholder and try to figure out *why* they need to know the HTML version.
Chances are, they probably actually don't, and what they might want instead is something like browser compatibility, which should be figured out based on individual features supported by the browser, rather than just the HTML version.

# Inaccessble links
An "inaccessible" link for us is one that does not respond successfully to a GET, either for transport-level reasons (timeout, unexpected connection closure, etc.), or by returning a response that is not in the 200 range (after the redirects are followed).

If the same URL (modulo fragment part) appears on a page multiple times, it will be counted as multiple links but will only be fetched once.

# Login form detection

For this, we simply look at whether the page contains an input field for a password.
This is already enough for most cases, but could lead to false positives (like sign-up pages, which also usually have a password field)

A more clever heursitic may or may not be required, depending on the intended use.
For this, we could look at e.g. form name, whether the page has words like "log in" or "sign in" or variations thereof, whether the form also has an email/phone/username field and so on.
These could also be taken with different weights (which could even be learned with ML, say, a Naive Bayes classifier) to produce a confidence value instead of just true/false, but this is very likely overkill.  

## Possible improvements

I unfortunately got somewhat sick while working on this project, and thus didn't have as much time as I would have liked to work on it, so there are quite a few of these.
For example, there are very few comments; I was planning to add GoDoc comments once I've settled on a final architecture, but ran out of time for that.

Some relatively low-hanging fruit would be to implement graceful shutdown, as the application is already using `context.Context` for its requests.

As mentioned above, the analysis heuristics are quite basic and could be improved as well.

The HTML template-based approach is also very basic: in a real scenario, this would probably be a microservice providing an API, and the page would be an SPA consuming it.
If we're feeling really fancy, the API could even provide a server-push indication of progress (how many links have been checked), so that the user can be shown something like a progress bar.