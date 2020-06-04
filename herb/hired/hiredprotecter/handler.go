package hiredprotecter

import (
	"net/http"
)

type HandlerConfig struct {
	StatusCode int
	RedirectTo string
	Headers    http.Header
	Body       string
}

func (c *HandlerConfig) CreateHandler() (http.Handler, error) {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body = ""
		if c.RedirectTo != "" {
			http.Redirect(w, r, c.RedirectTo, c.StatusCode)
			return
		}
		if c.Headers != nil {
			for name := range c.Headers {
				for k := range c.Headers[name] {
					w.Header().Add(name, c.Headers[name][k])
				}
			}
		}
		if c.StatusCode != 0 {
			w.WriteHeader(c.StatusCode)
		}

		if c.Body == "" {
			body = http.StatusText(c.StatusCode)
		} else {
			body = c.Body
		}
		_, err := w.Write([]byte(body))
		if err != nil {
			panic(err)
		}
	}), nil
}
