package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/namrahov/hesen-go/model"
	"github.com/namrahov/hesen-go/repo"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type IUserService interface {
	GetUserIfExist(ctx context.Context, sessionId string) (*model.User, error)
	SaveUser(ctx context.Context, u *model.UserRegister) (*model.User, error)
	SaveSession(ctx context.Context, sessionId string, userId string) error
	AlreadyLoggedIn(r *http.Request) bool
}

type UserService struct {
	UserRepo repo.IUserRepo
}

func (up *UserService) AlreadyLoggedIn(r *http.Request) bool {
	cookie, err := r.Cookie("session-id")

	if err != nil {
		return false
	}

	userAddress, err := up.GetUserIfExist(r.Context(), cookie.Value)
	if err == nil && userAddress != nil {
		return true
	} else {
		return false
	}

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

	u, err := up.UserRepo.GetUserById(sessions[len(sessions)-1].UserId)

	if err1 != nil {
		logger.Errorf("ActionLog.GetUserIfExist.error: cannot get user for user id %v, %v",
			sessions[len(sessions)-1].UserId, err)
		return nil, errors.New(fmt.Sprintf("%s.can't-get-user", model.Exception))
	}

	logger.Info("ActionLog.GetUserIfExist.success")

	return u, nil
}

func (up *UserService) SaveUser(ctx context.Context, u *model.UserRegister) (*model.User, error) {
	logger := ctx.Value(model.ContextLogger).(*log.Entry)
	logger.Info("ActionLog.SaveUser.start")

	ps, err1 := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.MinCost)
	if err1 != nil {
		logger.Errorf("ActionLog.SaveUser.error: cannot encrypt password for user id")
		return nil, errors.New(fmt.Sprintf("%s.can't-encrypt-password", model.Exception))
	}

	user := model.User{
		UserName:  u.UserName,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Password:  ps,
	}

	us, err := up.UserRepo.SaveUser(&user)
	if err != nil {
		logger.Errorf("ActionLog.GetUserIfExist.error: cannot save user for user id %v", err)
		return nil, errors.New(fmt.Sprintf("%s.can't-save-user", model.Exception))
	}

	logger.Info("ActionLog.SaveUser.success")

	return us, nil
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
