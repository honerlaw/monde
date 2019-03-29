package render

import (
	"net/http"
	"log"
)

type FixedPool struct {
	Pool
	ch    chan *JSVM
	path  string
	proxy http.Handler
}

// newEnginePool return pool of JS vms.
func NewFixedPool(filePath string, size int, proxy http.Handler) *FixedPool {
	log.Printf("initializing fixed pool. size: %d, path: %s", size, filePath)
	pool := &FixedPool{
		path:  filePath,
		ch:    make(chan *JSVM, size),
		proxy: proxy,
	}

	go func() {
		for i := 0; i < size; i++ {
			pool.ch <- newJSVM(filePath, proxy)
		}

		log.Print("finished initializing pool")
	}()

	return pool
}

func (o *FixedPool) get() *JSVM {
	return <-o.ch
}

func (o *FixedPool) put(ot *JSVM) {
	o.ch <- ot
}

func (o *FixedPool) drop(ot *JSVM) {
	ot.Stop()
	ot = nil
	o.ch <- newJSVM(o.path, o.proxy)
}
