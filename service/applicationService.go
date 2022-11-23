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
	"sync"
	"time"
)

type IService interface {
	GetApplication(ctx context.Context, id int64) (*model.Application, error)
	ChangeStatus(ctx context.Context, id int64, request model.ChangeStatusRequest) *model.ErrorResponse
	SaveApplication(ctx context.Context, application *model.Application) (*model.Application, error)
	GetFilterInfo(ctx context.Context) (*model.FilterInfo, error)
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

func (s *Service) GetFilterInfo(ctx context.Context) (*model.FilterInfo, error) {
	stratTime := time.Now()

	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetFilterInfo.start")

	var applicationWithDistinctCourtName *[]model.Application
	var applicationWithDistinctJudgeName *[]model.Application
	var err error

	var courts []string
	var judges []string

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		applicationWithDistinctCourtName, err = s.ApplicationRepo.GetDistinctCourtName(&wg)
		for _, application := range *applicationWithDistinctCourtName {
			if application.CourtName != "" {
				courts = append(courts, application.CourtName)
			}
		}
	}()

	if err != nil {
		logger.Errorf("ActionLog.GetFilterInfo.error: cannot get applicationWithDistinctCourtName info %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-application-with-distinct-court-name", model.Exception))
	}

	go func() {
		applicationWithDistinctJudgeName, err = s.ApplicationRepo.GetDistinctJudgeName(&wg)
		for _, application := range *applicationWithDistinctJudgeName {
			if application.JudgeName != "" {
				judges = append(judges, application.JudgeName)
			}
		}
	}()

	if err != nil {
		logger.Errorf("ActionLog.GetFilterInfo.error: cannot get applicationWithDistinctJudgeName info %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-application-with-distinct-judge-name", model.Exception))
	}

	wg.Wait()

	filterInfo := model.FilterInfo{
		Courts: courts,
		Judges: judges,
	}

	logger.Info("ActionLog.GetFilterInfo.success")

	fmt.Println(time.Now().Sub(stratTime))

	return &filterInfo, nil
}
