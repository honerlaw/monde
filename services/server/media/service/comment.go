package service

import (
	"services/server/media/repository"
	"services/server/media/model"
	"strings"
	"time"
	"github.com/pkg/errors"
)

type CommentRequest struct {
	Comment         string `form:"comment"`
	ParentCommentID string `form:"parent_comment_id"`
}

type CommentResponse struct {
	ID              string            `json:"id"`
	UserID          string            `json:"user_id"`
	MediaID         string            `json:"media_id"`
	ParentCommentID string            `json:"parent_comment_id"`
	Comment         string            `json:"comment"`
	CreatedAt       time.Time         `json:"created_at"`
	Children        []CommentResponse `json:"children"`
}

type CommentService struct {
	commentRepository *repository.CommentRepository
}

func NewCommentService(commentRepository *repository.CommentRepository) (*CommentService) {
	return &CommentService{
		commentRepository: commentRepository,
	}
}

// @todo we need to fetch user info from the user id
func (service *CommentService) GetByMediaID(id string) ([]CommentResponse, error) {
	comments, err := service.commentRepository.GetByMediaID(id)
	if err != nil {
		return nil, err;
	}

	// @todo this does not keep order
	resultMap := make(map[string]*CommentResponse)

	for _, comment := range comments {

		// does not exist in map, so add it
		if _, ok := resultMap[comment.ID]; !ok {
			resp := CommentResponse{
				Children: make([]CommentResponse, 0),
			}
			resp.ID = comment.ID
			resp.CreatedAt = comment.CreatedAt
			resp.UserID = comment.UserID
			resp.ParentCommentID = comment.ParentCommentID
			resp.MediaID = comment.MediaID
			resp.Comment = comment.Comment

			resultMap[comment.ID] = &resp
		}
	}

	for _, pcom := range comments {
		comment, _ := resultMap[pcom.ID]
		if parentComment, ok := resultMap[pcom.ParentCommentID]; ok {
			parentComment.Children = append(parentComment.Children, *comment)

			// this comment, became a child, so remove it from the map entirely
			// this basically leaves only the "roots" that don't have parents in the map
			delete(resultMap, pcom.ID)
		}
	}

	// convert the map to an array of values
	roots := make([]CommentResponse, 0)
	for _, resp := range resultMap {
		roots = append(roots, *resp)
	}

	return roots, nil
}

func (service *CommentService) GetByID(id string) (*model.Comment, error) {
	return service.commentRepository.GetByID(id)
}

func (service *CommentService) Create(id string, userID string, req *CommentRequest) (error) {
	comment := model.Comment{
		MediaID: id,
		UserID:  userID,
		Comment: req.Comment,
	}

	if len(strings.TrimSpace(req.ParentCommentID)) != 0 {
		parentComment, _ := service.GetByID(strings.TrimSpace(req.ParentCommentID))
		if parentComment != nil {
			comment.ParentCommentID = parentComment.ID
		}
	}

	if len(strings.TrimSpace(req.Comment)) == 0 {
		return errors.New("comment can not be empty")
	}

	return service.Save(&comment)
}

func (service *CommentService) Save(comment *model.Comment) (error) {
	return service.commentRepository.Save(comment)
}
