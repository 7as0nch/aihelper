package service

import (
	"context"
	"strconv"
	"time"

	"github.com/example/aichat/backend/api/userfeedback/v1"
	"github.com/example/aichat/backend/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
)

type UserFeedbackService struct {
	v1.UnimplementedUserFeedbackServer
	uc  *biz.UserFeedbackUsecase
	log *log.Helper
}

func NewUserFeedbackService(uc *biz.UserFeedbackUsecase, logger log.Logger) *UserFeedbackService {
	return &UserFeedbackService{
		uc:  uc,
		log: log.NewHelper(logger),
	}
}

func (s *UserFeedbackService) CreateFeedback(ctx context.Context, req *v1.CreateFeedbackRequest) (*v1.CreateFeedbackReply, error) {
	feedback := &biz.UserFeedback{
		// UserID:  req.UserId,
		Content: req.Content,
		Type:    req.Type,
		Status:  req.Status,
		Rating:  req.Rating,
	}
	
	createdFeedback, err := s.uc.CreateFeedback(ctx, feedback)
	if err != nil {
		s.log.Errorf("Failed to create feedback: %v", err)
		return nil, err
	}
	
	return &v1.CreateFeedbackReply{
		Feedback: &v1.UserFeedbackData{
			// Id:        createdFeedback.ID,
			// UserId:    createdFeedback.UserID,
			Content:   createdFeedback.Content,
			Type:      createdFeedback.Type,
			Status:    createdFeedback.Status,
			Rating:    createdFeedback.Rating,
			CreatedAt: createdFeedback.CreatedAt.Format(time.RFC3339),
			UpdatedAt: createdFeedback.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserFeedbackService) GetFeedback(ctx context.Context, req *v1.GetFeedbackRequest) (*v1.GetFeedbackReply, error) {
	feedback, err := s.uc.GetFeedback(ctx, uint64(req.Id))
	if err != nil {
		s.log.Errorf("Failed to get feedback: %v", err)
		return nil, err
	}
	
	return &v1.GetFeedbackReply{
		Feedback: &v1.UserFeedbackData{
			// Id:        feedback.ID,
			// UserId:    feedback.UserID,
			Content:   feedback.Content,
			Type:      feedback.Type,
			Status:    feedback.Status,
			Rating:    feedback.Rating,
			CreatedAt: feedback.CreatedAt.Format(time.RFC3339),
			UpdatedAt: feedback.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserFeedbackService) ListFeedbacks(ctx context.Context, req *v1.ListFeedbacksRequest) (*v1.ListFeedbacksReply, error) {
	feedbacks, total, err := s.uc.ListFeedbacks(ctx, strconv.FormatUint(uint64(req.UserId), 10), req.Page, req.PageSize)
	if err != nil {
		s.log.Errorf("Failed to list feedbacks: %v", err)
		return nil, err
	}
	
	var feedbackList []*v1.UserFeedbackData
	for _, feedback := range feedbacks {
		feedbackList = append(feedbackList, &v1.UserFeedbackData{
			// Id:        feedback.ID,
			// UserId:    feedback.UserID,
			Content:   feedback.Content,
			Type:      feedback.Type,
			Status:    feedback.Status,
			Rating:    feedback.Rating,
			CreatedAt: feedback.CreatedAt.Format(time.RFC3339),
			UpdatedAt: feedback.UpdatedAt.Format(time.RFC3339),
		})
	}
	
	return &v1.ListFeedbacksReply{
		Feedbacks: feedbackList,
		Total:     total,
	}, nil
}

func (s *UserFeedbackService) UpdateFeedback(ctx context.Context, req *v1.UpdateFeedbackRequest) (*v1.UpdateFeedbackReply, error) {
	feedback := &biz.UserFeedback{
		// ID:      req.Id,
		// UserID:  req.UserId,
		Content: req.Content,
		Type:    req.Type,
		Status:  req.Status,
		Rating:  req.Rating,
	}
	
	updatedFeedback, err := s.uc.UpdateFeedback(ctx, feedback)
	if err != nil {
		s.log.Errorf("Failed to update feedback: %v", err)
		return nil, err
	}
	
	return &v1.UpdateFeedbackReply{
		Feedback: &v1.UserFeedbackData{
			// Id:        updatedFeedback.ID,
			// UserId:    updatedFeedback.UserID,
			Content:   updatedFeedback.Content,
			Type:      updatedFeedback.Type,
			Status:    updatedFeedback.Status,
			Rating:    updatedFeedback.Rating,
			CreatedAt: updatedFeedback.CreatedAt.Format(time.RFC3339),
			UpdatedAt: updatedFeedback.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *UserFeedbackService) DeleteFeedback(ctx context.Context, req *v1.DeleteFeedbackRequest) (*v1.DeleteFeedbackReply, error) {
	err := s.uc.DeleteFeedback(ctx, uint64(req.Id))
	if err != nil {
		s.log.Errorf("Failed to delete feedback: %v", err)
		return nil, err
	}
	
	return &v1.DeleteFeedbackReply{
		Success: true,
	}, nil
}

func (s *UserFeedbackService) ListFeedbacksByStatus(ctx context.Context, req *v1.ListFeedbacksByStatusRequest) (*v1.ListFeedbacksByStatusReply, error) {
	feedbacks, total, err := s.uc.ListFeedbacksByStatus(ctx, req.Status, req.Page, req.PageSize)
	if err != nil {
		s.log.Errorf("Failed to list feedbacks by status: %v", err)
		return nil, err
	}
	
	var feedbackList []*v1.UserFeedbackData
	for _, feedback := range feedbacks {
		feedbackList = append(feedbackList, &v1.UserFeedbackData{
			// Id:        feedback.ID,
			// UserId:    feedback.UserID,
			Content:   feedback.Content,
			Type:      feedback.Type,
			Status:    feedback.Status,
			Rating:    feedback.Rating,
			CreatedAt: feedback.CreatedAt.Format(time.RFC3339),
			UpdatedAt: feedback.UpdatedAt.Format(time.RFC3339),
		})
	}
	
	return &v1.ListFeedbacksByStatusReply{
		Feedbacks: feedbackList,
		Total:     total,
	}, nil
}