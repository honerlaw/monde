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

		// on demand initializes a new JSVM for each request
		r.Pool = &OnDemandPool{
			path:  filePath,
			proxy: proxy,
		}
	}
	return r
}

// @todo use some html templates for the 500 error
func (r *React) Handle(c *gin.Context) error {
	defer func() {
		if r := recover(); r != nil {
			SendHtml(c, &SendHtmlOptions{
				StatusCode: http.StatusInternalServerError,
				Html:       "Internal Server Error",
			})
		}
	}()

	vm := r.get()
	reactMeta := c.MustGet("react-meta").(ReactMeta)

	// set the error no the meta props, if there is an error to set
	setError(c, &reactMeta)

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
				Html:       re.Html,
			})
		} else if len(re.Error) != 0 {
			log.Print(re.Error)

			return SendHtml(c, &SendHtmlOptions{
				Headers: HtmlHeaders{
					"X-RENDER-TIME": time.Since(start).String(),
				},
				StatusCode: http.StatusInternalServerError,
				Html:       "Internal Server Error",
			})
		}
	case <-time.After(2 * time.Second):
		r.drop(vm)
		return SendHtml(c, &SendHtmlOptions{
			StatusCode: http.StatusInternalServerError,
			Html:       "Timeout while rendering",
		})
	}
	return nil
}

func setError(c *gin.Context, meta *ReactMeta) {

	// there is already an error set so do nothing
	if _, ok := meta.Props["error"]; ok {
		return
	}

	// set the error from the context if it exists
	if err, ok := c.Get("error"); ok {
		meta.Props["error"] = err
		return
	}

	// last if there is an error in the url query params (e.g. we redirected with an error)
	params := c.Request.URL.Query()
	if err, ok := params["error"]; ok {
		meta.Props["error"] = err
	}
}
