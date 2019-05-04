package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"strings"
	"fmt"
	"regexp"
)

type SelectPage struct {
	Page uint
	Count uint
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap   = regexp.MustCompile("([a-z0-9])([A-Z])")

func GetSelectPage(c *gin.Context) (*SelectPage) {
	params := c.Request.URL.Query()

	selectPage := &SelectPage{
		Page: 0,
		Count: 50,
	}

	if page, ok := params["page"]; ok {
		parsed, err := strconv.ParseUint(page[0], 10, 32)
		if err == nil {
			selectPage.Page = uint(parsed)
		}
	}
	if count, ok := params["count"]; ok {
		parsed, err := strconv.ParseUint(count[0], 10, 32)
		if err == nil {
			selectPage.Count = uint(parsed)
		}
	}

	// set a max cap for the number of items to return per page
	if selectPage.Count > 50 {
		selectPage.Count = 50
	}

	return selectPage
}

func Redirect(c *gin.Context, path string) {
	statusCode := http.StatusFound

	// append an error if we set one so it can be picked up and parsed in the next route
	if val, ok := c.Get("error"); ok {
		if strings.Index(path, "?") == -1 {
			path = fmt.Sprintf("%s?error=%s", path, val)
		} else {
			path = fmt.Sprintf("%s&error=%s", path, val)
		}
	}


	if code, ok := c.Get("statusCode"); ok {
		statusCode = code.(int)
	}

	c.Redirect(statusCode, path)
	c.Abort()
}

func RedirectWithError(c *gin.Context, path string, code int, error string) {
	c.Set("statusCode", code)
	c.Set("error", error)
	Redirect(c, path)
}

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake  = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
