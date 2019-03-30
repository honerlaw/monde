package util

import (
	"github.com/gin-gonic/gin"
	"strconv"
)

type SelectPage struct {
	Page int
	Count int
}

func GetSelectPage(c *gin.Context) (*SelectPage) {
	params := c.Request.URL.Query()

	selectPage := &SelectPage{
		Page: 0,
		Count: 50,
	}

	if page, ok := params["page"]; ok {
		parsed, err := strconv.ParseInt(page[0], 10, 32)
		if err == nil {
			selectPage.Page = int(parsed)
		}
	}
	if count, ok := params["count"]; ok {
		parsed, err := strconv.ParseInt(count[0], 10, 32)
		if err == nil {
			selectPage.Count = int(parsed)
		}
	}

	// set a max cap for the number of items to return per page
	if selectPage.Count > 50 {
		selectPage.Count = 50
	}

	return selectPage
}
