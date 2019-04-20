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

type CommentService struct {
	commentRepository *repository.CommentRepository
}

func NewCommentService(commentRepository *repository.CommentRepository) (*CommentService) {
	return &CommentService{
		commentRepository: commentRepository,
	}
}

func (service *CommentService) GetByMediaID(id string) ([]model.Comment, error) {
	return service.commentRepository.GetByMediaID(id)
}

func (service *CommentService) GetByID(id string) (*model.Comment, error) {
	return service.commentRepository.GetByID(id)
}

func (service *CommentService) Create(id string, req CommentRequest) (error) {
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
		MediaID: id,
		ParentCommentID: *parentCommentId,
		Comment: req.Comment,
	})
}

func (service *CommentService) Save(comment *model.Comment) (error) {
	return service.commentRepository.Save(comment)
}
