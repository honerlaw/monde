package render

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HtmlHeaders map[string]string

type SendHtmlOptions struct {
	Headers    HtmlHeaders
	StatusCode int
	Html       string
}

func SendHtml(c *gin.Context, options *SendHtmlOptions) (error) {
	c.Writer.WriteHeader(options.StatusCode)
	c.Writer.Header().Set("Content-Type", "text/html")

	// add the headers if they exist, do it after so previous headers can be overwritten
	if options.Headers != nil {
		for key, value := range options.Headers {
			c.Writer.Header().Set(key, value)
		}
	}

	_, err := c.Writer.Write([]byte(options.Html))
	return err;
}

func RenderPage(c *gin.Context, statusCode int, props gin.H) {
	if props == nil {
		props = gin.H{}
	}

	// always add the auth payload to the props, whether it exists or not
	props["authPayload"], _ = c.Get("JWT_AUTH_PAYLOAD")
	props["statusCode"] = statusCode

	c.Set("react-meta", ReactMeta{
		StatusCode: statusCode,
		Props:      props,
	})
}

func RenderNoRoute(c *gin.Context) {
	RenderPage(c, http.StatusNotFound, nil)
}
