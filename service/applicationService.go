package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/namrahov/hesen-go/model"
	"github.com/namrahov/hesen-go/repo"
	log "github.com/sirupsen/logrus"
)

type IService interface {
	GetApplication(ctx context.Context, id int64) (*model.Application, error)
}

type Service struct {
	ApplicationRepo repo.IApplicationRepo
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
