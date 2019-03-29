package render

import "net/http"

type OnDemandPool struct {
	Pool
	path  string
	proxy http.Handler
}

func (f *OnDemandPool) get() *JSVM {
	return newJSVM(f.path, f.proxy)
}

func (f OnDemandPool) put(c *JSVM) {
	c.Stop()
}

func (f *OnDemandPool) drop(c *JSVM) {
	f.put(c)
}

