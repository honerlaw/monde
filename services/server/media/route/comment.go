package route

import (
	"github.com/gin-gonic/gin"
	"services/server/media/service"
	"services/server/core/util"
	"services/server/user/middleware"
)

type CommentRoute struct {
	commentService *service.CommentService
}

func NewCommentRoute(commentService *service.CommentService) (*CommentRoute) {
	return &CommentRoute{
		commentService: commentService,
	}
}

func (route *CommentRoute) Post(c *gin.Context) {
	id, ok := c.Params.Get("id")
	if !ok {
		util.Redirect(c, "/");
		return;
	}

	var req service.CommentRequest
	if err := c.ShouldBind(&req); err != nil {
		c.Set("error", err.Error())
		util.Redirect(c, "/media/view/"+id)
		return
	}

	payload := c.MustGet("JWT_AUTH_PAYLOAD")
	if payload == nil {
		util.Redirect(c, "/")
		return
	}
	auth := payload.(*middleware.AuthPayload)

	err := route.commentService.Create(id, auth.ID, &req)
	if err != nil {
		c.Set("error", err.Error())
	}

	util.Redirect(c, "/media/view/"+id)
}
