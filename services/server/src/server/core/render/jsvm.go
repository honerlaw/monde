package render

import (
	crand "crypto/rand"
	"net/http"
	"github.com/dop251/goja_nodejs/eventloop"
	"github.com/olebedev/gojax/fetch"
	"io/ioutil"
	"github.com/dop251/goja"
	"encoding/binary"
	"math/rand"
	"github.com/fatih/structs"
	"log"
	"reflect"
	"encoding/json"
)

// denotes the expected response from the global main function in bundle.js
// this should always be kept in sync with the type definition in the bundle
type Resp struct {
	Error string `json:"error"`
	Html  string `json:"html"`
}

// JSVM wraps goja EventLoop
type JSVM struct {
	*eventloop.EventLoop
	ch chan Resp
	fn goja.Callable
}

// newJSVM loads bundle.js into context.
func newJSVM(filePath string, proxy http.Handler) *JSVM {
	vm := &JSVM{
		EventLoop: eventloop.NewEventLoop(),
		ch:        make(chan Resp, 1),
	}

	vm.EventLoop.Start()
	err := fetch.Enable(vm.EventLoop, proxy)

	bundle, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	vm.EventLoop.RunOnLoop(func(_vm *goja.Runtime) {
		var seed int64
		if err := binary.Read(crand.Reader, binary.LittleEndian, &seed); err != nil {
			log.Panicf("Could not read random bytes: %v", err)
		}
		_vm.SetRandSource(goja.RandSource(rand.New(rand.NewSource(seed)).Float64))

		_, err := _vm.RunScript("bundle.js", string(bundle))
		if err != nil {
			panic(err)
		}

		if fn, ok := goja.AssertFunction(_vm.Get("main")); ok {
			vm.fn = fn
		} else {
			log.Println("Could not find main global function in bundle.js")
			panic(err)
		}

		_vm.Set("__goServerCallback__", func(call goja.FunctionCall) goja.Value {
			obj := call.Argument(0).Export().(map[string]interface{})
			re := &Resp{}
			for _, field := range structs.Fields(re) {
				if n := field.Tag("json"); len(n) > 1 {
					field.Set(obj[n])
				}
			}
			vm.ch <- *re
			return nil
		})
	})

	return vm
}

// Handle handles http requests
func (r *JSVM) Handle(req map[string]interface{}) <-chan Resp {
	r.EventLoop.RunOnLoop(func(vm *goja.Runtime) {
		// @todo actually debug this and potentially remove it
		data, _ := json.Marshal(req) // convert to json string, weird issues with ToValue otherwise
		r.fn(nil, vm.ToValue(string(data)), vm.ToValue("__goServerCallback__"))
	})
	return r.ch
}

func (r *JSVM) recursiveToValue(vm *goja.Runtime, req map[string]interface{}) goja.Value  {
	newmap := make(map[string]interface{})
	for k, v := range req {
		if reflect.ValueOf(req[k]).Kind() == reflect.Map {
			val, ok := req[k].(map[string]interface{})
			if ok {
				newmap[k] = r.recursiveToValue(vm, val)
			} else {
				newmap[k] = vm.ToValue(val)
			}
		} else {
			newmap[k] = vm.ToValue(v)
		}
	}

	return vm.ToValue(newmap)
}
