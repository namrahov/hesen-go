package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/namrahov/hesen-go/model"
	"github.com/namrahov/hesen-go/repo"
	"github.com/namrahov/hesen-go/util"
	log "github.com/sirupsen/logrus"
	"net/http"
)

type IService interface {
	GetApplication(ctx context.Context, id int64) (*model.Application, error)
	ChangeStatus(ctx context.Context, id int64, request model.ChangeStatusRequest) *model.ErrorResponse
	SaveApplication(ctx context.Context, application *model.Application) (*model.Application, error)
}

type Service struct {
	ApplicationRepo repo.IApplicationRepo
	ValidationUtil  util.IValidationUtil
	CommentRepo     repo.ICommentRepo
}

func (s *Service) SaveApplication(ctx context.Context, application *model.Application) (*model.Application, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetApplication.start")

	application, err := s.ApplicationRepo.SaveApplication(application)
	if err != nil {
		logger.Errorf("ActionLog.SaveApplication.error: cannot save application for application id %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-save-application", model.Exception))
	}

	return application, nil
}

func (s *Service) GetApplication(ctx context.Context, id int64) (*model.Application, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetApplication.start")

	application, err := s.ApplicationRepo.GetApplicationById(id)
	if err != nil {
		logger.Errorf("ActionLog.GetApplication.error: cannot get application for application id %d, %v", id, err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-application", model.Exception))
	}

	logger.Info("ActionLog.GetApplication.success")

	return application, nil
}

func (s *Service) ChangeStatus(ctx context.Context, id int64, request model.ChangeStatusRequest) *model.ErrorResponse {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.ChangeStatus.start")

	application, err := s.ApplicationRepo.GetApplicationById(id)
	if err != nil {
		logger.Errorf("ActionLog.ChangeStatus.error: cannot get application for application id %d, %v", id, err)
		return &model.ErrorResponse{Code: err.Error(), Status: http.StatusInternalServerError}
	}

	err = s.ValidationUtil.ValidationApplicationStatus(application.Status, request.Status)
	if err != nil {
		log.Warn(fmt.Sprintf("ActionLog.ChangeStatus.error: %s -> %s is not possible", application.Status, request.Status))
		return &model.ErrorResponse{
			Code:   fmt.Sprintf("%s.Invalid status change from %s to %s", model.Exception, application.Status, request.Status),
			Status: http.StatusForbidden,
		}
	}

	comment := model.Comment{
		Commentator:   "Nurlan",
		Description:   request.Description,
		CommentType:   model.Internal,
		ApplicationId: application.Id,
	}

	if request.Status == model.Hold {
		err := s.CommentRepo.SaveComment(&comment)
		if err != nil {
			logger.Errorf("ActionLog.SaveComment.error: could not save comment for application id %d - %v", application.Id, err)
			return &model.ErrorResponse{Code: err.Error(), Status: http.StatusInternalServerError}
		}
	}

	application.Status = request.Status
	_, err = s.ApplicationRepo.UpdateApplication(application)
	if err != nil {
		return &model.ErrorResponse{Code: err.Error(), Status: http.StatusInternalServerError}
	}

	logger.Info("ActionLog.ChangeStatus.end")
	return nil
}
