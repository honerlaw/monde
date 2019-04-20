package route

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"services/server/media/service"
	"services/server/core/render"
	"services/server/core/util"
	"services/server/user/middleware"
	"services/server/core/service/aws"
	"services/server/media/repository"
	"strings"
)

type ListRoute struct {
	mediaService *service.MediaService
}

func NewListRoute(mediaService *service.MediaService) (*ListRoute) {
	return &ListRoute{
		mediaService: mediaService,
	}
}

func (route *ListRoute) Get(c *gin.Context) {
	payload := c.MustGet("JWT_AUTH_PAYLOAD")

	if payload == nil {
		util.Redirect(c, "/")
		return
	}

	uploads := []service.MediaResponse{}

	// fetch requested media info for given page
	data, err := route.mediaService.GetByUserID(payload.(*middleware.AuthPayload).ID, util.GetSelectPage(c))
	if err != nil {
		render.RenderPage(c, http.StatusInternalServerError, nil)
		return
	}

	// there is a chance that the lambda has not started the job processing yet, so the media info won't exist
	// in those cases, we should append a pending upload in its place
	pending := route.getPendingUploadIfNeeded(c, data)
	if pending != nil {
		uploads = append(uploads, *pending)
	}

	uploads = append(uploads, route.mediaService.ConvertMediaData(data, func(datum *repository.MediaData, resp *service.MediaResponse) {
		resp.TranscodingStatus = aws.GetETService().GetJobStatus(datum.Media.JobID)
		resp.CanPublish = datum.Media.CanPublish()
		resp.IsPublished = datum.Media.Published
	})...)

	props := gin.H{
		"uploads": uploads,
	}

	render.RenderPage(c, http.StatusOK, props)
}


func (route *ListRoute) getPendingUploadIfNeeded(c *gin.Context, data []repository.MediaData) (*service.MediaResponse) {
	params := c.Request.URL.Query()
	bucket, okBucket := params["bucket"]
	key, okKey := params["key"]
	if okBucket && okKey {
		pieces := strings.Split(key[0], "/")
		videoId := pieces[len(pieces)-1]
		canAddPending := len(data) == 0 || data[0].Media.ID != videoId

		// basically, we don't have the latest info from the trannscoder, but the file was definitely uploaded
		// so we should append the info anyways...
		if canAddPending && aws.GetS3Service().FileExists(bucket[0], key[0]) {
			return &service.MediaResponse{
				ID:                videoId,
				CanPublish:        false,
				Title:             "",
				Description:       "",
				TranscodingStatus: "pending",
				Hashtags:          []string{},
				IsPublished:       false,
				Thumbnails:        []string{},
				Videos:            []service.MediaVideoResponse{},
			}
		}
	}
	return nil;
}