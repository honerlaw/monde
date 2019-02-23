package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/musawirali/preact-rpc/goclient"
)

func main() {

	// connect to the react rendering server
	goclient.Connect("tcp", "0.0.0.0:9000");

	r := gin.Default()
	r.GET("/", func(c *gin.Context) {

		resp, err := goclient.RenderComponent("Index", nil, map[string](interface{}){
			"toWhat": "Universe",
		})
		if err != nil {
			panic(err)
		}

		c.Writer.WriteHeader(http.StatusOK)
		c.Writer.Write([]byte(resp.Html));
	})
	r.Run() // listen and serve on 0.0s.0.0:8080
}
