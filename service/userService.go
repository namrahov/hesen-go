package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/namrahov/hesen-go/model"
	"github.com/namrahov/hesen-go/repo"
	log "github.com/sirupsen/logrus"
)

type IUserService interface {
	GetUserIfExist(ctx context.Context, sessionId string) (*model.User, error)
	SaveUser(ctx context.Context, u *model.User) (*model.User, error)
	SaveSession(ctx context.Context, sessionId string, userId string) error
}

type UserService struct {
	UserRepo repo.IUserRepo
}

func (up *UserService) GetUserIfExist(ctx context.Context, sessionId string) (*model.User, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.GetUserIfExist.start")

	sessionsAddresses, err1 := up.UserRepo.GetSessionBySessionId(sessionId)
	sessions := *sessionsAddresses

	if err1 != nil {
		logger.Errorf("ActionLog.GetUserIfExist.error: cannot get session for session id %v, %v", sessionId, err1)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-session", model.Exception))
	}
	//the user doesn't exist
	if len(sessions) == 0 {
		return nil, nil
	}

	fmt.Println("sessions[0]=", sessions[0])

	u, err := up.UserRepo.GetUserById(sessions[0].UserId)

	if err1 != nil {
		logger.Errorf("ActionLog.GetUserIfExist.error: cannot get user for user id %v, %v", sessions[0].UserId, err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-user", model.Exception))
	}

	logger.Info("ActionLog.GetUserIfExist.success")

	return u, nil
}

func (up *UserService) SaveUser(ctx context.Context, u *model.User) (*model.User, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.SaveUser.start")

	user, err := up.UserRepo.SaveUser(u)
	if err != nil {
		logger.Errorf("ActionLog.GetUserIfExist.error: cannot save user for user id %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-save-user", model.Exception))
	}

	logger.Info("ActionLog.SaveUser.success")

	return user, nil
}

func (up *UserService) SaveSession(ctx context.Context, sessionId string, userId string) error {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.SaveUser.start")

	session := model.Session{
		SessionId: sessionId,
		UserId:    userId,
	}

	err := up.UserRepo.SaveSession(&session)
	if err != nil {
		logger.Errorf("ActionLog.SaveSession.error: cannot save user for user id %v", err)
		return errors.New(fmt.Sprintf("%s.can't-save-session", model.Exception))
	}

	logger.Info("ActionLog.SaveUser.success")

	return nil
}
