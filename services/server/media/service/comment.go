package service

import (
	"services/server/media/repository"
	"services/server/media/model"
	"strings"
)

type CommentRequest struct {
	Comment         string `form:"comment"`
	ParentCommentID string `form:"parent_comment_id"`
}

type CommentResponse struct {
	model.Comment
	Children []CommentResponse `json:"children"`
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

	resultMap := make(map[string]CommentResponse)

	for _, comment := range comments {

		// does not exist in map, so add it
		if _, ok := resultMap[comment.ID]; !ok {
			resp := CommentResponse{
				Children: make([]CommentResponse, 0),
			}
			resp.ID = comment.ID
			resp.CreatedAt = comment.CreatedAt
			resp.UpdatedAt = comment.UpdatedAt
			resp.DeletedAt = comment.DeletedAt
			resp.ParentCommentID = comment.ParentCommentID
			resp.MediaID = resp.MediaID
			resp.Comment = resp.Comment
		}
	}

	// build the hierarchy
	for _, pcom := range comments {
		comment, _ := resultMap[pcom.ID]
		if parentComment, ok := resultMap[pcom.ParentCommentID]; ok {
			parentComment.Children = append(parentComment.Children, comment)

			// this comment, became a child, so remove it from the map entirely
			// this basically leaves only the "roots" that don't have parents in the map
			delete(resultMap, pcom.ID)
		}
	}

	// convert the map to an array of values
	roots := make([]CommentResponse, len(resultMap))
	for _, resp := range resultMap {
		roots = append(roots, resp)
	}

	return roots, nil
}

func (service *CommentService) GetByID(id string) (*model.Comment, error) {
	return service.commentRepository.GetByID(id)
}

func (service *CommentService) Create(id string, userID string, req CommentRequest) (error) {
	var parentCommentId *string
	if len(strings.TrimSpace(req.ParentCommentID)) == 0 {
		parentCommentId = nil
	} else {
		comment, _ := service.GetByID(*parentCommentId)
		if comment != nil {
			parentCommentId = &comment.ID
		}
	}

	return service.Save(&model.Comment{
		MediaID:         id,
		UserID:          userID,
		ParentCommentID: *parentCommentId,
		Comment:         req.Comment,
	})
}

func (service *CommentService) Save(comment *model.Comment) (error) {
	return service.commentRepository.Save(comment)
}
