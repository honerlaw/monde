package render

import (
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
	"log"
)

type ReactMeta struct {
	Props      map[string]interface{}
	StatusCode int
}

type Pool interface {
	get() *JSVM
	put(*JSVM)
	drop(*JSVM)
}

type React struct {
	Pool
	debug bool
	path  string
}

// NewReact initialized React struct
func NewReact(filePath string, debug bool, proxy http.Handler) *React {
	r := &React{
		debug: debug,
		path:  filePath,
	}
	if !debug {
		r.Pool = NewFixedPool(filePath, runtime.NumCPU(), proxy)
	} else {
		// Use onDemandPool to load full react
		// app each time for any http requests.
		// Useful to debug the app.
		r.Pool = &OnDemandPool{
			path:  filePath,
			proxy: proxy,
		}
	}
	return r
}

// Handle handles all HTTP requests which
// have not been caught via static file
// handler or other middlewares.
func (r *React) Handle(c *gin.Context) error {
	defer func() {
		if r := recover(); r != nil {
			SendHtml(c, &SendHtmlOptions{
				StatusCode: http.StatusInternalServerError,
				Html: "Internal Server Error",
			})
		}
	}()

	vm := r.get()
	reactMeta := c.MustGet("react-meta").(ReactMeta)

	start := time.Now()
	select {
	case re := <-vm.Handle(map[string]interface{}{
		"url":     c.Request.URL.String(),
		"headers": c.Request.Header,
		"props":   reactMeta.Props,
	}):
		r.put(vm)

		if len(re.Error) == 0 {
			return SendHtml(c, &SendHtmlOptions{
				StatusCode: reactMeta.StatusCode,
				Html: re.Html,
			})
		} else if len(re.Error) != 0 {
			log.Print(re.Error)

			return SendHtml(c,  &SendHtmlOptions{
				Headers: HtmlHeaders{
					"X-RENDER-TIME": time.Since(start).String(),
				},
				StatusCode: http.StatusInternalServerError,
				Html: "Internal Server Error",
			})
		}
	case <-time.After(2 * time.Second):
		r.drop(vm)
		return SendHtml(c, &SendHtmlOptions{
			StatusCode: http.StatusInternalServerError,
			Html: "Timeout while rendering",
		})
	}
	return nil
}
